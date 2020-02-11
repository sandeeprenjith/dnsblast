package main

import (
	"flag"
	"fmt"
	"github.com/google/goterm/term"
	"github.com/olekukonko/tablewriter"
	"github.com/sandeeprenjith/dnsblast/qry"
	"log"
	"os"
	"strconv"
	"time"
)

// For debug
//func checkpoint(n int) {
//	log.Println("Checkpoint: ", n)
//}

// Struct for the data returned by the send_query function (to channel).
// the data includes sum of QPS and average round trip time.
type Results struct {
	QPS int
	RTT time.Duration
}

func send_qry(server string, rate int, port string, duration int, threads int, limiter <-chan time.Time, res chan Results, ender <-chan time.Time) {
	var QPS int             // Variable to hold QPS
	var RTT []time.Duration // Variable to hold Latency

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

mainLoop:
	// The eternal for loop runs till program is killed by the ender ticker

	for {
		select { //Break main loop after duration expires
		case <-ender:
			break mainLoop
		default:
			num := rate / threads
			//		print(string(num))
			responses := make(chan qry.Response, num) // Channel to hold DNS responses
			// loop which runs for a maximum of "-rate" specified by user.

		rateLoop: // Issue #2 rate limit execution. Added to use in break statement.
			for i := 1; i <= num; i++ {
				select { // Issue #2 Rate limit execution
				case <-limiter:
					break rateLoop
				default:
					qname := qry.PQname(3, i)                            // Creating a predictable Qname
					qry.SimpleQuery(server, port, qname, "A", responses) // Query the specified server with the predictable qname
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
	port := flag.String("p", "853", "The destination UDP port")
	duration := flag.Int("l", 60, "Duration to run load")
	threads := flag.Int("t", 4, "Number of threads")
	flag.Parse()
	if *server == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// Formatting user given data to print
	user_inp := [][]string{
		[]string{"Target Server", *server + ":" + *port},
		[]string{"Send Rate", strconv.Itoa(*rate) + " Queries/Sec"},
		[]string{"Threads", strconv.Itoa(*threads)},
		[]string{"Duration of test", strconv.Itoa(*duration) + " Sec"},
	}
	table1 := tablewriter.NewWriter(os.Stdout)
	for _, ui := range user_inp {
		table1.Append(ui)
	}
	table1.SetAlignment(tablewriter.ALIGN_LEFT)
	table1.Render() // Display info as a table
	fmt.Println(" ")
	fmt.Println("EXECUTING TEST")
	fmt.Println("+-----------------------------------------------------------+")

	limiter := time.Tick(time.Second) // Ticker used for rate limiting packets per second
	res := make(chan Results, *threads)

	var total_qps int
	var total_rtt time.Duration
	// Channel which gets a value when the timer specified with duration has expired
	// Used to end the main forloop in the send_qry function
	ender := time.Tick(time.Duration(*duration) * time.Second)
	// Create as many goroutines as specified by "-t" argument
	for i := 1; i <= *threads; i++ {
		go send_qry(*server, *rate, *port, *duration, *threads, limiter, res, ender)
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
	tabledata := [][]string{
		[]string{term.Yellowf(strconv.Itoa(total_qps)), term.Yellowf(total_rtt.String())},
	}
	table2 := tablewriter.NewWriter(os.Stdout)
	table2.SetHeader([]string{"Average Queries/Sec", "Average Latency"})
	for _, dat := range tabledata {
		table2.Append(dat)
	}

	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println(" ")
	table2.Render()

	os.Exit(0)
}
