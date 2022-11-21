package app

import (
	"log"
	"os"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"
)

type ConfigApp struct {
	ServerAddr      string
	BaseURL         string
	FileStoragePath string
}

func NewConfigApp() *ConfigApp {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, load default values")
	}

	const (
		defaultServerAddr      = "localhost:8080"
		defaultBaseURL         = "http://127.0.0.1:8080"
		defaultFileStoragePath = ""
	)

	var (
		serverAddr      string
		baseURL         string
		fileStoragePath string
	)

	flag.StringVarP(&serverAddr, "a", "a", defaultServerAddr, "-a to set server address")
	flag.StringVarP(&baseURL, "b", "b", defaultBaseURL, "-b to set base url")
	flag.StringVarP(&fileStoragePath, "f", "f", defaultFileStoragePath, "-f to set location storage files")

	flag.Parse()

	if os.Getenv("SERVER_ADDRESS") != "" {
		serverAddr = os.Getenv("SERVER_ADDRESS")
	}

	if os.Getenv("BASE_URL") != "" {
		baseURL = os.Getenv("BASE_URL")
	}

	if os.Getenv("FILE_STORAGE_PATH") != "" {
		fileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	}

	return &ConfigApp{
		ServerAddr:      serverAddr,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
	}
}

var (
	Config = NewConfigApp()
	Memory = storage.NewMemory(Config.FileStoragePath)
)
