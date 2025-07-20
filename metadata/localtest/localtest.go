package main

import (
	"fmt"

	cfgg "github.com/killerquinn/vdmeta/metadata/config"
	"github.com/killerquinn/vdmeta/metadata/resolver"
)

func main() {

	cnfg := cfgg.LoadConf()
	if cnfg == nil {
		fmt.Println("config is nil, check LoadConf")
		return
	}

	extractor := resolver.NewExtractor(cnfg)
	result, err := extractor.ExtractIg("https://www.instagram.com/reels/DLZiiAju4tl/")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Author:", result.Author)
	fmt.Println("Video Links:", result.VideoLink[103]) //103 - best quality, best frames per second rate/101 - fps rate worse
	fmt.Println(len(result.VideoLink))
}
