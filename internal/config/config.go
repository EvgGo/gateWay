package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	HTTPServer  HTTPServer          `yaml:"http_server"`
	Env         string              `yaml:"env" env-default:"local"`
	LogFile     string              `yaml:"log_file" env:"LOG_FILE"`
	LogLevel    string              `yaml:"log_level" env-default:"local" env:"LOG_LEVEL"`
	Options     tokenOptions        `yaml:"token_options"`
	RedisClient *RedisConfig        `yaml:"redis"`
	Mailer      string              `yaml:"mailersend_api_key" env-required:"true"`
	Hcaptcha    string              `yaml:"hcaptcha_secret" env-required:"true"`
	Auth        AuthGRPCConfig      `yaml:"auth" env-required:"true"`
	WorkSpace   WorkSpaceGRPCConfig `yaml:"work_space" env-required:"true"`
	DialConfig  DialConfig          `yaml:"dial"`
}

type tokenOptions struct {
	JWTRefreshTTL      time.Duration `yaml:"token_refresh_ttl" env-required:"true"`
	JWTAccessTTL       time.Duration `yaml:"token_access_ttl" env-required:"true"`
	JWTVerificationTTL time.Duration `yaml:"token_verification_ttl" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type AuthGRPCConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type WorkSpaceGRPCConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type RedisConfig struct {
	Host     string `yaml:"redis_host" env-required:"true"`
	Port     int    `yaml:"redis_port" env-required:"true"`
	Password string `yaml:"redis_password" env-required:"true"`
}

type DialConfig struct {
	Attempts          int
	PerAttemptTimeout time.Duration // таймаут на одну попытку Dial
	BaseBackoff       time.Duration // стартовый бэкофф
	MaxBackoff        time.Duration // верхняя граница бэкоффа
	UseTLS            bool          // если true - TLS
	ExtraOptions      []grpc.DialOption
}

func MustLoad(name string) *Config {
	path := os.Getenv(name)

	if path == "" {
		fmt.Println("Путь конфига пуст, выбран путь по умолчанию")
		path = "./configs/consumers/bgoperator.yml"
		//panic("путь в конфигу пуст")
	}

	return MustLoadPath(path)
}

type HTTPServer struct {
	Address           string        `yaml:"address" env-default:"0.0.0.0:8082"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout" `
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env-default:"2s"`
	MaxHeaderBytes    int           `yaml:"max_header_bytes" env-default:"1048576"`   // 1MB
	RequestBodyLimit  int64         `yaml:"request_body_limit" env-default:"4194304"` // 4MB
	KeepAlives        bool          `yaml:"keep_alives"    env-default:"true"`
}

func MustLoadPath(configPath string) *Config {
	// проверяем существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("конфиг файл не существует: " + configPath)
	}

	// Читаем файл полностью в память
	raw, err := os.ReadFile(configPath)
	if err != nil {
		panic("не удалось прочитать файл конфига: " + err.Error())
	}

	// Расширяем в нем все ${VARS}, os.Getenv("VARS") или "" если не задана
	expanded := os.ExpandEnv(string(raw))

	// Декодируем развернутый YAML
	var cfg Config
	if err = yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		panic("не получилось распарсить YAML: " + err.Error())
	}

	// читаем env теги
	if err = cleanenv.ReadEnv(&cfg); err != nil {
		panic("не получилось прочитать конфиг: " + err.Error())
	}

	return &cfg
}
