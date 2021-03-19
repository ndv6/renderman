package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gookit/cache"
	"github.com/sirupsen/logrus"
)

type redishash struct {
	url   string
	pwd   string
	dbNum int
	pool  *redis.Pool
}

func redishashNew(url, pwd string, dbNum int) *redishash {
	c := &redishash{
		url:   url,
		pwd:   pwd,
		dbNum: dbNum,
	}
	return c
}

func (c *redishash) Clear() error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("FlushDb")
	return err
}

func (c *redishash) Close() error {
	return c.pool.Close()
}

func splitKey(key string) []string {
	return strings.SplitN(key, ":", 2)
}

func ttlKey(key string) string {
	return fmt.Sprintf("%s_ttl", key)
}

func (c *redishash) Del(key string) error {
	slc := splitKey(key)
	if len(slc) == 1 {
		_, err := c.exec("Del", key)
		if err != nil {
			logrus.Errorf("redis error: %s\n", err.Error())
		}
		return err
	}

	hname := slc[0]
	key = slc[1]
	_, err := c.exec("HDel", hname, key)
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
	}
	_, err = c.exec("HDel", hname, ttlKey(key))
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
	}
	return err
}

func (c *redishash) DelMulti(key []string) error {
	m := map[string][]interface{}{}
	for _, k := range key {
		slc := splitKey(k)
		if len(slc) == 1 {
			m[""] = append(m[""], k)
			continue
		}
		hname := slc[0]
		hfield := slc[1]
		m[hname] = append(m[hname], hfield, ttlKey(hfield))
	}
	conn := c.pool.Get()
	defer conn.Close()
	err := conn.Send("Multi")
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
		return err
	}
	for hname, keys := range m {
		if hname == "" {
			_, err := conn.Do("Del", keys...)
			if err != nil {
				logrus.Errorf("redis error: %s\n", err.Error())
				return err
			}
			continue
		}
		for _, k := range keys {
			_, err := conn.Do("HDel", hname, k)
			if err != nil {
				logrus.Errorf("redis error: %s\n", err.Error())
				return err
			}
		}
	}
	_, err = redis.Values(conn.Do("Exec"))
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
	}
	return err
}

func (c *redishash) Get(key string) interface{} {
	slc := splitKey(key)
	if len(slc) == 1 {
		_, err := c.exec("Get", key)
		if err != nil {
			logrus.Errorf("redis error: %s\n", err.Error())
		}
		return err
	}
	now := time.Now().Unix()
	hname := slc[0]
	key = slc[1]
	ttlkey := ttlKey(key)
	infs, err := redis.Values(c.exec("HMGet", hname, key, ttlkey))
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
		return nil
	}
	v := infs[0]
	vttl := infs[1]
	if v == nil {
		return nil
	}
	ttl, err := redis.Int64(vttl, nil)
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
		return nil
	}

	if now < ttl {
		return v
	}
	_, err = c.exec("HDel", hname, key)
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
		return nil
	}
	_, err = c.exec("HDel", hname, ttlkey)
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
		return nil
	}
	return nil
}

func (c *redishash) GetMulti(key []string) map[string]interface{} {
	now := time.Now().Unix()
	res := map[string]interface{}{}
	m := map[string][]interface{}{}
	for _, k := range key {
		slc := splitKey(k)
		if len(slc) == 1 {
			m[""] = append(m[""], k)
			continue
		}
		hname := slc[0]
		hfield := slc[1]
		m[hname] = append(m[hname], hfield, ttlKey(hfield))
	}
	conn := c.pool.Get()
	defer conn.Close()
	for hname, keys := range m {
		if hname == "" {
			list, err := redis.Values(conn.Do("MGet", keys...))
			if err != nil {
				logrus.Errorf("redis error: %s\n", err.Error())
				return nil
			}
			for i, e := range list {
				if str, ok := keys[i].(string); ok {
					res[str] = e
				}
			}
			continue
		}
		args := []interface{}{hname}
		args = append(args, keys...)
		list, err := redis.Values(conn.Do("HMGet", args...))
		if err != nil {
			logrus.Errorf("redis error: %s\n", err.Error())
			return nil
		}
		for i := 0; i < len(list); i += 2 {
			var (
				resKey = fmt.Sprintf("%s:%s", hname, keys[i])
				v      = list[i]
				vttl   = list[i+1]
			)
			if v == nil {
				res[resKey] = v
				continue
			}

			ttl, err := redis.Int64(vttl, nil)
			if err != nil {
				logrus.Errorf("redis error: %s\n", err.Error())
				return nil
			}

			if now < ttl {
				res[resKey] = v
				continue
			}
			res[resKey] = nil
			_, err = conn.Do("HDel", hname, keys[i])
			if err != nil {
				logrus.Errorf("redis error: %s\n", err.Error())
				return nil
			}
			_, err = conn.Do("HDel", hname, keys[i+1])
			if err != nil {
				logrus.Errorf("redis error: %s\n", err.Error())
				return nil
			}

		}
	}
	return res
}

