# DNS Blast
> A DNS performance testing utility

**Note:** Not in a usable stage yet 

*Tue Oct 29 21:37:12 DST 2019*

## Installation 

```
go get github.com/sandeeprenjith/dnsblast

```

## Usage

```

$ dnsblast -h
Usage of dnsblast:
  -len int
        Duration to run load (default 60)
  -port string
        The destination UDP port (default "53")
  -rate int
        Packets per second to send (default 100)
  -server string
        The address of the target server (default "127.0.0.1")

```

## Sample Output

```
╭nd33p@trv-wl-srenjith ~/go/src/github.com/sandeeprenjith/dnsblast ‹master›
╰$ go run main.go -server 192.168.130.9 -rate 10000 -len 10
2019/11/14 14:54:21 QPS:  1892  Latency:  234.101µs
2019/11/14 14:54:22 QPS:  1867  Latency:  230.924µs
2019/11/14 14:54:23 QPS:  1827  Latency:  237.156µs
2019/11/14 14:54:24 QPS:  1865  Latency:  232.912µs
2019/11/14 14:54:25 QPS:  1814  Latency:  237.668µs
2019/11/14 14:54:26 QPS:  1872  Latency:  232.288µs
2019/11/14 14:54:27 QPS:  1726  Latency:  250.058µs
2019/11/14 14:54:28 QPS:  1831  Latency:  241.003µs
2019/11/14 14:54:29 QPS:  1858  Latency:  232.634µs

```

