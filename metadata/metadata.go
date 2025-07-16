package metadata

import (
	"fmt"

	"github.com/killerquinn/vdmeta/metadata/models"
)

type Providers interface {
	ExtractYt(rawUrl string) string
	ExtractIg(rawUrl string) (*models.InstagramContent, error)
}

type Opts struct {
	providers Providers
}

func NewOpts(provider Providers) (*Opts, error) {
	if provider == nil {
		return &Opts{}, nil
	}
	return &Opts{
		providers: provider,
	}, nil
}

func (o *Opts) GetInstaLinks(rawUrl string) (directLink string, author string, err error) {
	const op = "getInstaLinks"
	content, err := o.providers.ExtractIg(rawUrl)
	if err != nil || content == nil {
		return "", "", fmt.Errorf("%s:%v", op, err)
	}
	return content.VideoLink[0], content.Author, nil
}
