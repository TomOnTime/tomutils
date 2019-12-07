package main

//	`
//	1.2.3.4 foo.comcast.com
//  15:12:11.0     /foo ->
//  5s               /bar
//  5m             /bar
//  16:22:11.0       /bim
//`

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/TomOnTime/tomutils/wtflow/models"
	"github.com/rs/dnscache"
)

var resolver = &dnscache.Resolver{}

func report(db *models.FlowDb, domain string) {
	fmt.Printf("Flow Report:\n")
	fmt.Printf("Item count: %v\n", len(db.ReportItems))
	prevtime := time.Time{}

	var timestr string

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	for _, fs := range db.ReportItems {
		var hostname string
		var ok bool

		if hostname, ok = stupidcache[fs.HostIP]; !ok {
			//fmt.Printf("CACHE MISS: %v %v\n", fs.HostIP, hostname)
			addrs, err := resolver.LookupAddr(context.Background(), fs.HostIP)
			if err != nil || len(addrs) == 0 {
				hostname = fs.HostIP
			} else {
				hostname = addrs[0]
			}
		}
		stupidcache[fs.HostIP] = hostname
		if strings.HasPrefix(hostname, "crawl") ||
			strings.HasSuffix(hostname, ".crawl.yahoo.net.") ||
			strings.HasSuffix(hostname, ".spider.yandex.com.") ||
			strings.HasSuffix(hostname, ".a.ahrefs.com.") ||
			strings.HasSuffix(hostname, ".search.qwant.com.") ||
			strings.HasSuffix(hostname, ".search.msn.com.") {
			continue
		}

		fmt.Println()

		if len(fs.Items) == 1 {
			fmt.Printf("%v %v (1 item):\n", fs.HostIP, hostname)
		} else {
			fmt.Printf("%v %v (%v items):\n", fs.HostIP, hostname, len(fs.Items))
		}

		//for a, b := range fs.Items {
		//fmt.Printf("%v %v\n", a, b.Referer)
		//}

		for i, f := range fs.Items {
			var asterix string
			var display_referer string

			// The first item is always a date, the others might be
			// displayed as durations.
			//timestr = f.Time.Format("2006-01-02 15:04:05")
			timestr = f.Time.In(location).Format("2006-01-02 15:04:05")

			if i > 0 {
				// The 2nd line on we have more to customize.
				// We can reduce the timestamp to a direction:
				d := f.Time.Sub(prevtime)
				//fmt.Printf("DELTA: %v  --- %v\n", d.String(), d)
				if d.Hours() < 1.5 {
					timestr = d.String()
				}
			}

			referer := f.Referer
			if referer == "-" || referer == "" {
				referer = ""
				display_referer = ""
			} else {
				display_referer = "REF=" + referer

				for j, x := range fs.Items[0:i] {
					oldurl := `https://` + domain + x.Path
					//fmt.Printf("DEBUG: compare %v,%v %q %q\n", i, j, referer, oldurl)
					if referer == oldurl {
						asterix = "*"
						if j-i == -1 {
							display_referer = ""
						} else {
							display_referer = fmt.Sprintf("REF=%v", j-i)
							//display_referer = fmt.Sprintf("REF=%v D=%v", j-i, referer)
						}
						break
					}
				}
			}

			url := `https://` + domain + f.Path
			fmt.Printf("      %20s %1s %v     %v\n",
				timestr,
				asterix,
				url,
				display_referer,
			)
			prevtime = f.Time
		}
	}
}

func main() {
	db := &models.FlowDb{
		IgnorePrefixes: []string{
			`/css/`,
			`/favicon`,
			`/fonts/`,
			`/images/`,
			`/js/`,
			// Specific files:
			`/index.json`,
			`/index.xml`,
			`/robots.txt`,
			`/sitemap.xml`,
			`/styles.css`,
		},
		IgnoreSuffixes: []string{
			`.png`,
			`.jpg`,
		},
	}

	stupidBegin(CACHEFILE)

	if len(os.Args) < 2 {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			//log.Println("line", s.Text())
			db.AddFlowFromString(s.Text())
		}
	} else {
		for _, fn := range os.Args[1:] {
			f, err := os.Open(fn)
			if err != nil {
				log.Printf("ERROR: %v\n", err)
				os.Exit(1)
			}
			s := bufio.NewScanner(f)
			for s.Scan() {
				// log.Println("line", s.Text())
				db.AddFlowFromString(s.Text())
			}
		}
	}

	report(db, os.Getenv("WTFLOW_DOMAIN"))
	//fmt.Printf("CACHE: %+v\n", stupidcache)

	stupidEnd(CACHEFILE)

}
