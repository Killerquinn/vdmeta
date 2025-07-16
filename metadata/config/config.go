package config

type Config struct {
	UserAgent      string
	Accept         string
	AcceptLanguage string
	MaxRetries     int
	Connection     string
	NeededTag      string
	TextKey        string
}

// config for get request to site
// to start use library you need to write method NewConfigLoader() to get avaible provider.
// in case ExtractIg() is only avaible now, because library in MVP version
// after you get needed provider, just call NewOpts() method to get direct link methods
// of course you can use just ExtractIg(), but it may cause bugs, NewOpts() was created in was created for scalability reasons, for future updates
func LoadConf() *Config {
	return &Config{
		UserAgent:      "Mozilla/5.0 (Linux; Android 13; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36",
		Accept:         "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng, application/json, text/javascript, */*; q=0.01",
		AcceptLanguage: "en-US,en;q=0.9,ru;q=0.8",
		MaxRetries:     6,
		Connection:     "keep-alive",
		NeededTag:      "script[type='application/json']",
		TextKey:        "video_versions",
	}
}
