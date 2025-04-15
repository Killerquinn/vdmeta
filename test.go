package main

import (
	"fmt"
	"net/url"
	"strings"
)

func main() {
	parsedurl, err := url.Parse("https://www.instagram.com/stories/highlights/17936732501814509/")
	if err != nil {
		fmt.Println(err)
	}

	videoFormat := strings.Trim(parsedurl.Path, "/")

	fmt.Println(parsedurl.Path)
	fmt.Println(videoFormat)
}
