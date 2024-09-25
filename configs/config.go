package configs

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Logger     LoggerConfig
	GRPCServer GRPCServerConfig
	HTTPServer HTTPServerConfig
	Database   DatabaseConfig
	Garantex   GarantexConfig
	Trace      TraceConfig
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Username string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `env:"DB_NAME"`
}

type LoggerConfig struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"DEBUG"`
}

type GRPCServerConfig struct {
	Host string `env:"GRPC_HOST"`
	Port string `env:"GRPC_PORT"`
}

type HTTPServerConfig struct {
	Host string `env:"HTTP_HOST"`
	Port string `env:"HTTP_PORT"`
}

func CreateAddr(host, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}

type GarantexConfig struct {
	URL     string        `env:"GARANTEX_URL"`
	Timeout time.Duration `env:"GARANTEX_TIMEOUT"`
}

type TraceConfig struct {
	Name string `env:"APP_NAME"`
	Host string `env:"TRACE_HOST"`
	Port string `env:"TRACE_PORT"`
}

func New(filenames ...string) (*Config, error) {

	var err error
	config := &Config{}

	flag.BoolFunc("name", "", func(s string) error {
		return os.Setenv("DB_NAME", s)
	})
	flag.BoolFunc("user", "", func(s string) error {
		return os.Setenv("DB_USER", s)
	})
	flag.BoolFunc("password", "", func(s string) error {
		return os.Setenv("DB_PASSWORD", s)
	})
	flag.BoolFunc("host", "", func(s string) error {
		return os.Setenv("DB_HOST", s)
	})
	flag.BoolFunc("port", "", func(s string) error {
		return os.Setenv("DB_PORT", s)
	})
	flag.Parse()

	err = godotenv.Load(filenames...)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
