package app

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

type ConfigApp struct {
	ServerAddr      string
	BaseURL         string
	FileStoragePath string
}

func NewConfigApp() func() *ConfigApp {

	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, load default values")
	}

	const (
		defaultServerAddr = "localhost:8080"
		defaultBaseURL    = "http://127.0.0.1:8080"
	)

	var (
		serverAddr      string
		baseURL         string
		fileStoragePath string
	)

	pflag.StringVar(&serverAddr, "a", defaultServerAddr, "-a to set server address")
	pflag.StringVar(&baseURL, "b", defaultBaseURL, "-b to set base url")
	pflag.StringVar(&fileStoragePath, "f", "", "-f to set location storage files")

	pflag.Parse()

	if os.Getenv("SERVER_ADDRESS") != "" {
		serverAddr = os.Getenv("SERVER_ADDRESS")
	}

	if os.Getenv("BASE_URL") != "" {
		baseURL = os.Getenv("BASE_URL")
	}

	if os.Getenv("FILE_STORAGE_PATH") != "" {
		fileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	}

	return func() *ConfigApp {
		return &ConfigApp{
			ServerAddr:      serverAddr,
			BaseURL:         baseURL,
			FileStoragePath: fileStoragePath,
		}
	}
}

var (
	cfg    = NewConfigApp()
	Config = cfg()
)
