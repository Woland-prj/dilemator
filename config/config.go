package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -.
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		PG      PG
		S3      S3
		Metrics Metrics
		Swagger Swagger
	}

	// App -.
	App struct {
		Name       string `env:"APP_NAME,required"`
		Version    string `env:"APP_VERSION,required"`
		PassCost   int    `env:"APP_PASS_COST" envDefault:"10"`
		TgBotToken string `env:"APP_TG_BOT_TOKEN,required"`
		TgBotName  string `env:"APP_TG_BOT_NAME,required"`
	}

	// HTTP -.
	HTTP struct {
		Port                 string `env:"HTTP_PORT,required"`
		APIPrefix            string `env:"HTTP_API_PREFIX" envDefault:"/api"`
		AssetsDir            string `env:"HTTP_ASSETS_DIR" envDefault:"./web/public/assets"`
		AssetsPrefix         string `env:"HTTP_ASSETS_PREFIX" envDefault:"/assets"`
		UsePreforkMode       bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
		CookieName           string `env:"HTTP_COOKIE_NAME" envDefault:"sessionToken"`
		CookieMaxAgeDays     int    `env:"HTTP_COOKIE_MAX_AGE" envDefault:"7"`
		CookieSecure         bool   `env:"HTTP_COOKIE_SECURE" envDefault:"false"`
		CookieSameSite       string `env:"HTTP_COOKIE_SAME_SITE" envDefault:"lax"`
		CorsAllowedHeaders   string `env:"HTTP_CORS_ALLOWED_HEADERS"  envDefault:"*"`
		CorsAllowedOrigins   string `env:"HTTP_CORS_ALLOWED_ORIGINS" envDefault:"*"`
		CorsAllowCredentials bool   `env:"HTTP_CORS_ALLOW_CREDENTIALS" envDefault:"false"`
		BodyLimitMb          int    `env:"HTTP_BODY_LIMIT_MB" envDefault:"50"`
	}

	// Log -.
	Log struct {
		Env      string `env:"LOG_ENV,required"`
		Multiple bool   `env:"LOG_MULTIPLE" envDefault:"false"`
		File     string `env:"LOG_FILE" envDefault:"logs/backend.log"`
	}

	// PG -.
	PG struct {
		PoolMax  int    `env:"PG_POOL_MAX,required"`
		User     string `env:"PG_USER,required"`
		Password string `env:"PG_PASSWORD,required"`
		DbName   string `env:"PG_DB_NAME,required"`
		SslMode  bool   `env:"PG_SSL_MODE" envDefault:"false"`
		Host     string `env:"PG_HOST,required"`
		Port     string `env:"PG_PORT,required"`
	}

	S3 struct {
		AccessKey            string `env:"S3_ACCESS_KEY,required"`
		SecretKey            string `env:"S3_SECRET_KEY,required"`
		BucketName           string `env:"S3_BUCKET_NAME,required"`
		Endpoint             string `env:"S3_ENDPOINT,required"`
		Region               string `env:"S3_REGION,required"`
		PresignLifetimeHours int    `env:"S3_PRESIGN_LIFETIME_HOURS" envDefault:"1"`
		BucketDomain         string `env:"S3_BUCKET_DOMAIN,required"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}
)

// NewConfig returns base config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
