# DNS Blast
> A DNS performance testing utility

*Note: * Not in a usable stage yet 
Tue Oct 29 21:37:12 DST 2019

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
$ dnsblast -server 192.168.130.9 -len 15 -rate 50
2019/10/29 21:35:24 QPS:  50  Latency:  234.78µs
2019/10/29 21:35:25 QPS:  50  Latency:  422.745µs
2019/10/29 21:35:26 QPS:  50  Latency:  382.235µs
2019/10/29 21:35:27 QPS:  50  Latency:  450.666µs
2019/10/29 21:35:28 QPS:  50  Latency:  393.156µs
2019/10/29 21:35:29 QPS:  50  Latency:  436.019µs
2019/10/29 21:35:30 QPS:  50  Latency:  341.156µs
2019/10/29 21:35:31 QPS:  50  Latency:  248.568µs
2019/10/29 21:35:32 QPS:  50  Latency:  391.901µs
2019/10/29 21:35:33 QPS:  50  Latency:  360.274µs
2019/10/29 21:35:34 QPS:  50  Latency:  356.764µs
2019/10/29 21:35:35 QPS:  50  Latency:  380.333µs
2019/10/29 21:35:36 QPS:  50  Latency:  228.019µs
2019/10/29 21:35:37 QPS:  50  Latency:  438.137µs
2019/10/29 21:35:38 QPS:  50  Latency:  428.941µs

```
