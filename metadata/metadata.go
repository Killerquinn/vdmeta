package metadata

type Providers interface {
	ExtractYt(rawUrl string) string
	ExtractIg(rawUrl string) string
}
