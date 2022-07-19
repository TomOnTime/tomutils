package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/TomOnTime/tomutils/vidnamer/filehash"
	"github.com/TomOnTime/tomutils/vidnamer/filminventory"
	"golang.org/x/exp/slices"
)

/*

Operations:

Find new MD5 hashs, draft YAML entries.

Find filenames that are wrong based on MD5 hash and whats in YAML.

NOT IMPLEMENTED: Report YAML entries to be deleted:
- Read YAML
- Read directory
- If any filenames are missing, output "delete

*/

func main() {
	var err error

	// Read the inventory.yaml file
	inventory, err := filminventory.FromYamlfile("../inventory.yaml")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("INVENYAML:\n%+v\n\n", inventory)

	// Read the md5.txt file
	md5db, err := filehash.FromFile("../md5.list")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("MD5DB:\n%v\n\n", md5db)

	// Hydrate inventory with md5hash info:

	// For each item missing from the inventory, output YAML:
	missingFilenames, missingSigs := missingFromInventory(inventory, md5db)
	_ = missingSigs
	if len(missingFilenames) > 0 {
		f, err := os.OpenFile("../missing.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(f, "---\n")
		for i, filename := range missingFilenames {
			_ = i
			// Parse the filename
			inventoryItem := filminventory.ParseFilename(filename)
			inventoryItem.Signature = missingSigs[i]
			// Output template yaml
			y, err := inventoryItem.ToYaml()
			if err != nil {
				panic(err)
			}
			fmt.Fprint(f, string(y))
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}

	// Audit keywords
	// Read the list of permitted keywords:
	permittedKeywords, err := readLines("../keywords.txt")
	if err != nil {
		log.Fatal(err)
	}
	invalidKeywords := auditKeywords(inventory, permittedKeywords)
	if len(invalidKeywords) > 0 {
		fmt.Printf("INVALID KEYWORDS: %s\n", strings.Join(invalidKeywords, " "))
	}

	// Audit tags
	permittedTags, err := readLines("../tags.txt")
	if err != nil {
		log.Fatal(err)
	}
	invalidTags := auditTags(inventory, permittedTags)
	if len(invalidTags) > 0 {
		fmt.Printf("INVALID TAGS: %s\n", strings.Join(invalidTags, " "))
	}

	// For each item in the inventory, if the filename is wrong, output
	// a correction.
	f, err := os.OpenFile("../fixit.sh", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	var newinventory []filminventory.Film
	for _, invItem := range inventory {
		//fmt.Printf("DEBUG: name=%s\n", invItem.Title)
		existing := filminventory.ExistingFilename(invItem.Signature, md5db)
		if existing == "" {
			continue
		}
		desired := invItem.DesiredFilename()
		// If the file name is wrong, output rename command.
		if desired != existing {
			fmt.Fprintln(f, makeRenameCmd(existing, desired))
		}
		newinventory = append(newinventory, invItem)
	}
	inventory = newinventory
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	// output NEW YAML
	err = writeNewInventoryYAML(inventory, "../inventory.yaml.NEW")
	if err != nil {
		log.Fatal(err)
	}

}

func writeNewInventoryYAML(inventory []filminventory.Film, filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "---\n")
	for _, inventoryItem := range inventory {
		y, err := inventoryItem.ToYaml()
		if err != nil {
			return err
		}
		fmt.Fprint(f, string(y))
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
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
	return fmt.Sprintf("mv -i \\\n   %q\\\n   %q", wrong, right)
}
