package filminventory

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/TomOnTime/tomutils/vidnamer/filehash"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/utf8string"

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
	// If it has a single song, list the title here
	SongTitle string `yaml:"SongTitle,omitempty"`
	// If known, the filename it downloaded to originally.
	OrigFName string `yaml:"OriginalFilename,omitempty"`
	// Source Path (dir/file.ext): Used to store the path to the filename when importing
	SourcePath string
	// Order (prefix for filenames)
	OrderPrefix string `yaml:"filext"`
}

var CSVHeader = "sha256,OrderPrefix,Title,SourceSite,ContentTags,MetaTags,Duration,hh,PartyScreen,URL,CreatorName,SongTitle,Filename,FileExt,OriginalFilename"

type FilmCSV struct {
	Sha256         string `csv:"sha256"`
	OrderPrefix    string `csv:"Order"`
	Title          string `csv:"Title"`
	SourceSite     string `csv:"SourceSite"`
	ContentTags    string `csv:"ContentTags"`
	MetaTags       string `csv:"MetaTags"`
	Duration       string `csv:"Duration"`
	Hh             string `csv:"hh"`
	PartyScreen    string `csv:"PartyScreen"`
	URL            string `csv:"URL"`
	CreatorName    string `csv:"CreatorName"`
	SongTitle      string `csv:"SongTitle"`
	Filename       string `csv:"Filename"`
	FileExt        string `csv:"FileExt"`
	DownloadedName string `csv:"DownloadedName"`
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

func FromManyFilepaths(filenames []string, sigs *filehash.DB) []Film {
	var result []Film
	for _, name := range filenames {
		_, file := filepath.Split(name)

		film := ParseFilename((file))
		film.SourcePath = name
		film.Signature = sigs.GetSignature(file)
		film.Filename = file
		film.FileExt = filepath.Ext(file)
		film.FileExt = filepath.Ext(file)

		result = append(result, film)
	}
	return result
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
	origtitle := title
	title = strings.TrimSpace(title)
	title = strings.TrimPrefix(title, "SpankBang.com_")
	title = strings.TrimPrefix(title, "xvideos.com_")

	// "site"
	var site string
	if len(major) > 1 {
		site = major[1]
	}
	site = strings.TrimSpace(site)
	if site == "" {
		if strings.Contains(strings.ToLower(origtitle), "spankbang") {
			site = "SpankBang"
		}
		if strings.Contains(strings.ToLower(origtitle), "xvideos") {
			site = "xvideos"
		}
	}

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
	//fmt.Printf("DEBUG: 1 fn=%q dparts=%v\n", f.Title, dparts)
	if f.Room != "" {
		n := f.Room
		dparts = append(dparts, n)
	}
	//fmt.Printf("DEBUG: 2 fn=%q dparts=%v\n", f.Title, dparts)
	if f.Test != "" {
		n := f.Test
		dparts = append(dparts, n)
	}
	//fmt.Printf("DEBUG: 3 fn=%q dparts=%v\n", f.Title, dparts)
	if len(f.Tags) > 0 {
		k := maps.Keys(f.Tags)
		sort.Strings(k)
		n := strings.Join(k, "-")
		dparts = append(dparts, n)
	}
	//fmt.Printf("DEBUG: 4 fn=%q dparts=%v\n", f.Title, dparts)
	if f.Duration != 0 {
		n := fmt.Sprintf("d%02d", f.Duration)
		dparts = append(dparts, n)
	} else {
		n := "dXX"
		dparts = append(dparts, n)
	}

	//fmt.Printf("DEBUG: FINAL fn=%q dparts=%v\n", f.Title, dparts)
	designation = strings.Join(dparts, "-")

	ext = f.FileExt
	//	fmt.Printf("DEBUG: ext1=%q\n", ext)
	if ext == "" {
		ext = "mp4"
	}
	//fmt.Printf("DEBUG: ext2=%q\n", ext)

	final := strings.Join([]string{title, site, keywords, designation}, "__") + "." + ext
	if !utf8string.NewString(final).IsASCII() {
		fmt.Printf("# WARNING: Non-ASCII string: %q\n", final)
	}
	return final

}
