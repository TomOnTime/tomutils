package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/TomOnTime/tomutils/tsms/basictwilio"
)

func main() {

	dryRun := flag.Bool("n", false, "dry run mode")
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Printf("Two args expected: TO and MESSAGE.\n")
		os.Exit(1)
	}
	to := flag.Args()[0]
	message := flag.Args()[1]

	if !*dryRun {
		basictwilio.SendSMS(to, message, *dryRun)
	}
}
