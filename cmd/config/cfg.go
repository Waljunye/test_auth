package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"time"
)

type Config struct {
	envLoaded *options
}

type options struct {
	PrivateKey          string `env:"PRIVATE_KEY,required"`
	PublicKey           string `env:"PUBLIC_KEY,required"`
	Salt                int    `env:"SALT,required"`
	HashKey             string `env:"HASH_KEY,required"`
	RefreshTokenExpires string `env:"REFRESH_TOKEN_EXPIRES,required"`
	AccessTokenExpires  string `env:"ACCESS_TOKEN_EXPIRES,required"`

	AccessTtl      int    `env:"ACCESS_TOKEN_TTL"`
	RefreshTtl     int    `env:"REFRESH_TOKEN_TTL"`
	AppMode        string `env:"APP_MODE,required"`
	DbUrl          string `env:"DB_URL,required"`
	MigrationsPath string `env:"MIGRATIONS_PATH,required"`
	PublicApiPort  int    `env:"PUBLIC_API_PORT,required"`

	KafkaServers       string `env:"KAFKA_SERVERS"`
	KafkaConsumerGroup string `env:"KAFKA_CONSUMER_GROUP"`
}

func LoadFromEnv(fallbackFile *string) (cfg *Config, err error) {
	cfg = &Config{envLoaded: &options{}}

	if fallbackFile != nil {
		err = godotenv.Load(*fallbackFile)
	}
	*cfg.envLoaded, err = env.ParseAs[options]()
	block, _ := pem.Decode([]byte(cfg.envLoaded.PrivateKey))

	_, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	block, _ = pem.Decode([]byte(cfg.envLoaded.PublicKey))

	_, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}

	return
}

func (cfg *Config) PrivateKey() *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(cfg.envLoaded.PrivateKey))

	privateKey, _ := x509.ParsePKCS8PrivateKey(block.Bytes)

	return privateKey.(*rsa.PrivateKey)
}

func (cfg *Config) PublicKey() *rsa.PublicKey {
	block, _ := pem.Decode([]byte(cfg.envLoaded.PublicKey))

	publicKey, _ := x509.ParsePKIXPublicKey(block.Bytes)

	return publicKey.(*rsa.PublicKey)
}

func (cfg *Config) HashKey() string {
	return cfg.envLoaded.HashKey
}

func (cfg *Config) Salt() int {
	return cfg.envLoaded.Salt
}

func (cfg *Config) RefreshTokenExpires() time.Duration {
	dur, err := time.ParseDuration(cfg.envLoaded.RefreshTokenExpires)

	if err != nil {
		return time.Duration(0)
	}

	return dur
}

func (cfg *Config) AccessTokenExpires() time.Duration {
	dur, err := time.ParseDuration(cfg.envLoaded.AccessTokenExpires)

	if err != nil {
		return time.Duration(0)
	}

	return dur
}

func (cfg *Config) AppMode() string {
	return cfg.envLoaded.AppMode
}

func (cfg *Config) DbUrl() string {
	return cfg.envLoaded.DbUrl
}

func (cfg *Config) MigrationsPath() string {
	return fmt.Sprintf("file://%s", cfg.envLoaded.MigrationsPath)
}

func (cfg *Config) PublicApiPort() int {
	return cfg.envLoaded.PublicApiPort
}

func (cfg *Config) KafkaServers() string {
	return cfg.envLoaded.KafkaServers
}

func (cfg *Config) KafkaConsumerGroup() string {
	return cfg.envLoaded.KafkaConsumerGroup
}
