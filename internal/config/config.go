package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	GrpcServer `yaml:"grpc_server"`
	Db
	UrlsDb
}

type HTTPServer struct {
	ServerPort string `yaml:"server_port"`
}

type GrpcServer struct {
	Addr string `yaml:"grpc_server_address"`
}

type Db struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"db_Name"`
	User     string `yaml:"db_User"`
	Password string `yaml:"db_Password"`
	Option   string `yaml:"db_option"`
}

type UrlsDb struct {
	Db string
}

func MustInit(configPath string) *Config {
	godotenv.Load(configPath)

	db := &Db{
		Driver:   MustGetEnv("DB_DRIVER"),
		Host:     MustGetEnv("DB_HOST"),
		Port:     MustGetEnv("DB_PORT"),
		Name:     MustGetEnv("DB_NAME"),
		User:     MustGetEnv("DB_USER"),
		Password: MustGetEnv("DB_PASSWORD"),
		Option:   MustGetEnv("DB_OPTION"),
	}

	var urlDb = buildDbConnectUrl(db)

	env := MustGetEnv("ENV")

	// определяем переменную grpc сервера в dev/prod режиме
	grpcAddrServer := MustGetEnv("GRPC_SERVER_ADDRESS")
	if env == "DEV" {
		grpcAddrServer = MustGetEnv("GRPC_SERVER_ADDRESS_DEV")
	}

	return &Config{
		Env: env,
		HTTPServer: HTTPServer{
			ServerPort: MustGetEnv("SERVER_PORT"),
		},
		GrpcServer: GrpcServer{
			Addr: grpcAddrServer,
		},
		Db: *db,
		UrlsDb: UrlsDb{
			Db: urlDb,
		},
	}
}

func PathDefault(workDir string, filename *string) string {
	if filename == nil {
		return filepath.Join(workDir, ".env")
	}

	return filepath.Join(workDir, *filename)
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("no variable in env: %s", key)
	}
	return value
}

func MustGetEnvAsInt(name string) int {
	valueStr := MustGetEnv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return -1
}

func ParseConfigPathFromCl(currentDir string) string {
	return PathDefault(currentDir, nil)
}

func buildDbConnectUrl(db *Db) string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s",
		db.Driver,
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.Option,
	)
}
