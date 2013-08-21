go-proxymate
============

A simple client/server tool to assist people who need to work on TCP proxies.

Installation
------------

`go get github.com/politician/go-proxymate`

Usage
-----

```bash
$ go build
$ ./go-proxymate --help
```

Example Output
--------------

```Bash
$ ./go-proxymate -s ":8081"
2013/08/21 02:24:53 Listening on :8081
2013/08/21 02:24:53 connection accepted: 127.0.0.1:50438
2013/08/21 02:24:53 CTRL-C to exit...
([]uint8) {
 00000000  70 69 6e 67 00 00 00 00  00 00 00 00 00 00 00 00  |ping............|
}
([]uint8) {
 00000000  70 6f 6e 67 00 00 00 00  00 00 00 00 00 00 00 00  |pong............|
 00000010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
}
([]uint8) {
 00000000  70 69 6e 67 00 00 00 00  00 00 00 00 00 00 00 00  |ping............|
}
([]uint8) {
 00000000  70 6f 6e 67 00 00 00 00  00 00 00 00 00 00 00 00  |pong............|
 00000010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
}
([]uint8) {
 00000000  70 69 6e 67 00 00 00 00  00 00 00 00 00 00 00 00  |ping............|
}
([]uint8) {
 00000000  70 6f 6e 67 00 00 00 00  00 00 00 00 00 00 00 00  |pong............|
 00000010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
}
^C2013/08/21 02:25:00 Got:  interrupt
2013/08/21 02:25:00 stopping client
2013/08/21 02:25:00 stopping server
2013/08/21 02:25:00 refusing new connections to :8081
2013/08/21 02:25:00 waiting for existing requests to drain...
2013/08/21 02:25:00 connection closed: 127.0.0.1:50438
2013/08/21 02:25:00 remaining connections have drained
2013/08/21 02:25:00 goodbye
```
