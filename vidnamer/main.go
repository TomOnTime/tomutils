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
	// Read the md5.txt file
	md5db, err := filehash.FromYamlFile("md5.list")
	if err != nil {
		log.Fatal(err)
	}

	// Hydrate inventory with md5hash info:

	// For each item missing from the inventory:
	missing := missingFromInventory(inventory, md5db)
	for _, filename := range missing {
		// Parse the filename
		inventoryItem := filminventory.ParseFilename(filename)
		// Output template yaml
			y, err := inventoryItem.ToYaml()
			if err != nil {
				panic(err)
			}
		fmt.Printf(string(y))
	}

	// For each item in the inventory...
	for _, invItem := range inventory {
		// Generate the intended filename.
		desired := invItem.DesiredFilename()
		// If the file name is wrong, output rename command.
		existing := invItem.ExistingFilename()
		if desired != existing {
			fmt.Printf(makeRenameCmd(existing, desired))
		}
	}
}

func missingFromInventory(inv []filminventory.Film, hashes []filehash.Info) []string {
	// invSigs  := list of hashes in inv.
	var invSigs []string
	for i := range inv {
		n := inv[i].Signature
		if n != "" {
			invSigs = append(invSigs, n)
		}
	}
	//sort.Strings(invSigs)

	// hasSigs := list of hashes in hashes.
	var hasSigs []string
	for i := range hashes {
		n := hashes[i].Signature
		if n != "" {
			hasSigs = append(hasSigs, n)
		}
	}
	//sort.Strings(hasSigs)

	// return items in hashes that are not in inv.
	var r []string
	for _, sig := range hasSigs {
		if !slices.Contains(invSigs, sig) {
			r = append(r, sig)
		}
	}
	return r
}

func makeRenameCmd(wrong, right string) string {
	return fmt.Sprintf("mv %q %q", wrong, right)
}
