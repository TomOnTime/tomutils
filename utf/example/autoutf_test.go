package main_test

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestReader(t *testing.T) {

	expected, err := ioutil.ReadFile("calblur8.htm")
	if err != nil {
		log.Fatal(err)
	}

	// The test files were generate with:
	//    for i in $(iconv  -l|grep -i utf)  ; do
	//        iconv -f UTF-8 -t $i calblur8.htm > calblur8.htm.$i
	//    done
	for _, experiment := range []string{
		"calblur8.htm.UTF-16",
		"calblur8.htm.UTF-16BE",
		"calblur8.htm.UTF-16LE",
		"calblur8.htm.UTF-8",
	} {
		actual, err := main.ReadFileUTF16(experiment)
		if err != nil {
			log.Fatal(err)
		}

		if string(expected) != string(actual) {
			t.Errorf("%v: expected %#v... got %#v...\n", experiment, string(expected)[:4], actual[:4])
		}
	}

}