func (c *redishash) Has(key string) bool {
	slc := splitKey(key)
	if len(slc) == 1 {
		one, err := redis.Int(c.exec("Exists", key))
		if err != nil {
			logrus.Errorf("redis error: %s\n", err.Error())
			return false
		}
		return one == 1
	}
	hname := slc[0]
	key = slc[1]
	one, err := redis.Int(c.exec("HExists", hname, key))
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
		return false
	}
	if one != 1 {
		return false
	}
	ttl, err := redis.Int64(c.exec("HGet", hname, ttlKey(key)))
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
		return false
	}
	return time.Now().Unix() < ttl
}

func (c *redishash) Set(key string, val interface{}, ttl time.Duration) error {
	slc := splitKey(key)
	if len(slc) == 1 {
		_, err := c.exec("SetEx", key, int64(ttl/time.Second), val)
		if err != nil {
			logrus.Errorf("redis error: %s\n", err.Error())
		}
		return err
	}
	hname := slc[0]
	key = slc[1]
	until := time.Now().Add(ttl).Unix()
	_, err := c.exec("HMSet", hname, key, val, ttlKey(key), until)
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
	}
	return err
}

func (c *redishash) SetMulti(values map[string]interface{}, ttl time.Duration) error {
	until := time.Now().Add(ttl).Unix()
	m := map[string]map[string]interface{}{}
	for k, v := range values {
		slc := splitKey(k)
		if len(slc) == 1 {
			im := m[""]
			if len(im) == 0 {
				im = map[string]interface{}{}
				m[""] = im
			}
			im[k] = v
			continue
		}
		hname := slc[0]
		hfield := slc[1]
		im := m[hname]
		if len(im) == 0 {
			im = map[string]interface{}{}
			m[hname] = im
		}
		im[hfield] = v
		im[ttlKey(hfield)] = until
	}
	conn := c.pool.Get()
	defer conn.Close()
	err := conn.Send("Multi")
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
		return err
	}
	for hname, keys := range m {
		if hname == "" {
			for k, v := range keys {
				_, err := conn.Do("SetEx", k, int64(ttl/time.Second), v)
				if err != nil {
					logrus.Errorf("redis error: %s\n", err.Error())
					return err
				}
			}
			continue
		}
		args := []interface{}{hname}
		for k, v := range keys {
			args = append(args, k, v)
		}
		_, err := conn.Do("HMSet", args...)
		if err != nil {
			logrus.Errorf("redis error: %s\n", err.Error())
			return err
		}
	}
	_, err = redis.Strings(conn.Do("Exec"))
	if err != nil {
		logrus.Errorf("redis error: %s\n", err.Error())
	}
	return err
}

func redishashConnect(url, pwd string, dbNum int) cache.Cache {
	c := redishashNew(url, pwd, dbNum)
	c.pool = newPool(c.url, c.pwd, c.dbNum)
	logrus.Infof("connect to server %s db is %d", c.url, c.dbNum)
	return c
}

func (c *redishash) exec(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	conn := c.pool.Get()
	defer conn.Close()
	st := time.Now()
	defer func() {
		logrus.Debugf(
			"operate redis cache. command: %s, key: %v, elapsed time: %.03f\n",
			commandName, args[0], time.Since(st).Seconds()*1000,
		)
	}()

	return conn.Do(commandName, args...)
}

// create new pool
func newPool(url, password string, dbNum int) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 5,
		// timeout
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", url)
			if err != nil {
				return nil, err
			}

			if password != "" {
				_, err := c.Do("AUTH", password)
				if err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			_, _ = c.Do("SELECT", dbNum)
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
