package main

/*

The Go JSON module can't not access unexported fields in a struct. So
how do you work with them?

This demonstrates the solution in http://choly.ca/post/go-json-marshalling/
where we have a 2nd struct that embeds the primary struct but adds
fields that will be used to expose the unexported fields.  We then write
MarshalJSON() and UnmarshalJSON() functions that do the right thing.

This also helps in situations where we have fields that require a custom
format only in JSON.
*/

import (
	"encoding/json"
	"fmt"
	"time"
)

// Cranberry stores data.
//  Visible: This field is exported and JSON displays it as usual.
//  invisible: This field is unexported but we want it to be included in JSON.
//  Custom: This field has a custom output format.  We store it as time.Time
//    but when it appears in JSON, it should be in Unix Epoch format.
type Cranberry struct {
	Visible   int       `json:"visible"`
	invisible int       // No tag here
	Custom    time.Time `json:"-"` // Don't output this field (we'll handle it in CranberryJSON).
}

// CranberryAlias is an alias of Cranberry. We use an alias because aliases
// are stripped of any functions and we need a struct without
// MarshalJSON/UnmarshalJSON defined, otherwise we'd get a recursive defintion.
type CranberryAlias Cranberry

// CranberryJSON represents out we represent Cranberry to the JSON package.
type CranberryJSON struct {
	*CranberryAlias       // All the exported fields.
	Invisible       int   `json:"invisible"`
	CustomUnixEpoch int64 `json:"epoch"`
	// FYI: The json tags "invisble" and "epoch" can be any valid JSON tag.
	// It is all a matter of how we want the JSON to be presented externally.
}

// MarshalJSON marshals a Cranberry. (struct to JSON)
func (u *Cranberry) MarshalJSON() ([]byte, error) {
	return json.Marshal(&CranberryJSON{
		CranberryAlias: (*CranberryAlias)(u),
		// Unexported or custom-formatted fields are listed here:
		Invisible:       u.invisible,
		CustomUnixEpoch: u.Custom.Unix(),
	})
}

// UnmarshalJSON unmarshals a Cranberry. (JSON to struct)
func (u *Cranberry) UnmarshalJSON(data []byte) error {
	temp := &CranberryJSON{
		CranberryAlias: (*CranberryAlias)(u),
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Copy the exported fields:
	*u = (Cranberry)(*(temp).CranberryAlias)
	// Each unexported field must be copied and/or converted individually:
	u.invisible = temp.Invisible
	u.Custom = time.Unix(temp.CustomUnixEpoch, 0) // Convert while copying.

	return nil
}

func main() {
	var out []byte
	var err error

	// Demonstration of marshalling: Marshal s (struct) to out ([]byte)
	fmt.Printf("Struct to JSON:\n")
	s := &Cranberry{Visible: 1, invisible: 2, Custom: time.Unix(1521492409, 0)}
	out, err = json.Marshal(s)
	if err != nil {
		panic(err)
	}
	fmt.Printf("      got=%v\n", string(out))
	fmt.Println(` expected={"visible":1,"invisible":2,"epoch":1521492409}`)

	// Demonstration of how to unmarshal: Unmarshal "out" ([]byte) to n (struct)
	fmt.Printf("JSON to struct:\n")
	var n = &Cranberry{}
	err = json.Unmarshal(out, n)
	if err != nil {
		panic(err)
	}
	fmt.Printf("      got=%+v\n", n)
	fmt.Println(` expected=&{Visible:1 invisible:2 Custom:2018-03-19 xx:46:49 xxxxx xxx}`)
}
