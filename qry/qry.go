package qry

import (
	"crypto/tls"
	"github.com/google/goterm/term"
	"github.com/miekg/dns"
	"log"
	"net"
	"sync"
	"time"
)

// Send a single DNS query
func SimpleQuery(server string,
	port string,
	qname string,
	qtype string,
	responses chan Response,
	proto string,
	wg *sync.WaitGroup,
	noverify bool) {
	s_server := net.JoinHostPort(server, port)
	qrytype := Qtype(qtype)
	question := new(dns.Msg)
	question.SetQuestion(dns.Fqdn(qname), qrytype)
	c := new(dns.Client)
	c.Dialer = &net.Dialer{
		Timeout: 5 * time.Second,
	}
	c.DialTimeout = 5 * time.Second
	c.Net = proto
	if proto == "tcp-tls" {
		var tlc tls.Config
		tlc.InsecureSkipVerify = noverify
		c.TLSConfig = &tlc
	}
	ans, rtt, err := c.Exchange(question, s_server)
	if err != nil {
		log.Println(term.Redf(err.Error()))
	} else {
		var R Response
		R.Rcode = Rcode(ans.Rcode)
		R.Rtt = rtt
		R.Qname = qname
		R.Server = s_server
		R.Qtype = qtype
		responses <- R
	}
	wg.Done()
}
