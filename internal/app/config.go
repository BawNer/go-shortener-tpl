package app

const (
	defaultServerAddr = "127.0.0.1:8080"
	defaultBaseURL    = ""
	defaultScheme     = "http"
)

type Config struct {
	ServerAddr string
	BaseURL    string
	Scheme     string
}

func NewConfig(conf Config) Config {
	if conf.ServerAddr == "" {
		conf.ServerAddr = defaultServerAddr
	}

	if conf.BaseURL == "" {
		conf.BaseURL = defaultBaseURL
	}

	if conf.Scheme == "" {
		conf.Scheme = defaultScheme
	}

	return conf
}
