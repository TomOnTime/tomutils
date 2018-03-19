package main

import (
	"encoding/json"
	"fmt"
)

// Apple stores data.
type Apple struct {
	Visible   int `json:"visible"`
	Invisible int `json:"invisible"`
}

// Banana stores data.
type Banana struct {
	Visible   int `json:"visible"`
	invisible int // Tag on unexported field wouldn't make sense `json:"invisible"`
}

// Test implements the solution in https://stackoverflow.com/questions/28035784
type Test struct {
	fieldA string
	FieldB int
	FieldC string
}

// TestJSON implements the solution in https://stackoverflow.com/questions/28035784
type TestJSON struct {
	FieldA string `json:"fieldA"`
	FieldB int    `json:"fieldB"`
	FieldC string `json:"fieldC"`
}

// MarshalJSON implements the solution in https://stackoverflow.com/questions/28035784
func (t *Test) MarshalJSON() ([]byte, error) {
	return json.Marshal(TestJSON{
		t.fieldA,
		t.FieldB,
		t.FieldC,
	})
}

// UnmarshalJSON implements the solution in https://stackoverflow.com/questions/28035784
func (t *Test) UnmarshalJSON(b []byte) error {
	temp := &TestJSON{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	t.fieldA = temp.FieldA
	t.FieldB = temp.FieldB
	t.FieldC = temp.FieldC

	return nil
}

// Cranberry implements the solution in http://choly.ca/post/go-json-marshalling/
type Cranberry struct {
	Visible   int `json:"visible"`
	invisible int // Tag on unexported field wouldn't make sense `json:"invisible"`
}

// MarshalJSON implements the solution in http://choly.ca/post/go-json-marshalling/
func (u *Cranberry) MarshalJSON() ([]byte, error) {
	type Alias Cranberry
	return json.Marshal(&struct {
		Invisible int `json:"invisible"`
		*Alias
	}{
		Invisible: u.invisible,
		Alias:     (*Alias)(u),
	})
}

// UnmarshalJSON implements the solution in http://choly.ca/post/go-json-marshalling/
func (u *Cranberry) UnmarshalJSON(data []byte) error {
	type Alias Cranberry
	// Each unexported field must be listed here twice, once to define the struct, then to fill it.
	temp := &struct {
		Invisible int `json:"invisible"`
		*Alias
	}{
		Invisible: u.invisible,
		Alias:     (*Alias)(u),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*u = (Cranberry)(*(temp).Alias) // Copy the exported fields.
	// Each unexported field must be copied individually:
	u.invisible = temp.Invisible

	return nil
}

func main() {

	var out []byte
	var err error

	// Everything exported and appears in JSON.
	{
		fmt.Printf("Apple:\n")
		s := &Apple{Visible: 1, Invisible: 2}
		out, err = json.Marshal(s)
		if err != nil {
			panic(err)
		}
		fmt.Printf(" got=%v\n", string(out))
		fmt.Println(` exp={"visible":1,"invisible":2}`)
	}

	// One field unexported.  Should not appear in JSON.
	{
		fmt.Printf("Banana:\n")
		s := &Banana{Visible: 1, invisible: 2}
		out, err = json.Marshal(s)
		if err != nil {
			panic(err)
		}
		fmt.Printf(" got=%v\n", string(out))
		fmt.Println(` exp={"visible":1}`)
	}

	// fieldA not exported but should appear in output.
	{
		fmt.Printf("Test:\n")
		s := &Test{fieldA: "a", FieldB: 2, FieldC: "3"}
		out, err = json.Marshal(s)
		if err != nil {
			panic(err)
		}
		fmt.Printf(" got=%v\n", string(out))
		fmt.Println(` exp={"fieldA":"a","fieldB":2,"fieldC":"3"}`)
	}

	// one field not exported but should appear in output.
	{
		fmt.Printf("Cranberry:\n")
		s := &Cranberry{Visible: 1, invisible: 2}
		out, err = json.Marshal(s)
		if err != nil {
			panic(err)
		}
		fmt.Printf(" got=%v\n", string(out))
		fmt.Println(` exp={"invisible":2,"visible":1}`)
		// Try unmarshalling too:
		var n = &Cranberry{}
		err = json.Unmarshal(out, n)
		if err != nil {
			panic(err)
		}
		fmt.Printf(" rev=%+v\n", n)
	}

}
