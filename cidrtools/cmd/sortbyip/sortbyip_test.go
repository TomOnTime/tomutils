package main

import (
	"strings"
	"testing"

	"github.com/TomOnTime/tomutils/cidrtools/ipimbo"
)

func run(t *testing.T, name string, a, b string) {
	t.Helper()
	ipimboHandle := ipimbo.New()
	var ipdbA ipimbo.ImboList

	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b) + "\n"

	// Read A
	for v := range ipimboHandle.ReadFile(strings.NewReader(a)) {
		ipdbA = append(ipdbA, v)
	}

	// Sort A
	ipdbA.Sort()

	result := ""
	for _, v := range ipdbA {
		result += v.String() + "\n"
	}

	if result != b {
		t.Errorf("test #%s failed:\nINPUT:\n%s\n\nEXPECTED:\n%s\nGOT:\n%s\n", name, a, b, result)
	}

}

func TestSort1(t *testing.T) {

	run(t, "simple4", `
0.1.2.1
1.2.3.4
30.1.1.1
4.1.1.1
`,
		`
0.1.2.1
1.2.3.4
4.1.1.1
30.1.1.1
`)

	run(t, "simple6", `
ff07:f2f8:a9c0::30
ff07:f2f8:a9c0::4
FF07:f2f8:a9c0::3
1007:f2f8:a9c0::3
0107:f2f8:a9c0::3
fe80::5054:ff:fef3:3410
2607:f2f8:a9c0::3
::1
`,
		`
::1
107:f2f8:a9c0::3
1007:f2f8:a9c0::3
2607:f2f8:a9c0::3
fe80::5054:ff:fef3:3410
ff07:f2f8:a9c0::3
ff07:f2f8:a9c0::4
ff07:f2f8:a9c0::30
`)

	run(t, "mixed", `
0107:f2f8:a9c0::3
fe80::5054:ff:fef3:3410
2607:f2f8:a9c0::3
1.2.3.4
4.1.1.1
30.1.1.1
::1
255.1.1.1
`, `
1.2.3.4
4.1.1.1
30.1.1.1
255.1.1.1
::1
107:f2f8:a9c0::3
2607:f2f8:a9c0::3
fe80::5054:ff:fef3:3410
`)

	run(t, "cidrs", `
0107:f2f8:a9c0::3/24
0107:f2f8:a9c0::3/56
0107:f2f8::/56
10.10.10.0/24
10.10.10.0/32
10.10.10.10/24
10.10.10.10/32
`, `
10.10.10.0/24
10.10.10.0/32
10.10.10.10/24
10.10.10.10/32
107:f2f8::/56
107:f2f8:a9c0::3/24
107:f2f8:a9c0::3/56
`)

}
