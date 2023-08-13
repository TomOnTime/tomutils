package filminventory

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

func mkSha256(filepath string) string {
	h := sha256.New()

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Cant open %q: %v", filepath, err)
	}

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (f Film) ToCSV() []string {
	var items []string

	items[0] = mkSha256((f.Filename))        // sha256
	items[1] = f.Title                       // Title
	items[2] = f.SourceSite                  // SourceSite
	items[3] = ""                            // ContentTags
	items[4] = ""                            // MetaTags
	items[5] = fmt.Sprintf("%d", f.Duration) // Duration
	items[6] = ""                            // hh
	items[7] = ""                            // PartyScreen
	items[8] = ""                            // URL
	items[9] = ""                            // CreatorName
	items[10] = ""                           // SongTitle
	items[11] = ""                           // OriginalFilename

	return items
}
