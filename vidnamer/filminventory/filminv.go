package filminventory

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/TomOnTime/tomutils/vidnamer/filehash"
	"golang.org/x/exp/maps"

	"gopkg.in/yaml.v2"
)

type Film struct {
	Signature  string            `yaml:"md5"` // The md5 has or "" if we don't know it.
	Filename   string            `yaml:"filename,omitempty"`
	Title      string            `yaml:"title"`
	Author     string            `yaml:"author,omitempty"`
	SourceSite string            `yaml:"sourcesite,omitempty"`
	Keywords   []string          `yaml:"keywords,omitempty"` // Keywords (topics in the video)
	Hh         string            `yaml:"hh,omitempty"`
	Room       string            `yaml:"room,omitempty"`
	Test       string            `yaml:"test,omitempty"`
	Duration   string            `yaml:"duration,omitempty"`
	FileExt    string            `yaml:"filext,omitempty"`
	Tags       map[string]string `yaml:"tags,omitempty"` // Meta tags
	//
}

func FromYamlfile(filename string) (r []Film, err error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Can't read %q: %w", filename, err)
	}
	err = yaml.Unmarshal(yamlFile, &r)
	if err != nil {
		return nil, fmt.Errorf("Can't parse %q: %w", filename, err)
	}
	return r, nil
}

func ParseFilename(filename string) Film {
	// `title__site__keywords__designation.mp4`

	foundFilename := filename

	// "extenstion"
	var ext string
	for _, e := range []string{".mp4", ".mpeg"} {
		if strings.HasSuffix(filename, e) {
			filename = strings.TrimSuffix(filename, e)
			filename = strings.TrimSpace(filename)
			ext = strings.TrimPrefix(e, ".")
			break
		}
	}
	//fmt.Printf("DEBUG: filename2 = %q\n", filename)

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
		//md5hash:    not set by this function
		Filename:   foundFilename,
		Title:      title,
		SourceSite: site,
		FileExt:    ext,
		Tags:       map[string]string{},
	}

	// "keywords"
	if len(major) > 2 {
		f.Keywords = strings.Split(major[2], "-")
	}

	// "designation"
	if len(major) > 3 {
		f.Duration = "XX"
		combined := major[3]
		combined = strings.ToLower(combined)
		combined = strings.Replace(combined, "_", "-", -1)
		for _, d := range strings.Split(combined, "-") {
			if d == "" {
				continue
			}
			//fmt.Printf("DEBUG: designation d=%v\n", d)
			if d == "hh1" || d == "h1" {
				f.Hh = "1"
			} else if d == "hh2" || d == "h2" {
				f.Hh = "2"
			} else if d == "finale" {
				f.Tags["finale"] = ""
			} else if d == "calm" {
				f.Tags["calm"] = ""
			} else if d == "humiliation" {
				f.Tags["calm"] = ""
			} else if d == "pornaddict" {
				f.Tags["pornaddict"] = ""
			} else if d == "long" {
				f.Tags["long"] = ""
			} else if d == "main" {
				f.Room = d
			} else if d == "side" {
				f.Room = d
			} else if d == "both" {
				f.Room = d
			} else if matched, err := regexp.MatchString(`^d\d+`, d); err == nil && matched {
				//fmt.Printf("DEBUG: matched, err: %v %v\n", matched, err)
				//} else if strings.HasPrefix(d, "d") { // FIXME: r`^d\d+`
				f.Duration = d[1:]
				// f.Test = ""
				// f.Duration = ""
			} else {
				fmt.Printf("WARNING: Unknown designation: %q\n", d)
			}
		}
	}

	// If "by X" in the title, split that out as the Author.
	//f.Author = ""

	return f
}

func (f Film) ToYaml() ([]byte, error) {
	h := []Film{f}
	return yaml.Marshal(&h)
}

// func (f Film) ExistingFilename() string {
// 	return f.Filename
// }
//		existing := invItem.ExistingFilename(invItem.Signature, md5db)
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
	//	if strings.HasPrefix(title, `'`) && strings.HasSuffix(title, `'`) {
	//		fmt.Printf("DEBUG: TITLE1=%q\n", title)
	//		title = title[1 : len(title)-2]
	//		fmt.Printf("DEBUG: TITLE2=%q\n", title)
	//	}

	if f.Author != "" {
		title = title + " by " + f.Author
	}

	site = f.SourceSite

	keywords = strings.Join(f.Keywords, "-")

	// Build the designation
	if f.Hh != "" {
		designation = strings.Join([]string{designation, "hh" + f.Hh}, "-")
	}
	if f.Room != "" {
		designation = strings.Join([]string{designation, f.Room}, "-")
	}
	if f.Test != "" {
		designation = strings.Join([]string{designation, f.Test}, "-")
	}
	if len(f.Tags) > 0 {
		t := strings.Join(maps.Keys(f.Tags), "-")
		designation = strings.Join([]string{designation, t}, "-")
	}

	ext = f.FileExt

	return strings.Join([]string{title, site, keywords, designation}, "__") + "." + ext

}
