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
		serverAddr      = defaultServerAddr
		baseURL         = defaultBaseURL
		fileStoragePath = ""
	)

	pflag.StringVar(&serverAddr, "a", defaultServerAddr, "Server Address")
	pflag.StringVar(&baseURL, "b", defaultBaseURL, "Base URL")
	pflag.StringVar(&fileStoragePath, "f", "", "File Location Storage")

	if os.Getenv("SERVER_ADDRESS") != "" {
		serverAddr = os.Getenv("SERVER_ADDRESS")
	}

	if os.Getenv("BASE_URL") != "" {
		baseURL = os.Getenv("BASE_URL")
	}

	if os.Getenv("FILE_STORAGE_PATH") != "" {
		fileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	}

	pflag.Parse()

	return func() *ConfigApp {
		return &ConfigApp{
			ServerAddr:      serverAddr,
			BaseURL:         baseURL,
			FileStoragePath: fileStoragePath,
		}
	}
}
