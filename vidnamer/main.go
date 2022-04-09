package main

import (
	"fmt"
	"log"

	"github.com/TomOnTime/tomutils/vidnamer/filehash"
	"github.com/TomOnTime/tomutils/vidnamer/filminventory"
	"golang.org/x/exp/slices"
)

func main() {
	var err error

	// Read the inventory.yaml file
	inventory, err := filminventory.FromYamlfile("inventory.yaml")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("INVENYAML:\n%+v\n\n", inventory)
	// Read the md5.txt file
	md5db, err := filehash.FromFile("md5.list")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("MD5DB:\n%v\n\n", md5db)

	// Hydrate inventory with md5hash info:

	// For each item missing from the inventory:
	missingFilenames, missingSigs := missingFromInventory(inventory, md5db)
	_ = missingSigs
	fmt.Printf("---\n")
	for i, filename := range missingFilenames {
		_ = i
		// Parse the filename
		//fmt.Printf("DEBUG: FILENAME: %v\n", filename)
		inventoryItem := filminventory.ParseFilename(filename)
		inventoryItem.Signature = missingSigs[i]
		// Output template yaml
		y, err := inventoryItem.ToYaml()
		if err != nil {
			panic(err)
		}
		fmt.Printf(string(y))
	}

	// // For each item in the inventory...
	for _, invItem := range inventory {

		existing := filminventory.ExistingFilename(invItem.Signature, md5db)
		desired := invItem.DesiredFilename()

		// If the file name is wrong, output rename command.
		if desired != existing {
			fmt.Println(makeRenameCmd(existing, desired))
		} else {
			fmt.Printf("echo %q\n", existing)
		}
	}

}

func missingFromInventory(inv []filminventory.Film, hashes []filehash.Info) ([]string, []string) {
	// invSigs  := list of hashes in inv.
	var invSigs []string
	for i := range inv {
		n := inv[i].Signature
		if n != "" {
			invSigs = append(invSigs, n)
		}
	}
	//sort.Strings(invSigs)
	//fmt.Printf("DEBUG: invSigs: %v\n", invSigs)

	// hasSigs := list of hashes in hashes.
	// var hasSigs []string
	// for i := range hashes {
	// 	n := hashes[i].Signature
	// 	if n != "" {
	// 		hasSigs = append(hasSigs, n)
	// 	}
	// }
	//sort.Strings(hasSigs)
	//fmt.Printf("DEBUG: hasSigs: %v\n", hasSigs)

	// return items in hashes that are not in inv.
	var r []string
	var s []string
	for _, h := range hashes {
		if !slices.Contains(invSigs, h.Signature) {
			r = append(r, h.Filename)
			s = append(s, h.Signature)
		}
	}
	//fmt.Printf("DEBUG: missing: %v\n", r)
	return r, s
}

func makeRenameCmd(wrong, right string) string {
	return fmt.Sprintf("mv \\\n   %q\\\n   %q", wrong, right)
}
