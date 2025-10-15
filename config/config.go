package config

import (
	"time"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv string `mapstructure:"APP_ENV"`
	AppPort string `mapstructure:"APP_PORT"`
	DBDSN string `mapstructure:"DB_DSN"`
	RedisAddr string `mapstructure:"REDIS_ADDR"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
	JWTRefreshSecret string `mapstructure:"JWT_REFRESH_SECRET"`
	JWTAccessTTL time.Duration `mapstructure:"JWT_ACCESS_TTL"`
	JWTRefreshTTL time.Duration `mapstructure:"JWT_REFRESH_TTL"`
	PayAPIKey string `mapstructure:"PAYNETWORKS_API_KEY"`
	PayWebhookSecret string `mapstructure:"PAYNETWORKS_WEBHOOK_SECRET"`
	PayReturnURL string `mapstructure:"PAYNETWORKS_RETURN_URL"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil { return nil, err }
	if cfg.AppPort == "" { cfg.AppPort = "8080" }
	return cfg, nil
}
