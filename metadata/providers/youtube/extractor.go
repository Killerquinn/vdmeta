package youtube

import (
	"fmt"
	"net/url"
	"strings"
)

//in process

func ExtractYt(rawUrl string) string {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if values, ok := parsedUrl.Query()["v"]; ok && len(values) > 0 {
		return values[0]
	}

	if parsedUrl.Host == "youtu.be" {
		return strings.TrimPrefix(parsedUrl.Path, "/")
	}

	return ""
}
