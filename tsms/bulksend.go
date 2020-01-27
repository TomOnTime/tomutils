package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/TomOnTime/tomutils/tsms/basictwilio"
)

func main() {

	dryRun := flag.Bool("n", false, "dry run mode")
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Printf("Three args expected: toFile messageFile.\n")
		os.Exit(1)
	}
	toFilename := flag.Args()[0]
	messageFilename := flag.Args()[1]

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

	toItems := strings.Split(string(toData), "\n")
	for _, to := range toItems {
		to = strings.TrimSpace(to)
		if to == "" {
			continue
		}
		basictwilio.SendSMS(to, message, *dryRun)
	}

}
