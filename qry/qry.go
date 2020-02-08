package qry

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"time"
)

func HandleErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func SimpleQuery(server string, port string, qname string, qtype string, responses chan Response) {
	s_server := server + ":" + port
	qrytype := Qtype(qtype)
	question := new(dns.Msg)
	question.SetQuestion(dns.Fqdn(qname), qrytype)
	c := new(dns.Client)
	c.Dialer = &net.Dialer{
		Timeout: 200 * time.Millisecond,
	}
	c.SingleInflight = true
	c.Net = "tcp-tls"
	ans, rtt, err := c.Exchange(question, s_server)
	HandleErr(err)
	var R Response
	R.Rcode = Rcode(ans.Rcode)
	R.Rtt = rtt
	R.Qname = qname
	R.Server = s_server
	R.Qtype = qtype
	responses <- R
}
