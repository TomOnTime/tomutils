package main

// Generate .in-addr.arpa and ip6.arpa filenames from a list of CIDR blocks.

import (
	"bufio"
	"fmt"
	"html/template"
	"net"
	"os"
	"strings"
)

func MakeArpaZone(n *net.IPNet) string {
	ones, bits := n.Mask.Size()
	var sl int
	if bits == 128 { // IPv6
		var h, name string
		for _, b := range n.IP.To16() {
			h = h + fmt.Sprintf("%02x", b)
		}
		for _, b := range h {
			name = string(b) + "." + name
		}
		if (ones % 4) == 0 {
			sl = len(name) - (ones / 4 * 2)
		} else {
			// If this is misaligned, include the extra nibble.
			sl = len(name) - (ones / 4 * 2) - 2
		}
		// TODO(tlim): assert: name[0:sl] are repetitions of "0."
		return name[sl:] + "ip6.arpa"
	} else if bits == 32 { // IPv4
		var name string
		oc := ones
		for _, b := range n.IP.To4() {
			name = fmt.Sprintf("%d.", b) + name
			if oc <= 8 {
				break
			}
			oc -= 8
		}
		return name + "in-addr.arpa"
	} else { // Error
		return "ERROR"
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t := scanner.Text()

		// Skip comments and blanks.
		if len(t) == 0 || t[0] == '#' {
			continue
		}

		// Comment
		fields := strings.SplitN(t, " ", 2)
		var comment string
		if len(fields) == 2 {
			comment = fields[1]
		}

		// address and netlength
		var combined = fields[0]
		_, ipnet, err := net.ParseCIDR(combined)
		if err != nil {
			fmt.Printf("// Skipping: %v (%s)\n", t, err)
			continue
		}

		z := MakeArpaZone(ipnet)

		const mastertemplate = `
  // {{.Address}}{{with .Comment}} ({{.}}){{else}}{{end}}
  zone "{{.Zone}}." {
    type master;
    file "master/{{.Zone}}.zone";
  };
`

		const slavetemplate = `
  // {{.Address}}{{with .Comment}} ({{.}}){{else}}{{end}}
  zone "{{.Zone}}." {
    type slave;
    file "slaves/{{.Zone}}.zone";
    masters { {{.Masters}} };
  };
`

		type ZoneInfo struct {
			Style   string
			Address string
			Comment string
			Zone    string
			Masters string
		}
		zi := ZoneInfo{}
		zi.Address = ipnet.String()
		zi.Comment = comment
		zi.Zone = z
		zi.Masters = "192.111.0.80; 10.8.2.5;"

		// TODO(tlim): This should be a flag.
		//tplate := mastertemplate
		tplate := slavetemplate

		tm := template.Must(template.New("zone").Parse(tplate))
		err = tm.Execute(os.Stdout, zi)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}
