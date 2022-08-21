package models

// TODO:
// When adding a Flow, check to see if the previous Flow or this
// usertoken is very distant. In that case, start a new Flows.

//	`
//	1.2.3.4 foo.comcast.com
//  15:12:11.0     /foo ->
//  5s               /bar
//  5m             /bar
//  16:22:11.0       /bim
//`

// https://golang.org/pkg/net/url/#URL
// https://github.com/Songmu/axslogparser/blob/master/axslogparser_test.go

import (
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Songmu/axslogparser"
)

// var MaxUserPause = time.Minute * 40
var MaxUserPause = time.Hour * 1

type Flows struct {
	HostIP string
	Items  []*Flow
}

type Flow struct {
	Time    time.Time
	Path    string
	Referer string
}

type FlowDb struct {
	items          map[string]*Flows
	ReportItems    []*Flows
	IgnorePrefixes []string
	IgnoreSuffixes []string
}

func (db *FlowDb) SkipPath(path string) bool {

	pathlc := strings.ToLower(path)
	for _, p := range db.IgnoreSuffixes {
		if strings.HasSuffix(pathlc, p) {
			return true
		}
	}

	for _, p := range db.IgnorePrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}

	return false
}

func (db *FlowDb) LookupUser(u string) (*Flows, bool) {
	if db.items == nil {
		db.items = map[string]*Flows{}
	}
	if fs, ok := db.items[u]; !ok {
		return nil, false
	} else {
		//fmt.Printf("DEBUG: LookupUser(%q) = %p\n", u, fs)
		return fs, true
	}
}

func (db *FlowDb) AddUser(u string, fs *Flows) *Flows {
	db.items[u] = fs
	db.ReportItems = append(db.ReportItems, fs)
	return fs
}

func (db *FlowDb) AddFlowFromString(line string) {
	//fmt.Println(line)
	logentry, err := axslogparser.Parse(line)
	if err != nil {
		log.Printf("BAD LINE: %v: %v\n", err, line)
		return
	}
	u, err := url.Parse(logentry.RequestURI)
	if err != nil {
		//log.Printf("BAD URI: %v: %v\n", err, line)
		return
	}
	//host := logentry.Host
	//path := u.Path
	path := u.String()
	if db.SkipPath(path) {
		return
	}
	var usertoken string // opaque string that changes with user.
	if logentry.User == "-" {
		usertoken = logentry.Host
	} else {
		usertoken = logentry.User + "@" + logentry.Host
	}

	//fmt.Printf("%v Host=%v Path=%v\n", usertoken, host, path)

	flow := &Flow{
		Time:    logentry.Time,
		Path:    path,
		Referer: logentry.Referer,
	}

	// Add it to the database.

	fs, ok := db.LookupUser(usertoken)
	if !ok {
		// User not found?  Add it.
		fs = db.AddUser(usertoken, &Flows{HostIP: logentry.Host})
	} else {
		if len(fs.Items) > 0 {
			// Previous item very old?  Start a new "user".
			prevtime := fs.Items[len(fs.Items)-1].Time
			if flow.Time.Sub(prevtime) > MaxUserPause {
				fs = db.AddUser(usertoken, &Flows{HostIP: logentry.Host})
			}
		}
	}

	// If the previous result was a redirect to this item because
	// of a missing trailing "/", delete it so we only store the
	// destination.
	if len(fs.Items) > 0 {
		topindex := len(fs.Items) - 1
		if flow.Path == (fs.Items[topindex].Path + "/") {
			fs.Items = fs.Items[:len(fs.Items)-1]
		}

	}
	fs.Items = append(fs.Items, flow)
}
