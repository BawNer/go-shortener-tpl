package app

const (
	defaultServerAddr = "127.0.0.1:8080"
	defaultBaseURL    = ""
)

type Config struct {
	ServerAddr string
	BaseURL    string
}

func NewConfig(conf Config) Config {
	if conf.ServerAddr == "" {
		conf.ServerAddr = defaultServerAddr
	}

	if conf.BaseURL == "" {
		conf.BaseURL = defaultBaseURL
	}

	return conf
}
