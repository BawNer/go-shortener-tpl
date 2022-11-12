package app

const (
	defaultServerAddr = "localhost:8080"
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
