package app

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"
)

type ConfigApp struct {
	ServerAddr      string
	BaseURL         string
	FileStoragePath string
	Secret          string
	DB              string
}

func NewConfigApp() *ConfigApp {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, load default values")
	}

	const (
		defaultServerAddr      = "localhost:8080"
		defaultBaseURL         = "http://127.0.0.1:8080"
		defaultFileStoragePath = ""
		defaultSecret          = "u-nya-nya-mo-ni-ni"
		defaultDB              = ""
	)

	// dbConn = postgresql://postgres:Ewelli55dxx@localhost:5432/shortener

	var (
		serverAddr      string
		baseURL         string
		fileStoragePath string
		secret          string
		db              string
	)

	flag.StringVarP(&serverAddr, "a", "a", defaultServerAddr, "-a to set server address")
	flag.StringVarP(&baseURL, "b", "b", defaultBaseURL, "-b to set base url")
	flag.StringVarP(&fileStoragePath, "f", "f", defaultFileStoragePath, "-f to set location storage files")
	flag.StringVarP(&secret, "s", "s", defaultSecret, "-s to secret key")
	flag.StringVarP(&db, "d", "d", defaultDB, "-d to set db address")

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

	if os.Getenv("SECRET") != "" {
		secret = os.Getenv("SECRET")
	}

	if os.Getenv("DATABASE_DSN") != "" {
		db = os.Getenv("DATABASE_DSN")
	}

	return &ConfigApp{
		ServerAddr:      serverAddr,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
		Secret:          secret,
		DB:              db,
	}
}

var (
	Config = NewConfigApp()
)
