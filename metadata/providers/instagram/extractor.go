package instagram

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"vdmeta/metadata/dto"
	"vdmeta/metadata/models"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrNotSupportedLink = errors.New("that link is not supported yet")
)

func ExtractIg(RawUrl string) ([]string, error) {
	content, err := ExtractParts(RawUrl)
	if err != nil || content.Type == "stories" || content.Type == "highlights" {

	}
	meta, err := ExtractLink(RawUrl)
	//TODO: add logic, rewrite ExtractMeta func
	return meta.Video
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

	default:
		return nil, fmt.Errorf("unsupported type of material: %s", parts[0])
	}
}

func SelectorRetryAdditional(rawUrl string) (string, bool, error) {
	jsonString := ""
	req, err := http.NewRequest("GET", rawUrl, nil)
	if err != nil {
		return "", false, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 13; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Referer", "https://www.instagram.com/")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng, application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("DNT", "1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Cache-Control", "no-cache")

	do, err := http.DefaultClient.Do(req)
	if err != nil || do == nil {
		return "", false, err
	}

	if do.StatusCode != 200 {
		return "", false, fmt.Errorf("connection not stable")
	}

	defer do.Body.Close()

	godoc, err := goquery.NewDocumentFromReader(do.Body)
	if err != nil {
		return "", false, fmt.Errorf("cannot convert response into goquery document")
	}

	godoc.Find("script[type='application/json']").Each(func(i int, s *goquery.Selection) {
		jsonText := s.Text()

		if !strings.Contains(jsonText, "video_versions") {
			return
		} else {
			jsonString = jsonText
		}
	})

	return jsonString, true, nil
}

func ExtractLink(rawUrl string) (*dto.IgMeta, error) {
	const op = "instagram.extract_meta"

	const maxRetries = 6
	jsonString := ""
	var urls []string
	var author string

	for i := 0; i < maxRetries; i++ {
		currentJson, found, err := SelectorRetryAdditional(rawUrl)
		if currentJson == "" || !found || err != nil {
			fmt.Println("new retry...")
			time.Sleep(time.Second)
		} else {
			jsonString = currentJson
			break
		}
	}

	regexpBlock := regexp.MustCompile(`"video_versions":(\[.*?\])`)
	authorRxpBlock := regexp.MustCompile(`"ig_artist":(\{.*?\})`)
	blockWithUsername := authorRxpBlock.FindAllStringSubmatch(jsonString, -1)
	blocks := regexpBlock.FindAllStringSubmatch(jsonString, -1)
	urlRxp := regexp.MustCompile(`"url":"([^"]+)"`)
	authorPost := regexp.MustCompile(`"username":"([^"]+)"`)
	for _, block := range blocks {
		urlBlock := urlRxp.FindAllStringSubmatch(block[0], 3)
		for _, u := range urlBlock {
			res := strings.ReplaceAll(u[1], "\\", "")
			urls = append(urls, res)
		}
	}
	urls = urls[:3]
	for _, block := range blockWithUsername {
		usernameBlock := authorPost.FindAllStringSubmatch(block[0], 1)
		for _, u := range usernameBlock {
			author = u[1]
		}
	}
	if len(urls) == 0 || author == "" {
		return nil, fmt.Errorf("there is no author or urls, retry it")
	}

	return &dto.IgMeta{
		VideoLink: urls,
		Author:    author,
	}, nil

}
