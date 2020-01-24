package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfreiberg/gotwilio"
)

func main() {

	accountSid, ok := os.LookupEnv("TWILIO_SID")
	if !ok {
		fmt.Printf("Env variable TWILIO_SID not set. Exiting.\n")
		os.Exit(1)
	}

	authToken, ok := os.LookupEnv("TWILIO_KEY")
	if !ok {
		fmt.Printf("Env variable TWILIO_KEY not set. Exiting.\n")
		os.Exit(1)
	}

	dryRun := flag.Bool("n", false, "dry run mode")
	flag.Parse()

	if len(flag.Args()) != 3 {
		fmt.Printf("Three args expected: from to phone# and message.\n")
		os.Exit(1)
	}
	from := flag.Args()[0]
	to := flag.Args()[1]
	message := flag.Args()[2]

	twilio := gotwilio.NewTwilioClient(accountSid, authToken)
	fmt.Printf("twilio.SendSMS(%q, %q, %q)\n", from, to, message)
	if !*dryRun {
		a, b, c := twilio.SendSMS(from, to, message, "", "")
		fmt.Printf("RESULT: %v %v %v\n", a, b, c)
	}
}
