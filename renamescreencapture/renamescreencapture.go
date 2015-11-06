package main

/*

renamescreencapture.go -- Rename OS X ScreenCapture files

El Capitan (MacOS X 10.11) removed the ability to select the exact
format of screencapture's filenames.  This utility renames such
files to be more unix-friendly: No spaces, 24-hour time.

Flags:

  -y  Actually rename the files. Otherwise it just reports.

Example:
  $ renamescreencapture *
  Renaming screencapture 2015-11-19 at 1.44.04 AM.png screencapture_2015-11-19_01.44.04.png
  ...test mode...
  Renaming screencapture 2015-11-03 at 10.54.57 PM.png screencapture_2015-11-03_22.54.57.png
  ...test mode...

  $ renamescreencapture -y *
  Renaming screencapture 2015-11-19 at 1.44.04 AM.png screencapture_2015-11-19_01.44.04.png
  Renaming screencapture 2015-11-03 at 10.54.57 PM.png screencapture_2015-11-03_22.54.57.png
*/

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var activate bool

const scapFilePrefix = "screencapture "
const scapFilePrefixNew = "screencapture_"

func main() {
	flag.BoolVar(&activate, "y", false, "Actually rename the files")
	flag.Parse()

	for i, arg := range flag.Args() {

		dir, file := filepath.Split(arg)
		ext := filepath.Ext(file)

		if !strings.HasPrefix(file, scapFilePrefix) {
			fmt.Printf("Skipping %v: %s (not a screencapture file)\n", i, arg)
			continue
		}

		datePart := file[len(scapFilePrefix) : len(file)-len(ext)]

		t, err := time.Parse("2006-01-02 at 3.04.05 PM", datePart)
		if err != nil {
			fmt.Printf("Skipping %v: %s (no date found)\n", i, arg)
			continue
		}

		newDatePart := t.Format("2006-01-02_15.04.05")
		newFile := filepath.Join(dir, scapFilePrefixNew+newDatePart+ext)

		fmt.Printf("Renaming %v %v\n", arg, newFile)
		if activate {
			os.Rename(arg, newFile)
		} else {
			fmt.Println("...test mode...")
		}
	}
}
