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

	for _, fs := range db.ReportItems {
		var hostname string
		var prevurl string

		fmt.Println()

		addrs, err := resolver.LookupAddr(context.Background(), fs.HostIP)
		if err != nil || len(addrs) == 0 {
			hostname = fs.HostIP
		} else {
			hostname = addrs[0]
		}

		if len(fs.Items) == 1 {
			fmt.Printf("%v %v (1 item):\n", fs.HostIP, hostname)
		} else {
			fmt.Printf("%v %v (%v items):\n", fs.HostIP, hostname, len(fs.Items))
		}
		for i, f := range fs.Items {
			var asterix string
			var display_referer string

			if i == 0 {
				timestr = f.Time.Format("2006-01-02 15:04:05")
				asterix = ""
			} else {
				d := f.Time.Sub(prevtime)
				if d.Hours() < 1 {
					timestr = d.String()
				}
				asterix = ""
			}
			prevtime = f.Time

			referer := f.Referer
			if referer == "-" {
				referer = ""
			}
			if referer == "" {
				display_referer = ""
			} else {
				if prevurl == referer {
					asterix = "*"
					display_referer = ""
				} else {
					display_referer = "REF=" + referer
				}
			}

			url := `https://` + domain + f.Path
			prevurl = url
			fmt.Printf("      %20s %1s %v     %v\n",
				timestr,
				asterix,
				url,
				display_referer,
			)
		}
	}
}

func main() {
	db := &models.FlowDb{
		IgnorePrefixes: []string{
			`/css/`,
			`/favicon`,
			`/images/`,
			`/js/`,
			`/styles.css`,
			`/index.json`,
			`/robots.txt`,
			`/sitemap.xml`,
		},
		IgnoreSuffixes: []string{
			`.png`,
			`.jpg`,
		},
	}

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

	report(db, "www.realpornmeets.com")
}
