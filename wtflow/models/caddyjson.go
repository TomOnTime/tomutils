package models

import (
	"encoding/json"
	"log"
	"math"
	"strings"
	"time"
)

type CaddyLogJSON struct {
	Level   string  `json:"level"`
	Ts      float64 `json:"ts"`
	Logger  string  `json:"logger"`
	Msg     string  `json:"msg"`
	Request struct {
		RemoteIP   string `json:"remote_ip"`
		RemotePort string `json:"remote_port"`
		Proto      string `json:"proto"`
		Method     string `json:"method"`
		Host       string `json:"host"`
		URI        string `json:"uri"`
		Headers    struct {
			SecChUaPlatform []string `json:"Sec-Ch-Ua-Platform"`
			Accept          []string `json:"Accept"`
			IfNoneMatch     []string `json:"If-None-Match"`
			IfModifiedSince []string `json:"If-Modified-Since"`
			SecChUaMobile   []string `json:"Sec-Ch-Ua-Mobile"`
			UserAgent       []string `json:"User-Agent"`
			SecFetchSite    []string `json:"Sec-Fetch-Site"`
			SecFetchMode    []string `json:"Sec-Fetch-Mode"`
			SecFetchDest    []string `json:"Sec-Fetch-Dest"`
			Referer         []string `json:"Referer"`
			AcceptEncoding  []string `json:"Accept-Encoding"`
			AcceptLanguage  []string `json:"Accept-Language"`
			SecChUa         []string `json:"Sec-Ch-Ua"`
		} `json:"headers"`
		TLS struct {
			Resumed     bool   `json:"resumed"`
			Version     int    `json:"version"`
			CipherSuite int    `json:"cipher_suite"`
			Proto       string `json:"proto"`
			ServerName  string `json:"server_name"`
		} `json:"tls"`
	} `json:"request"`
	UserID      string  `json:"user_id"`
	Duration    float64 `json:"duration"`
	Size        int     `json:"size"`
	Status      int     `json:"status"`
	RespHeaders struct {
		Server []string `json:"Server"`
		Etag   []string `json:"Etag"`
	} `json:"resp_headers"`
}

func (db *FlowDb) AddFlowFromCaddyJSON(line string) {
	//fmt.Println(line)
	var data CaddyLogJSON
	err := json.Unmarshal([]byte(line), &data)
	if err != nil {
		log.Printf("BAD JSON: %v: %v\n", err, line)
		return
	}

	// Convert the timestamp:
	integ, decim := math.Modf(data.Ts)
	ts := time.Unix(int64(integ), int64(decim*(1e9)))

	flow := &Flow{
		Time:    ts,
		Path:    data.Request.URI,
		Referer: strings.Join(data.Request.Headers.Referer, "\t"),
	}

	// Guess a unique user:
	var usertoken string // opaque string that changes with user.
	if data.UserID == "" {
		usertoken = data.Request.RemoteIP
	} else {
		usertoken = data.UserID + "@" + data.Request.RemoteIP
	}

	// Add it to the database.

	fs, ok := db.LookupUser(usertoken)
	if !ok {
		// User not found?  Add it.
		fs = db.AddUser(usertoken, &Flows{HostIP: data.Request.RemoteIP})
	} else {
		if len(fs.Items) > 0 {
			// Previous item very old?  Start a new "user".
			prevtime := fs.Items[len(fs.Items)-1].Time
			if flow.Time.Sub(prevtime) > MaxUserPause {
				fs = db.AddUser(usertoken, &Flows{HostIP: data.Request.RemoteIP})
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
