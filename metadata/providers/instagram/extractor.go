package instagram

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/killerquinn/vdmeta/metadata/config"
	"github.com/killerquinn/vdmeta/metadata/models"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrNotSupportedLink = errors.New("that link is not supported yet")
)

type ConfLoader struct {
	cfg *config.Config
}

func NewConfigLoader() (*ConfLoader, error) {
	cfg := config.LoadConf()
	if cfg == nil {
		return nil, fmt.Errorf("cannot load config")
	}
	return &ConfLoader{
		cfg: cfg,
	}, nil
}

func (c *ConfLoader) ExtractIg(RawUrl string) (*models.InstagramContent, error) {
	content, err := c.ExtractParts(RawUrl)
	if err != nil || content.Type == "stories" || content.Type == "highlights" {
		return nil, fmt.Errorf("i support instagram, but i cant reconize your link, can u sure is it right, please?")
	}
	meta, err := c.ExtractLink(RawUrl)
	if err != nil {
		return nil, err
	}
	return &models.InstagramContent{
		VideoLink: meta.VideoLink,
		Author:    meta.Author,
	}, nil
}

func (c *ConfLoader) ExtractParts(RawUrl string) (*models.IgMeta, error) {
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
		return &models.IgMeta{
			Type: "post",
			ID:   parts[1],
		}, nil

	case "reel":
		if len(parts) < 2 {
			return nil, fmt.Errorf("missing reel ID in link")
		}
		return &models.IgMeta{
			Type: "reel",
			ID:   parts[1],
		}, nil

	default:
		return nil, fmt.Errorf("unsupported type of material: %s", parts[0])
	}
}

func (c *ConfLoader) SelectorRetryAdditional(rawUrl string) (string, bool, error) {
	jsonString := ""
	req, err := http.NewRequest("GET", rawUrl, nil)
	if err != nil {
		return "", false, err
	}

	req.Header.Set("User-Agent", c.cfg.UserAgent)
	req.Header.Set("Referer", "https://www.instagram.com/")
	req.Header.Set("Accept", c.cfg.Accept)
	req.Header.Set("Accept-Language", c.cfg.AcceptLanguage)
	req.Header.Set("Connection", c.cfg.Connection)
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

	godoc.Find(c.cfg.NeededTag).Each(func(i int, s *goquery.Selection) {
		jsonText := s.Text()

		if !strings.Contains(jsonText, c.cfg.TextKey) {
			return
		} else {
			jsonString = jsonText
		}
	})

	return jsonString, true, nil
}

func (c *ConfLoader) ExtractLink(rawUrl string) (*models.InstagramContent, error) {
	const op = "instagram.extract_meta"

	maxRetries := c.cfg.MaxRetries
	jsonString := ""
	var urls []string
	var author string

	for i := 0; i < maxRetries; i++ {
		currentJson, found, err := c.SelectorRetryAdditional(rawUrl)
		if currentJson == "" || !found || err != nil {
			fmt.Println("new retry...")
			time.Sleep(time.Second)
		} else {
			jsonString = currentJson
			break
		}
	}
	mustCompile := fmt.Sprintf(`%s:(\[.*?\])`, c.cfg.TextKey)
	regexpBlock := regexp.MustCompile(mustCompile)
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
	if len(urls) == 0 {
		return nil, fmt.Errorf("there is no urls, retry it")
	}

	return &models.InstagramContent{
		VideoLink: urls,
		Author:    author,
	}, nil

}
