package main

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Content ...
type Content struct {
	Title   string
	Content string
	Header  http.Header
}

// Marshal ...
func (c *Content) Marshal() string {
	content, err := json.Marshal(c)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return string(content)
}

// UnMarshal ...
func (c *Content) UnMarshal(val []byte) {
	err := json.Unmarshal(val, c)
	if err != nil {
		logrus.Error(err)
	}
}
