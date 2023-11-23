package config

import (
	"embed"
	"os"

	_ "embed"

	"gopkg.in/yaml.v2"

	libconfig "github.com/pecolynx/golang-webapi-boilerplate/lib/config"
	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
)

type AppConfig struct {
	Name        string `yaml:"name" validate:"required"`
	HTTPPort    int    `yaml:"httpPort" validate:"required"`
	MetricsPort int    `yaml:"metricsPort" validate:"required"`
}

type AuthConfig struct {
	SigningKey          string `yaml:"signingKey"`
	AccessTokenTTLMin   int    `yaml:"accessTokenTtlMin" validate:"gte=1"`
	RefreshTokenTTLHour int    `yaml:"refreshTokenTtlHour" validate:"gte=1"`
}

type ShutdownConfig struct {
	TimeSec1 int `yaml:"timeSec1" validate:"gte=1"`
	TimeSec2 int `yaml:"timeSec2" validate:"gte=1"`
}

type DebugConfig struct {
	Gin  bool `yaml:"ginMode"`
	Wait bool `yaml:"wait"`
}

type Config struct {
	App      *AppConfig               `yaml:"app" validate:"required"`
	DB       *libconfig.DBConfig      `yaml:"db" validate:"required"`
	Auth     *AuthConfig              `yaml:"auth" validate:"required"`
	Trace    *libconfig.TraceConfig   `yaml:"trace" validate:"required"`
	CORS     *libconfig.CORSConfig    `yaml:"cors" validate:"required"`
	Shutdown *ShutdownConfig          `yaml:"shutdown" validate:"required"`
	Log      *libconfig.LogConfig     `yaml:"log" validate:"required"`
	Swagger  *libconfig.SwaggerConfig `yaml:"swagger" validate:"required"`
	Debug    *DebugConfig             `yaml:"debug"`
}

//go:embed local.yml
//go:embed production.yml
var config embed.FS

func LoadConfig(env string) (*Config, error) {
	filename := env + ".yml"
	confContent, err := config.ReadFile(filename)
	if err != nil {
		return nil, liberrors.Errorf("config.ReadFile. filename: %s, err: %w", filename, err)
	}

	confContent = []byte(os.ExpandEnv(string(confContent)))
	conf := &Config{}
	if err := yaml.Unmarshal(confContent, conf); err != nil {
		return nil, liberrors.Errorf("yaml.Unmarshal. filename: %s, err: %w", filename, err)
	}

	if err := libdomain.Validator.Struct(conf); err != nil {
		return nil, liberrors.Errorf("Validator.Structl. filename: %s, err: %w", filename, err)
	}

	return conf, nil
}
