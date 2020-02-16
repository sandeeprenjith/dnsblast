# DNS Blast
> A DNS performance testing utility


[![Go Report Card](https://goreportcard.com/badge/github.com/sandeeprenjith/dnsblast)](https://goreportcard.com/report/github.com/sandeeprenjith/dnsblast)

| Currently Supported Protocols |
| ----------------------------- |
| UDP DNS                      	|
| TCP DNS                      	|
| TCP DNS over TLS             	|


*Mon Feb 10 17:18:41 DST 2020*

## Installation 

#### Download the executable

##### v1
Go to [downloads](https://github.com/sandeeprenjith/dnsblast/tree/v1/builds)

#### Build for your platform

##### Requirements

* go

```
go get github.com/sandeeprenjith/dnsblast

```

#### Build for other platforms 

##### Requirements

* go
* make

```
$ git clone https://github.com/sandeeprenjith/dnsblast.git
$ cd dnsblast
$ make
```
This creates a directory named "builds". The directory contains archives with executables for different platforms. 

```
╰$ tree builds
builds
├ dnsblast-v1-darwin-386.tar.gz
├ dnsblast-v1-darwin-amd64.tar.gz
├ dnsblast-v1-linux-386.tar.gz
├ dnsblast-v1-linux-amd64.tar.gz
├ dnsblast-v1-linux-arm.tar.gz
├ dnsblast-v1-linux-arm64.tar.gz

```
The archives contain the executable for the platform which the name of the archive suggests.

```
$ tar -tf builds/dnsblast-v1-linux-386.tar.gz
dnsblast
```


## Usage

```
$ ./dnsblast -h
Usage of ./dnsblast:
  -c int
        Value 0 for random QNAMES (for uncached responses), 100 for Predictable QNAMES (for cached responses)
  -l int
        Duration to run load (default 60)
  -p string
        The destination UDP port (default "53")
  -proto string
        Protocol to use for DNS queries ( udp, tcp, tls) (default "udp")
  -r int
        Packets per second to send (default 100)
  -s string
        [Required] The address of the target server
  -t int
        Number of threads (default 4)

```

## Sample Output

> Tested against coredns running DNS over TLS with errata plugin configured to give fake responses.

```

$ ./dnsblast -s 192.168.130.9 -r 100 -t 1 -c 100 -proto tls -l 5
Putting entries in cache; hoping no 0 ttl responses
100 / 100 [-------------------------------------------------------] 100.00% 24 p/s

EXECUTING TEST
+-----------------------------------------------------------+
2020/02/13 22:35:37 QPS/Thread:  51  Latency:  6.645784ms
2020/02/13 22:35:38 QPS/Thread:  81  Latency:  8.286109ms
2020/02/13 22:35:39 QPS/Thread:  76  Latency:  12.941233ms
2020/02/13 22:35:40 QPS/Thread:  79  Latency:  15.883037ms
2020/02/13 22:35:41 QPS/Thread:  71  Latency:  21.37718ms
+-----------------------------------------------------------+

  REPORT
+---------------------+-------------------------+
| Target Server       | tls://192.168.130.9:853 |
| Test                | Cached Responses        |
| Send Rate           | 100 Queries/Sec         |
| Threads             | 1                       |
| Duration of test    | 5 Sec                   |
| Protocol            | TCP-TLS                 |
| Average Queries/Sec | 71                      |
| Average Latency     | 13.026668ms             |
+---------------------+-------------------------+
```
## Credit where due

* https://github.com/miekg/dns
