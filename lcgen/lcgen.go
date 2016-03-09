package main

/*
  go run *.go [--debug] command
	  command
		a or announce:  the initial announcement
		r or reminder:  the reminder the friday before
		t or today":  morning of the broadcast
		w or watch:  after the broadcast is released to YouTube

		--debug: Show any "DELETE" line.
*/

import (
	"flag"
	"fmt"
	"strings"
)

var announce = `

******* lc_00EPISODE_NUM_announce

Next on LISA Conversations: GUEST_NAME on "ORIG_TITLE"

<<<<<<<<<< lc_00EPISODE_NUM_announce
Our next guest will be GUEST_NAME. We'll be discussing GUEST_PRONOUN talk from ORIG_CONFNAME titled <i>ORIG_TITLE</i>.

<b>Watch live!</b> We'll be recording the episode on EPISODE_DATE at EPISODE_TIME.
Particpate in the live Q&amp;A by submitting your questions during the broadcast.
Pre-registration is recommended but not required.  Register and/or watch via
[this link](ONAIR_LINK).

* Homework: Watch GUEST_PRONOUN talk ahead of time.
  * <i>ORIG_TITLE</i>
  * Recorded at ORIG_CONFNAME
  * [Talk Description](ORIG_LISTING)
  * [Video and Slides](ORIG_SLIDES)
  * [YouTube](ORIG_YOUTUBE)

* <b>Watch live!</b>
  * LISA Conversations Episode #EPISODE_NUM
  * Co-hosts: Lee Damon and Thomas Limoncelli
  * Guest: GUEST_NAME
  * Will be recorded: EPISODE_DATE at EPISODE_TIME ([convert](EPISODE_CONVERT))

The recorded episode will be available shortly afterwards on YouTube.

You won't want to miss this!

`

var reminder = `

******* lc_00EPISODE_NUM_reminder

Reminder: Do your homework for next week's LISA Conversations: GUEST_NAME on "ORIG_TITLE"

<<<<<<<<<< lc_00EPISODE_NUM_announce
This weekend is a good time to watch the video we'll be discussing
on the next episode of LISA conversations:
GUEST_NAME's talk from ORIG_CONFNAME titled <i>ORIG_TITLE</i>.

* Homework: Watch GUEST_PRONOUN talk ahead of time.
  * <i>ORIG_TITLE</i>
  * GUEST_NAME
  * Recorded at ORIG_CONFNAME
  * [Talk Description](ORIG_LISTING)
  * [Video and Slides](ORIG_SLIDES)
  * [YouTube](ORIG_YOUTUBE)

Then you'll be prepared when we record the episode on EPISODE_DATE at EPISODE_TIME ([convert](EPISODE_CONVERT)).  Register (optional) and watch via [this link](ONAIR_LINK).  Watching live makes it possible to participate in the Q&amp;A.

The recorded episode will be available shortly afterwards on YouTube.

You won't want to miss this!
`

var today = `

******* lc_00EPISODE_NUM_today

Watch us live today! [LISA Conversations](https://www.usenix.org/conference/lisa16/lisa-conversations) Episode EPISODE_NUM: GUEST_NAME on "ORIG_TITLE"


<<<<<<<<<< lc_00EPISODE_NUM_today
Today (EPISODE_DATE) we'll be recording episode #EPISODE_NUM of LISA Conversations.  Join the Google Hangout and submit questions live via [this link](ONAIR_LINK).

Our guest will be GUEST_NAME. We'll be discussing GUEST_PRONOUN talk <i>ORIG_TITLE</i> from ORIG_CONFNAME.


* The video we'll be discussing:
  * <i>ORIG_TITLE</i>
  * GUEST_NAME
  * Recorded at ORIG_CONFNAME
  * [Talk Description](ORIG_LISTING)
  * [Video and Slides](ORIG_SLIDES)
  * [YouTube](ORIG_YOUTUBE)

* <b>Watch us record the episode live!</b>
  * EPISODE_DATE at EPISODE_TIME ([convert](EPISODE_CONVERT))
  * LISA Conversations Episode #EPISODE_NUM
  * Co-hosts: Lee Damon and Thomas Limoncelli
  * Guest: GUEST_NAME
  * [Join us](ONAIR_LINK).

The recorded episode will be available shortly afterwards on YouTube.

You won't want to miss this!
`

