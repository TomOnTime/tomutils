package filminventory

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/TomOnTime/tomutils/vidnamer/filehash"
	"golang.org/x/exp/maps"

	"gopkg.in/yaml.v2"
)

type Film struct {
	Signature  string            `yaml:"md5"` // The md5 has or "" if we don't know it.
	Filename   string            `yaml:"filename,omitempty"`
	FileExt    string            `yaml:"filext"`
	URL        string            `yaml:"url,omitempty"`
	SourceSite string            `yaml:"sourcesite,omitempty"`
	Title      string            `yaml:"title"`
	Author     string            `yaml:"author,omitempty"`
	Keywords   []string          `yaml:"keywords"` // Keywords (topics in the video)
	Tags       map[string]string `yaml:"tags"`     // Meta tags (describes the video's qualities)
	Duration   int               `yaml:"duration"`
	Hh         int               `yaml:"hh,omitempty"`
	Room       string            `yaml:"room,omitempty"`
	Test       string            `yaml:"test,omitempty"` // Good test of strength? (level)
}

func FromYamlfile(filename string) (r []Film, err error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can't read %q: %w", filename, err)
	}
	err = yaml.Unmarshal(yamlFile, &r)
	if err != nil {
		return nil, fmt.Errorf("can't parse %q: %w", filename, err)
	}
	return r, nil
}

func ParseFilename(filename string) Film {
	// `title__site__keywords__designation.mp4`

	foundFilename := filename

	// "extenstion"
	var ext string
	for _, e := range []string{"mp4", "mpg", "mpeg", "mov"} {
		if strings.HasSuffix(filename, "."+e) {
			filename = strings.TrimSuffix(filename, "."+e)
			filename = strings.TrimSuffix(filename, ".") // In case filename..mp4
			filename = strings.TrimSpace(filename)       // in case filename\ .mp4
			ext = e
			break
		}
	}
	//fmt.Printf("DEBUG: filename2 = %q\n", filename)
	//fmt.Printf("DEBUG: ext = %q\n", ext)

	major := strings.Split(filename, `__`)

	// "title"
	title := major[0]
	title = strings.TrimSpace(title)

	// "site"
	var site string
	if len(major) > 1 {
		site = major[1]
	}
	site = strings.TrimSpace(site)

	f := Film{
		//Signature:    not set by this function
		Filename:   foundFilename,
		Title:      title,
		SourceSite: site,
		FileExt:    ext,
		Tags:       map[string]string{},
	}

	// "keywords"
	if len(major) > 2 {
		f.Keywords = strings.Split(major[2], "-")
		for i := range f.Keywords {
			if f.Keywords[i] == "0fvoiceover" {
				f.Keywords[i] = "fvoiceover"
			}
		}
	}

	// "designation"
	if len(major) > 3 {
		combined := major[3]
		combined = strings.TrimSuffix(combined, "."+ext)
		combined = strings.ToLower(combined)
		combined = strings.Replace(combined, "_", "-", -1)
		for _, d := range strings.Split(combined, "-") {
			if d == "" {
				continue
			}

			if d == "hh1" || d == "h1" {
				f.Hh = 1
			} else if d == "hh2" || d == "h2" {
				f.Hh = 2

			} else if d == "main" {
				f.Room = d
			} else if d == "side" {
				f.Room = d
			} else if d == "both" {
				f.Room = d

			} else if matched, err := regexp.MatchString(`^d\d+`, d); err == nil && matched {
				dur, _ := strconv.Atoi(d[1:])
				f.Duration = dur

			} else {
				f.Tags[d] = ""
			}
		}
	}

	// TODO(tlim) If "by X" in the title, split that out as the Author.
	//f.Author = ""

	return f
}

func (f Film) ToYaml() ([]byte, error) {
	h := []Film{f}
	return yaml.Marshal(&h)
}

func ExistingFilename(signature string, md5db []filehash.Info) string {
	for _, m := range md5db {
		if m.Signature == signature {
			return m.Filename
		}
	}
	return ""
}

func (f Film) DesiredFilename() string {

	var title, site, keywords, designation, ext string

	title = f.Title
	title = strings.TrimSpace(title)

	if f.Author != "" {
		title = title + " by " + f.Author
	}

	site = f.SourceSite

	keywords = strings.Join(f.Keywords, "-")

	// Build the designation
	var dparts []string
	if f.Hh != 0 {
		n := fmt.Sprintf("hh%d", f.Hh)
		dparts = append(dparts, n)
	}
	if f.Room != "" {
		n := f.Room
		dparts = append(dparts, n)
	}
	if f.Test != "" {
		n := f.Test
		dparts = append(dparts, n)
	}
	if len(f.Tags) > 0 {
		k := maps.Keys(f.Tags)
		sort.Strings(k)
		n := strings.Join(k, "-")
		dparts = append(dparts, n)
	}
	if f.Duration != 0 {
		n := fmt.Sprintf("d%02d", f.Duration)
		dparts = append(dparts, n)
	} else {
		n := "dXX"
		dparts = append(dparts, n)
	}
	designation = strings.Join(dparts, "-")

	ext = f.FileExt
	//	fmt.Printf("DEBUG: ext1=%q\n", ext)
	if ext == "" {
		ext = "mp4"
	}
	//fmt.Printf("DEBUG: ext2=%q\n", ext)

	return strings.Join([]string{title, site, keywords, designation}, "__") + "." + ext

}
