package source

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type DocHandler func(title, body string) error

type Redirect struct {
	Title string `xml:"title,attr"`
}

type Page struct {
	Title string `xml:"title"`
	Redirect Redirect `xml:"redirect"`
	Text string `xml:"revision>text"`
}

func CanonicalizeTitle(title string) string {
	can := strings.ToLower(title)
	can = strings.Replace(can, " ", "_", -1)
	can = url.QueryEscape(can)
	return can
}

var filter, _ = regexp.Compile("^file:.*|^talk:.*|^special:.*|^wikipedia:.*|^wiktionary:.*|^user:.*|^user_talk:.*")

func LoadWiki(wikipedia string, maxIndexCount int, handler DocHandler)(err error){
	var f *os.File
	f, err = os.Open(wikipedia)
	if err != nil {
		return
	}
	defer f.Close()

	d := xml.NewDecoder(f)
	count := 0
	for {
		if maxIndexCount != -1 && count >= maxIndexCount {
			break
		}
		var tok xml.Token
		tok, err = d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			return err
		}

		switch ty := tok.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement := ty.Name.Local
			// ...and its name is "page"
			if inElement == "page" {
				var p Page
				// decode a whole chunk of following XML into the
				// variable p which is a Page (se above)
				_ = d.DecodeElement(&p, &ty)

				// Do some stuff with the page.
				p.Title = CanonicalizeTitle(p.Title)
				m := filter.MatchString(p.Title)
				if !m && p.Redirect.Title == "" {
					err = handler(p.Title, p.Text)
					if err != nil {
						break
					}
					count++
				}
			}
		default:
		}
	}

	fmt.Println("count =", count)
	return
}
