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
	result, err := extractor.ExtractIg("https://www.instagram.com/reels/DJKBATZAvyA/")
	if err != nil {
		fmt.Println("your library is piece of shit")
		return
	}

	fmt.Println("Author:", result.Author)
	fmt.Println("Video Links:", result.VideoLink[0])
}
