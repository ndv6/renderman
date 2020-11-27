To build this launcher:

```
docker build -f chromium.dockerfile -t go-rod/launcher:1.0.0 .
```

To run this launcher:
```
 docker run -it -p 9222:9222 --name rod-launcher go-rod/launcher:1.0.0
```