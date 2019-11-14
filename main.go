package main

import (
	"flag"
	"github.com/sandeeprenjith/dnsblast/qry"
	"log"
	"os"
	"time"
)

/* For debug
func checkpoint(n int) {
	log.Println("Checkpoint: ", n)
}
*/

func main() {
	// Getting input from user
	server := flag.String("server", "127.0.0.1", "The address of the target server")
	num := flag.Int("rate", 100, "Packets per second to send")
	port := flag.String("port", "53", "The destination UDP port")
	duration := flag.Int("len", 60, "Duration to run load")
	flag.Parse()

	limiter := time.Tick(time.Second) // Ticker used for rate limiting packets per second

	// Exit program when the time specified by "-len" is passed.
	ender := time.Tick(time.Duration(*duration) * time.Second)
	go func() {
		<-ender
		os.Exit(0)
	}()

	var QPS int             // Variable to hold QPS
	var RTT []time.Duration // Variable to hold Latency

	// The eternal for loop runs till program is killed by the ender ticker
	for {
		responses := make(chan qry.Response, *num) // Channel to hold DNS responses
		// loop which runs for a maximum of "-rate" specified by user.
	rateLoop: // Issue #2 rate limit execution. Added to use in break statement.
		for i := 1; i <= *num; i++ {
			select { // Issue #2 Rate limit execution
			case <-limiter:
				break rateLoop
			default:
				qname := qry.PQname(3, i)                              // Creating a predictable Qname
				qry.SimpleQuery(*server, *port, qname, "A", responses) // Query the specified server with the predictable qname
			}
		}
		close(responses)

		// Iterate through the responses channel and if RCODE is NOERROR, increment QPS and append Latency to array
		for resp := range responses {
			if resp.Rcode == "NOERROR" {
				QPS++
				RTT = append(RTT, resp.Rtt)
			}
		}
		//<-limiter // Limit the execution. Will block till 1 second passes

		var sumRTT time.Duration // Varable to hold sum of latency values
		var avgRTT time.Duration // Variable to hold avg of latency values

		// Calculate sum of latency ( for avg calculation) and calculate avg
		for x := 1; x < len(RTT); x++ {
			sumRTT = sumRTT + RTT[x]
		}
		avgRTT = sumRTT / time.Duration(len(RTT))

		log.Println("QPS: ", QPS, " Latency: ", avgRTT)
		QPS = 0                                 // Reinitialize QPS for next iteration
		RTT = []time.Duration{time.Duration(0)} // Reinitialize latency array for next iteration
	}

}
