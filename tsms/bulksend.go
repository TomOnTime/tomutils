package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
		fmt.Printf("Three args expected: fromNumber toFile messageFile.\n")
		os.Exit(1)
	}
	from := flag.Args()[0]
	toFilename := flag.Args()[1]
	messageFilename := flag.Args()[2]

	toData, err := ioutil.ReadFile(toFilename)
	if err != nil {
		fmt.Printf("Can not read %q: %v\n", toFilename, err)
		os.Exit(1)
	}

	messageRaw, err := ioutil.ReadFile(messageFilename)
	if err != nil {
		fmt.Printf("Can not read %q: %v\n", toFilename, err)
		os.Exit(1)
	}
	message := strings.TrimSpace(string(messageRaw))

	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	toItems := strings.Split(string(toData), "\n")

	for _, to := range toItems {
		if to == "" {
			continue
		}
		fmt.Printf("twilio.SendSMS(%q, %q, %q)\n", from, to, message)
		if !*dryRun {
			resp, excp, err := twilio.SendSMS(from, to, message, "", "")
			if err == nil {
				fmt.Printf("SUCCESS!\n")
				fmt.Printf("RESPONSE: %v\n", resp)
				fmt.Printf("EXCEPTION: %v\n", excp)
			} else {
				fmt.Printf("FAILURE: %v\n", err)
				fmt.Printf("RESPONSE: %v\n", resp)
				fmt.Printf("EXCEPTION: %v\n", excp)
			}
		}
	}

}
