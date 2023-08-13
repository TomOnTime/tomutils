package filminventory

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TomOnTime/tomutils/vidnamer/filehash"
	"golang.org/x/exp/maps"
)

func CsvToFilm(c FilmCSV) (f Film) {

	// Fix up the SourceSite
	sourceSite := c.SourceSite
	if sourceSite == "" {
		if strings.Contains(c.Filename, "SpankBang") {
			sourceSite = "SpankBang"
		} else if strings.Contains(c.Filename, "xvideos") {
			sourceSite = "xvideos"
		}
	}
	if strings.Contains(strings.ToLower(sourceSite), "spankbang") {
		sourceSite = "SpankBang"
	}
	if strings.Contains(strings.ToLower(sourceSite), "xvideos") {
		sourceSite = "xvideos"
	}

	// Fix up titles
	title := c.Title
	title = strings.TrimPrefix(title, "SpankBang.com_")
	title = strings.TrimPrefix(title, "xvideos.com_")

	f.Signature = c.Sha256
	f.Filename = c.Filename
	f.OrderPrefix = c.OrderPrefix
	f.FileExt = c.FileExt
	f.URL = c.URL
	f.SourceSite = sourceSite
	f.Title = title
	f.Author = c.Sha256
	f.Keywords = strings.Split(c.ContentTags, ",")
	f.Tags = mkTags(c.MetaTags)
	f.Duration = parseIntLazy(c.Duration)
	f.Hh = parseIntLazy(c.Hh)
	f.Room = c.PartyScreen
	//f.Test = c.Test
	f.SongTitle = c.SongTitle
	f.OrigFName = c.DownloadedName
	//f.SourcePath = c.SourcePath
	return f
}

func mkTags(s string) map[string]string {
	m := map[string]string{}
	for _, part := range strings.Split(s, ",") {
		m[part] = ""
	}
	return m
}

func parseIntLazy(s string) int {
	if s == "" {
		return 0
	}
	intVar, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return intVar
}

func (f Film) ToCSV(sigDB *filehash.DB) []string {
	var items []string = []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}

	hh := ""
	if f.Hh != 0 {
		hh = fmt.Sprintf("%d", f.Hh)
	}

	duration := fmt.Sprintf("%d", f.Duration) // Duration
	if _, ok := f.Tags["dxx"]; ok {
		delete(f.Tags, "dxx")
		duration = "0"
	}
	if f.Duration == 0 {
		duration = ""
	}

	items[0] = sigDB.GetSignature(f.Filename)       // sha256
	items[1] = ""                                   // OrderPrefix
	items[2] = f.Title                              // Title
	items[3] = f.SourceSite                         // SourceSite
	items[4] = strings.Join(f.Keywords, ",")        // ContentTags
	items[5] = strings.Join(maps.Keys(f.Tags), ",") // MetaTags
	items[6] = duration                             // Duration
	items[7] = hh                                   // hh
	items[8] = f.Room                               // PartyScreen
	items[9] = f.URL                                // URL
	items[10] = f.Author                            // CreatorName
	items[11] = f.SongTitle                         // SongTitle
	items[12] = f.Filename                          // Filename
	items[13] = f.FileExt                           // FileExt
	items[14] = f.OrigFName                         // OriginalFilename

	return items
}
