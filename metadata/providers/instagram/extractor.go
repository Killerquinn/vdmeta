package instagram

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"vdmeta/metadata/dto"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrNotSupportedLink = errors.New("that link is not supported yet")
)

func ExtractIg(RawUrl string) string {
	parsedurl, err := url.Parse(RawUrl)
	if err != nil {
		fmt.Printf("error %d", err)
		return ""
	}

	re := regexp.MustCompile(`^([^/]+)/([^/]+)/([^/]+)/`)
	match := re.FindStringSubmatch(parsedurl.Path)
	if len(match) > 2 {
		switch {
		case match[0] == "reel" || match[0] == "reels" || match[0] == "p":
			if len(match[1]) > 0 && len(match[1]) < 25 {
				ExtractMeta(RawUrl, match[1])
			}
		case match[0] == "stories":
			if len(match[2]) > 0 && len(match[2]) < 25 && match[1] == "highlights" {
				ExtractMeta(RawUrl, match[2])
			}
			fmt.Printf("stories actually not supported, only highlight version err: %v", ErrNotSupportedLink)

		default:
			fmt.Println(ErrNotSupportedLink)
		}

	}

	return ""
}

func ExtractMeta(url string, id string) *dto.IgMeta {
	const op = "instagram.extract_meta"

	htmlpage, err := http.Get(url)
	if err != nil {
		log.Fatal(err, op)
		return nil
	}

	defer htmlpage.Body.Close()

	var metapack *dto.IgMeta

	querydc, err := goquery.NewDocumentFromReader(htmlpage.Body)
	if err != nil {
		log.Fatal(err, op)
		return nil
	}

	querydc.Find("a")

}
