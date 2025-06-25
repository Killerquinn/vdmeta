package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"vdmeta/metadata/dto"
	"vdmeta/metadata/models"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrNotSupportedLink = errors.New("that link is not supported yet")
)

func ExtractIg(RawUrl string) string {
	content, err := ExtractParts(RawUrl)
	//TODO: add logic, rewrite ExtractMeta func
}

func ExtractParts(RawUrl string) (*models.InstagramContent, error) {
	u, err := url.Parse(RawUrl)

	if err != nil {
		return nil, fmt.Errorf("invalid url: %v", err)
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty path")
	}

	switch parts[0] {
	case "p":
		if len(parts) < 2 {
			return nil, fmt.Errorf("missing post ID in link")
		}
		return &models.InstagramContent{
			Type: "post",
			ID:   parts[1],
		}, nil

	case "reel":
		if len(parts) < 2 {
			return nil, fmt.Errorf("missing reel ID in link")
		}
		return &models.InstagramContent{
			Type: "reel",
			ID:   parts[1],
		}, nil

	case "stories":
		if len(parts) < 3 {
			return nil, fmt.Errorf("missing username of story owner in link")
		}
		return &models.InstagramContent{
			Type:   "stories",
			IgUser: parts[1],
			ID:     parts[2],
		}, nil

	default:
		return nil, fmt.Errorf("unsupported type of material: %s", parts[0])
	}
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
		if strings.Contains(s.Text(), "data") {
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

	return &dto.IgMeta{
		Video: videourl,
	}

}
