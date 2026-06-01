package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	EnvLocal      = "local"
	EnvStaging    = "staging"
	EnvProduction = "production"
)

type Bootstrap struct {
	Server   Server   `yaml:"server"`
	Data     Data     `yaml:"data"`
	JWT      JWT      `yaml:"jwt"`
	BooksKey BooksKey `yaml:"booksKey"`
	Email    Email    `yaml:"email"`
	R2       R2       `yaml:"r2"`
	OpenAI   OpenAI   `yaml:"openai"`
	AWS      AWS      `yaml:"aws"`
}

type Server struct {
	HTTP HTTP `yaml:"http"`
	CORS CORS `yaml:"cors"`
}

type HTTP struct {
	Addr string `yaml:"addr"`
}

type CORS struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
}

type Data struct {
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
}

type Database struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWT struct {
	AdminSecret  string `yaml:"admin_secret"`
	ClientSecret string `yaml:"client_secret"`
}

type BooksKey struct {
	AESKey string `yaml:"aes_key"`
	AESIV  string `yaml:"aes_iv"`
	RSAPub string `yaml:"rsa_pub_key"`
}

type Email struct {
	APIKey string `yaml:"api_key"`
	From   string `yaml:"from"`
}

type R2 struct {
	AccountID       string `yaml:"account_id"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	BucketName      string `yaml:"bucket_name"`
	Endpoint        string `yaml:"endpoint"`
}

type OpenAI struct {
	APIKey string `yaml:"api_key"`
}

type AWS struct {
	Region          string `yaml:"region"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

var envBindings = map[string]string{
	"server.http.addr":      "APP_PORT",
	"data.database.dsn":     "DATABASE_DSN",
	"data.redis.addr":       "REDIS_ADDR",
	"data.redis.password":   "REDIS_PASSWORD",
	"jwt.client_secret":     "CLIENT_SECRET",
	"jwt.admin_secret":      "ADMIN_SECRET",
	"booksKey.aes_key":      "AES_KEY",
	"booksKey.aes_iv":       "AES_IV",
	"booksKey.rsa_pub_key":  "RSA_PUB_KEY",
	"email.api_key":         "SENDGRID_API_KEY",
	"email.from":            "EMAIL_FROM",
	"r2.account_id":         "R2_ACCOUNT_ID",
	"r2.access_key_id":      "R2_ACCESS_KEY_ID",
	"r2.secret_access_key":  "R2_SECRET_ACCESS_KEY",
	"r2.bucket_name":        "R2_BUCKET_NAME",
	"openai.api_key":        "OPENAI_API_KEY",
	"aws.region":            "AWS_REGION",
	"aws.access_key_id":     "AWS_ACCESS_KEY_ID",
	"aws.secret_access_key": "AWS_SECRET_ACCESS_KEY",
}

func init() {
	loadDotEnv()
}

func loadDotEnv() {
	_ = godotenv.Load()
	env := os.Getenv("APP_ENV")
	if env != "" && env != EnvLocal {
		if err := godotenv.Load(fmt.Sprintf(".env.%s", env)); err == nil {
			return
		}
	}
	cwd, _ := os.Getwd()
	paths := []string{
		filepath.Join(cwd, ".env"),
		filepath.Join(cwd, "../.env"),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			_ = godotenv.Load(p)
		}
	}
}

func Load(configPath string) (*Bootstrap, error) {
	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	overrideFromEnv(v)

	var bc Bootstrap
	if err := v.Unmarshal(&bc); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &bc, nil
}

func overrideFromEnv(v *viper.Viper) {
	for key, env := range envBindings {
		if val := os.Getenv(env); val != "" {
			if key == "server.http.addr" {
				if _, err := strconv.Atoi(val); err == nil {
					val = ":" + val
				}
			}
			v.Set(key, val)
		}
	}
	if os.Getenv("REDIS_ADDR") == "" {
		if host := os.Getenv("REDIS_HOST"); host != "" {
			port := os.Getenv("REDIS_PORT")
			if port == "" {
				port = "6379"
			}
			v.Set("data.redis.addr", host+":"+port)
		}
	}
	if os.Getenv("REDIS_DB") != "" {
		if db, err := strconv.Atoi(os.Getenv("REDIS_DB")); err == nil {
			v.Set("data.redis.db", db)
		}
	}
	if os.Getenv("DATABASE_DSN") == "" {
		host := os.Getenv("DB_HOST")
		if host != "" {
			port := os.Getenv("DB_PORT")
			user := os.Getenv("DB_USER")
			password := os.Getenv("DB_PASSWORD")
			dbName := os.Getenv("DB_NAME")
			if port == "" {
				port = "5432"
			}
			dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=prefer", host, port, user, password, dbName)
			v.Set("data.database.dsn", dsn)
		}
	}
}
