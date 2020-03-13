package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/google/goterm/term"
	"github.com/olekukonko/tablewriter"
	"github.com/sandeeprenjith/dnsblast/qry"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Struct for the data returned by the send_query function (to channel).
// the data includes sum of QPS and average round trip time.
type Results struct {
	QPS int
	RTT time.Duration
}

func send_qry(qnamelist []string, server string, rate int, port string, duration int, threads int, limiter <-chan time.Time, res chan Results, ender <-chan time.Time, proto string, chr int) {

	//handling panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(term.Redf("Looks like the program paniced."))
		}
	}()

	var QPS int             // Variable to hold QPS
	var RTT []time.Duration // Variable to hold Latency
	var qname string
	var result Results
	var resultset []Results
	var sumRTT time.Duration // Varable to hold sum of latency values
	var avgRTT time.Duration // Variable to hold avg of latency values
	var total_qps int
	var total_avg_rtt time.Duration
	var avg_avg_rtt time.Duration
	var avg_total_qps int
	var final_results Results
	var qps_denominator int
	total_qps = 0
	total_avg_rtt = 0
	avg_avg_rtt = 0
	num := rate / threads

mainLoop:
	// The eternal for loop runs till program is killed by the ender ticker
	for {
		select { //Break main loop after duration expires
		case <-ender:
			break mainLoop
		default:
			responses := make(chan qry.Response, num) // Channel to hold DNS responses
			// loop which runs for a maximum of "-rate" specified by user.

		rateLoop: // Issue #2 rate limit execution. Added to use in break statement.
			for i := 1; i <= num; i++ {
				select { // Issue #2 Rate limit execution
				case <-limiter:
					break rateLoop
				default:
					if len(qnamelist) == 0 {
						if chr == 100 {
							// Creating a predictable Qname
							qname = qry.PQname(3, i)
						} else if chr == 0 {
							// Creating an unpredictable Qname
							qname = qry.RQname(3)
						}
					} else {
						rand.Seed(time.Now().UnixNano())
						qname = qnamelist[rand.Intn(len(qnamelist))]
					}
					qry.SimpleQuery(server, port, qname, "A", responses, proto) // Query the specified server with the predictable qname
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

			// Calculate sum of latency ( for avg calculation) and calculate avg
			for x := 1; x < len(RTT); x++ {
				sumRTT = sumRTT + RTT[x]
			}
			if sumRTT == 0 {
				log.Println(term.Redf("No usable results"))
			} else {
				avgRTT = sumRTT / time.Duration(len(RTT))

				//Send results to a chanel in the type Result
				result.QPS = QPS
				result.RTT = avgRTT
				resultset = append(resultset, result)

				log.Println(term.Cyanf("QPS/Thread: "), term.Yellowf(strconv.Itoa(QPS)), term.Cyanf(" Latency: "), term.Yellowf(avgRTT.String())) // Print result per iteration ( minumum rate times/sec)
				QPS = 0                                                                                                                           // Reinitialize QPS for next iteration
				RTT = []time.Duration{time.Duration(0)}                                                                                           // Reinitialize latency array for next iteration
			}
		}
	}
	for x := range resultset {
		total_qps = total_qps + resultset[x].QPS
		total_avg_rtt = total_avg_rtt + resultset[x].RTT
	}
	if len(resultset) == 0 {
		os.Exit(1)
	}
	qps_denominator = len(resultset) / threads

	//Handling zero division panic in a way that does not affect the final result
	if qps_denominator == 0 {
		qps_denominator = 1
		fmt.Println(term.Redf("A thread ran into an issue. Consider reducing the number of threads"))

	}
	avg_total_qps = total_qps / qps_denominator
	avg_avg_rtt = total_avg_rtt / time.Duration(len(resultset))
	final_results.QPS = avg_total_qps
	final_results.RTT = avg_avg_rtt
	res <- final_results
}

func main() {
	// Getting input from user
	server := flag.String("s", "", "[Required] The address of the target server")
	rate := flag.Int("r", 100, "Packets per second to send")
	port := flag.String("p", "53", "The destination UDP port")
	duration := flag.Int("l", 60, "Duration to run load")
	threads := flag.Int("t", 4, "Number of threads")
	protocol := flag.String("proto", "udp", "Protocol to use for DNS queries ( udp, tcp, tls)")
	infile := flag.String("f", "", "Input file with query names")
	chr := flag.Int("c", 0, "Value 0 for random QNAMES (for uncached responses), 100 for Predictable QNAMES (for cached responses)")
	flag.Parse()

	var proto string
	var test string
	var qnamelist []string // Variable to store qnames from query file

	if *infile != "" {
		//Handle query names from file
		inFile, err := os.Open(*infile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer inFile.Close()
		scanner := bufio.NewScanner(inFile)
		for scanner.Scan() {
			qnamelist = append(qnamelist, scanner.Text())
		}
	}

	if *protocol == "udp" {
		proto = "udp"
	} else if *protocol == "tcp" {
		proto = "tcp"
	} else if *protocol == "tls" {
		proto = "tcp-tls"
	}

	if *chr == 0 {
		test = "Uncached Responses"
	} else if *chr == 100 {
		test = "Cached Responses"
	} else {
		fmt.Println("Value of '-c' flag must be either 0 or 100")
		os.Exit(1)
	}
	// Changing test name when using input files
	if *infile != "" {
		test = "Queries from file"
	}
	// Just in case user specifies DNS over TLS but does not modify port
	qport := *port
	if proto == "tcp-tls" {
		if *port == "53" {
			qport = "853"
		} else {
			qport = *port
		}
	}
	if *server == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	num := *rate / *threads
	//Priming the cache for cached response testing
	if *chr == 100 {
		fmt.Println("Putting entries in cache; hoping no 0 ttl responses")
		bar := pb.StartNew(num)
		for i := 1; i <= num; i++ {
			unused := make(chan qry.Response, num)
			qname := qry.PQname(3, i)
			qry.SimpleQuery(*server, qport, qname, "A", unused, proto)
			bar.Increment()
		}
		bar.Finish()
	}

	fmt.Println(" ")
	fmt.Println(term.Cyanf("EXECUTING TEST"))
	fmt.Println(term.Yellowf("+-----------------------------------------------------------+"))

	limiter := time.Tick(time.Second) // Ticker used for rate limiting packets per second
	res := make(chan Results, *threads)

	var total_qps int
	var total_rtt time.Duration
	// Channel which gets a value when the timer specified with duration has expired
	// Used to end the main forloop in the send_qry function
	ender := time.Tick(time.Duration(*duration) * time.Second)
	// Create as many goroutines as specified by "-t" argument
	for i := 1; i <= *threads; i++ {
		go send_qry(qnamelist, *server, *rate, qport, *duration, *threads, limiter, res, ender, proto, *chr)
	}
	sleepval := *duration + 1
	time.Sleep(time.Duration(sleepval) * time.Second)
	// Iterate over resuls channel and calculate QPS and RTT
	total_qps = 0
	total_rtt = 0

	close(res)
	for each_res := range res {
		total_qps = total_qps + each_res.QPS
		total_rtt = total_rtt + each_res.RTT
	}

	// Formatting final results into an ascii table
	// Formatting user given data to print
	tabledata := [][]string{
		[]string{"Target Server", *protocol + "://" + *server + ":" + qport},
		[]string{"Test", test},
		[]string{"Send Rate", strconv.Itoa(*rate) + " Queries/Sec"},
		[]string{"Threads", strconv.Itoa(*threads)},
		[]string{"Duration of test", strconv.Itoa(*duration) + " Sec"},
		[]string{"Protocol", strings.ToUpper(proto)},
		[]string{"Average Queries/Sec", term.Yellowf(strconv.Itoa(total_qps))},
		[]string{"Average Latency", term.Yellowf(total_rtt.String())},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, dat := range tabledata {
		table.Append(dat)
	}
	time.Sleep(3 * time.Second)
	fmt.Println(term.Yellowf("+-----------------------------------------------------------+"))
	fmt.Println(" ")
	fmt.Println(term.Cyanf("  REPORT"))
	table.Render()

	os.Exit(0)
}
