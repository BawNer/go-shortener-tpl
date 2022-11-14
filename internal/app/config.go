package app

const (
	defaultServerAddr = "localhost:8080"
	defaultBaseURL    = "http://127.0.0.1:8080"
)

type Config struct {
	ServerAddr string
	BaseURL    string
}

func NewConfig(conf *Config) *Config {

	if conf.ServerAddr == "" {
		conf.ServerAddr = defaultServerAddr
	}

	if conf.BaseURL == "" {
		conf.BaseURL = defaultBaseURL
	}

	return conf
}
