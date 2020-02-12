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

```

$ dnsblast -s 192.168.130.9 -r 3000 -l 10
2019/11/14 14:58:55 QPS:  1858  Latency:  237.799µs
2019/11/14 14:58:56 QPS:  1858  Latency:  236.285µs
2019/11/14 14:58:57 QPS:  1847  Latency:  233.041µs
2019/11/14 14:58:58 QPS:  1881  Latency:  232.893µs
2019/11/14 14:58:59 QPS:  1867  Latency:  233.41µs
2019/11/14 14:59:00 QPS:  1866  Latency:  232.374µs
2019/11/14 14:59:01 QPS:  1880  Latency:  229.484µs
2019/11/14 14:59:02 QPS:  1882  Latency:  230.578µs
2019/11/14 14:59:03 QPS:  1810  Latency:  237.021µs

```

