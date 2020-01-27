package basictwilio

import (
	"fmt"
	"os"

	"github.com/nyaruka/phonenumbers"
	"github.com/sfreiberg/gotwilio"
)

func SendSMS(toNumber, message string, dryRun bool) error {

	accountSid, ok := os.LookupEnv("TWILIO_ACCOUNT_SID")
	if !ok {
		fmt.Printf("Env variable TWILIO_ACCOUNT_SID not set. Exiting.\n")
		os.Exit(1)
	}

	authToken, ok := os.LookupEnv("TWILIO_AUTH_TOKEN")
	if !ok {
		fmt.Printf("Env variable TWILIO_AUTH_TOKEN not set. Exiting.\n")
		os.Exit(1)
	}

	fromNumber, ok := os.LookupEnv("TWILIO_DEF_FROM")
	if !ok {
		fmt.Printf("Env variable TWILIO_DEF_FROM not set. Exiting.\n")
		os.Exit(1)
	}

	num, err := phonenumbers.Parse(toNumber, "US")
	if err != nil {
		return err
	}
	toNumber = fmt.Sprintf("+%d%d", num.GetCountryCode(), num.GetNationalNumber())

	shortMsg := message[0:20]
	if len(message) > 20 {
		shortMsg = shortMsg + "..."
	}

	twilio := gotwilio.NewTwilioClient(accountSid, authToken)
	if dryRun {
		fmt.Printf("twilio.SendSMS(%q, %q, %q)\n", fromNumber, toNumber, shortMsg)
	} else {
		resp, excp, err := twilio.SendSMS(fromNumber, toNumber, message, "", "")
		if err == nil {
			fmt.Printf("SUCCESS! %v\n", toNumber)
			fmt.Printf("    RESPONSE: %v\n", resp)
		} else {
			fmt.Printf("FAILURE: %v\n", toNumber)
			fmt.Printf("    ERROR: %v\n", err)
			fmt.Printf("    RESPONSE: %v\n", resp)
			fmt.Printf("    EXCEPTION: %v\n", excp)
		}
		fmt.Println()
	}

	return err
}
