package resolver

import (
	"fmt"

	"github.com/killerquinn/vdmeta/metadata/config"
	"github.com/killerquinn/vdmeta/metadata/models"
	"github.com/killerquinn/vdmeta/metadata/providers/instagram"
)

type Extractor struct {
	cfg *config.Config
}

func NewExtractor(cfg *config.Config) *Extractor {
	return &Extractor{cfg: cfg}
}

// main func, others doesnt work
// . First that u need is initilize config by config.LoadConf(), after it initialize new extractor by resolver.NewExtrector, then use extractor.(needed extractor to use'instagram is only works right now')
func (ex *Extractor) ExtractIg(RawUrl string) (*models.InstagramContent, error) {

	conf, err := instagram.NewConfigLoader(ex.cfg)
	if err != nil {
		return nil, err
	}
	content, err := conf.ExtractParts(RawUrl)
	if err != nil || content.Type == "stories" || content.Type == "highlights" {
		return nil, fmt.Errorf("i support instagram, but i cant reconize your link, can u sure is it right, please? err - %v", err)
	}
	meta, err := conf.ExtractLink(RawUrl)
	if err != nil {
		return nil, err
	}
	return &models.InstagramContent{
		VideoLink: meta.VideoLink,
		Author:    meta.Author,
	}, nil
}

func (ex *Extractor) ExtractYt(RawUrl string) string {
	return "youtube doesnt supports right now"
}
