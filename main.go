package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	reqs         int
	max          int
	mailingid    int
	formatcoupon bool
	found        int
	scrapeurl    string
)

func init() {
	flag.IntVar(&reqs, "reqs", 500, "Total requests")                                                                     // How many MailingID vaules you want to check. This scraper only scrapes values higher than this one. If you want old or more coupons use a lower number for the MailingID and increase the total requests.
	flag.IntVar(&max, "concurrent", 45, "Maximum concurrent requests")                                                    // Prevents "Too many open files" error and timeouts.
	flag.IntVar(&mailingid, "mailingid", 27917361605, "Mailing ID to start on")                                           // Get a recent MailingID from an Arby's promo email by clicking to view print version and getting the query value for 'MailingID' from the url. URL Example from 10/25/2017: http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=27917361605
	flag.StringVar(&scrapeurl, "scrapeurl", "http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=", "URL to scape") // This was found by viewing a printable version of an Arby's promo email and removing the extra query values from the URL.
	flag.BoolVar(&formatcoupon, "formatcoupon", true, "Remove the 'Offer valid only at:' text and expiration dates.")     // Remove the "Offer valid only at:" text and expiration dates.
}

type Response struct {
	*http.Response
	err error
}

// Dispatcher
func dispatcher(reqChan chan *http.Request) {
	defer close(reqChan)
	for i := 0; i < reqs; i++ {
		req, err := http.NewRequest("GET", scrapeurl+strconv.Itoa(mailingid+i), nil)
		if err != nil {
			log.Println(err)
		}
		reqChan <- req
	}
}

// Worker Pool
func workerPool(reqChan chan *http.Request, respChan chan Response) {
	t := &http.Transport{}
	for i := 0; i < max; i++ {
		go worker(t, reqChan, respChan)
	}
}

// Worker
func worker(t *http.Transport, reqChan chan *http.Request, respChan chan Response) {
	for req := range reqChan {
		resp, err := t.RoundTrip(req)
		r := Response{resp, err}
		respChan <- r
	}
}

// Formatter
func formatCoupon(couponHTML string) (string, error) {
	reg, err := regexp.Compile("(?i)Offer expires (\\d{1,2}/\\d{1,2}/\\d{2,4})") // Remove expiration dates.
	if err != nil {
		return "", err
	}
	return reg.ReplaceAllString(strings.Replace(couponHTML, "Offer valid only at:", "", -1), ""), nil // Remove "valid only at location" and expiration dates. Store id query value must not be set in url.
}

// Consumer
func consumer(respChan chan Response) int64 {
	var (
		conns      int64
		HTMLString string
	)
	for conns < int64(reqs) {
		select {
		case r, ok := <-respChan:
			if ok {
				if r.err != nil {
					log.Println(r.err)
				} else {
					if r.ContentLength != 117 { // a content length of 117 is the 'not found' error page.
						log.Println(r.Request.URL.String())
						body, err := ioutil.ReadAll(r.Body)
						if formatcoupon {
							HTMLString, err = formatCoupon(string(body))
							if err != nil {
								log.Fatal(err)
							}
							body = []byte(HTMLString)
						}
						// Save the offer to disk
						err = ioutil.WriteFile(r.Request.URL.Query().Get("MailingID")+".html", body, 0644) // MailingID is used as a unique file name.
						if err != nil {
							log.Fatal(err)
						}
						found = found + 1
					}
					if err := r.Body.Close(); err != nil {
						log.Println(r.err)
					}
				}
				conns++
			}
		}
	}
	return conns
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	reqChan := make(chan *http.Request)
	respChan := make(chan Response)
	start := time.Now()
	go dispatcher(reqChan)
	go workerPool(reqChan, respChan)
	conns := consumer(respChan)
	took := time.Since(start)
	ns := took.Nanoseconds()
	av := ns / conns
	average, err := time.ParseDuration(fmt.Sprintf("%d", av) + "ns")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Connections:\t%d\nConcurrent:\t%d\nTotal time:\t%s\nAverage time:\t%s\nTotal results:\t%d\n", conns, max, took, average, found)
}
