package filminventory

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Film struct {
	Signature  string   // The md5 has or "" if we don't know it.
	filename   string   // The existing/found filename
	Title      string   `yaml:"title"`
	Author     string   `yaml:"author"`
	SourceSite string   `yaml:"sourcesite"`
	Keywords   []string `yaml:"keywords"`
	Hh         string   `yaml:"hh"`
	Room       string   `yaml:"room"`
	Test       string   `yaml:"test"`
	Duration   string   `yaml:"duration"`
	Format     string   `yaml:"format"`
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

	major := strings.Split(filename, `__`)
	// "title"
	title := major[0]
	// "site"
	var source string
	if len(major) > 0 {
		source = major[1]
	}
	// "extenstion"
	var ext string
	for _, e := range []string{".mp4", ".mpeg"} {
		if strings.HasPrefix(filename, e) {
			filename = strings.TrimSuffix(filename, e)
			ext = e
			break
		}
	}

	f := Film{
		//md5hash:    not set by this function
		filename:   filename,
		Title:      title,
		SourceSite: source,
		Format:     ext,
	}

	// "keywords"
	if len(major) > 1 {
		f.Keywords = strings.Split(major[1], "-")
	}

	// "designation"
	if len(major) > 2 {
		for _, d := range strings.Split(major[2], "-") {
			if d == "hh1" {
				f.Hh = "1"
			} else if d == "hh2" {
				f.Hh = "2"
			} else if strings.HasPrefix(d, "d") { // FIXME: r`^d\d+`
				f.Duration = d[1:]
				// f.Room = ""
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
	return yaml.Marshal(&f)
}
func (f Film) DesiredFilename() string {
	return ""
}
func (f Film) ExistingFilename() string {
	return ""
}
