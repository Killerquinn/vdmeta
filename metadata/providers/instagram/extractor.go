package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://www.instagram.com/")
	req.Header.Set("Origin", "https://www.instagram.com")

	// maybe smth headers more needs to it...

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	jsonData := ""

	htmldoc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	htmldoc.Find("script[type='application/json']").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "__additionalDataLoaded") {
			jsonData = s.Text()
		}
	})

	if jsonData == "" {
		fmt.Println("json data is empty")
	}

	var reelData map[string]interface{}
	err = json.Unmarshal([]byte(jsonData), &reelData)
	if err != nil {
		fmt.Println(err)
	}

	items, ok := reelData["items"].([]interface{})
	if !ok {
		fmt.Println("have no items")
	}

	item, ok := items[0].(map[string]interface{})
	if !ok {
		fmt.Println("have no item")
	}

	videourls, ok := item["video_versions"].([]interface{})
	if !ok {
		fmt.Println("have no videourls")
	}

	videourl, ok := videourls[0].(map[string]interface{})["url"].(string)
	if !ok {
		fmt.Println("have no videourl")
	}

	return &dto.IgMeta{}

}
