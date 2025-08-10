package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port              string `mapstructure:"PORT"`
	DatabaseURL       string `mapstructure:"DATABASE_URL"`
	RedisURL          string `mapstructure:"REDIS_URL"`
	JwtSecret         string `mapstructure:"JWT_SECRET"`
	AccessTTLMinutes  int    `mapstructure:"JWT_ACCESS_TTL_MIN"`
	RefreshTTLHours   int    `mapstructure:"JWT_REFRESH_TTL_HOURS"`
	BcryptCost        int    `mapstructure:"BCRYPT_COST"`
	AppEnv            string `mapstructure:"APP_ENV"`
	RateLimitRequests int    `mapstructure:"RATE_LIMIT_REQUESTS"`
	RateLimitWindowS  int    `mapstructure:"RATE_LIMIT_WINDOW_SECONDS"`

	AccessTTL       time.Duration `mapstructure:"-"`
	RefreshTTL      time.Duration `mapstructure:"-"`
	RateLimitWindow time.Duration `mapstructure:"-"`
}

func LoadEnv(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "")
	viper.SetDefault("REDIS_URL", "")
	viper.SetDefault("JWT_SECRET", "dev-secret")
	viper.SetDefault("JWT_ACCESS_TTL_MIN", 15)
	viper.SetDefault("JWT_REFRESH_TTL_HOURS", 168)
	viper.SetDefault("BCRYPT_COST", 12)
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("RATE_LIMIT_REQUESTS", 10)
	viper.SetDefault("RATE_LIMIT_WINDOW_SECONDS", 60)

	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.AccessTTL = time.Duration(cfg.AccessTTLMinutes) * time.Minute
	cfg.RefreshTTL = time.Duration(cfg.RefreshTTLHours) * time.Hour
	cfg.RateLimitWindow = time.Duration(cfg.RateLimitWindowS) * time.Second

	return &cfg, nil
}