var watch = `
lc_00EPISODE_NUM_watch

LISA Conversations Episode EPISODE_NUM: GUEST_NAME on "ORIG_TITLE"

<<<<<<<<<< lc_00EPISODE_NUM_watch
Episode EPISODE_NUM of [LISA Conversations](https://www.usenix.org/conference/lisa16/lisa-conversations)
is GUEST_NAME, who presented
<i>ORIG_TITLE</i>
at ORIG_CONFNAME.

* Watch the Episode here:
  * [LISA Conversations Episode #EPISODE_NUM with GUEST_NAME](EPISODE_VIDEO)
	* Co-hosts: Lee Damon and Thomas Limoncelli
	* Guest: GUEST_NAME
  * Recorded EPISODE_DATE

* In this episode we discuss GUEST_PRONOUN talk:
  * <i>ORIG_TITLE</i>
  * Recorded at ORIG_CONFNAME
  * [Talk Description](ORIG_LISTING)
  * [Video and Slides](ORIG_SLIDES)
  * [YouTube](ORIG_YOUTUBE)

You won't want to miss this!
<<<<<<<<<< lc_00EPISODE_NUM_watch

`

var tokens = []struct {
	f string
	t string
}{
	//{"EPISODE_NUM", "7"},
	//{"GUEST_NAME", " Kris Buytaert"},
	//{"GUEST_PRONOUN", "his"},
	//{"ORIG_TITLE", "DevOps: The past and future are here. It's just not evenly distributed (yet)"},
	//{"ORIG_CONFNAME", "LISA '11"},
	//{"ORIG_LISTING", "http://static.usenix.org/legacy/events/lisa11/tech/tech.html#Buytaert"},
	//{"ORIG_SLIDES", "https://www.usenix.org/conference/lisa11/devops-past-and-future-are-here-its-just-not-evenly-distributed-yet"},
	//{"ORIG_YOUTUBE", "https://www.youtube.com/watch?v=p-8UBYMMjp8"},
	//{"EPISODE_YOUTUBE", "https://www.youtube.com/watch?v=IzNzYbLtHLM"},
	//{"EPISODE_DATE", "Tuesday, February 23, 2016"},
	//{"EPISODE_TIME", "11:30 am PST/2:30 pm EST"},

	{"EPISODE_NUM", "8"},
	{"GUEST_NAME", "Caskey Dickson"},
	{"GUEST_PRONOUN", "his"},
	{"ORIG_TITLE", "Why Your Manager LOVES Technical Debt and What to Do About It"},
	{"ORIG_CONFNAME", "LISA '15"},
	{"ORIG_LISTING", "DELETE"},
	{"ORIG_SLIDES", "https://www.usenix.org/conference/lisa15/conference-program/presentation/dickson"},
	{"ORIG_YOUTUBE", "DELETE"},
	{"ONAIR_LINK", "https://plus.google.com/events/cjudags5uialvq78jq67u9vbla8"},
	{"EPISODE_YOUTUBE", "DELETE"},
	{"EPISODE_DATE", "Tuesday, March 29, 2016"},
	{"EPISODE_TIME", "3:30-4:30 p.m. Pacific Time"},
	{"EPISODE_CONVERT", "http://www.timeanddate.com/worldclock/converted.html?iso=20160329T1530&p1=224&p2=179"},

	//{"EPISODE_NUM", "9"},
	//{"GUEST_NAME", "kc claffy"},
	//{"GUEST_PRONOUN", "her"},
	//{"ORIG_TITLE", "Named Data Networking"},
	//{"ORIG_CONFNAME", "LISA '15"},
	//{"ORIG_LISTING", "DELETE"},
	//{"ORIG_SLIDES", "https://www.usenix.org/conference/lisa15/conference-program/presentation/claffy"},
	//{"ORIG_YOUTUBE", "DELETE"},
	//{"ONAIR_LINK", "FILLIN"},
	//{"EPISODE_YOUTUBE", "DELETE"},
	//{"EPISODE_DATE", "Tuesday, April 26, 2016"},
	//{"EPISODE_TIME", "3:30-4:30 p.m. Pacific Time"},
	//{"EPISODE_CONVERT", "http://www.timeanddate.com/worldclock/converted.html?iso=20160426T1530&p1=224&p2=179"},
}

func help() {
	fmt.Println("subcommands: announce, reminder, today, watch")
}

var show_deletes bool

func init() {
	flag.BoolVar(&show_deletes, "debug", false, "show DELETE items")
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		help()
	} else {

		// Select which template:
		var s string
		switch flag.Arg(0)[0:1] {
		case "a": // announce
			s = announce
		case "r": // reminder
			s = reminder
		case "t": // today
			s = today
		case "w": // watch
			s = watch
		default:
			s = ""
			help()
		}

		// Substitute all variables.
		for _, b := range tokens {
			s = strings.Replace(s, b.f, b.t, -1)
		}

		// If a line containst the text DELETE, remove the entire line.
		lines := strings.Split(s, "\n")
		for i := len(lines) - 1; i > 0; i-- {
			if (!show_deletes) && strings.Index(lines[i], "DELETE") != -1 {
				lines = append(lines[:i], lines[i+1:]...)
			}
		}
		s = strings.Join(lines, "\n")

		// Print.
		fmt.Println(strings.TrimSpace(s))
	}

}
