package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

const (
	EnvironmentLocal      = "local"
	EnvironmentDebug      = "debug"
	EnvironmentProduction = "production"
)

type Config struct {
	Instance   string `yaml:"instance" env:"INSTANCE" env-required:"true"`
	Env        string `yaml:"env" env:"ENV" env-required:"true"`
	UseKafka   bool   `yaml:"use_kafka" env:"USE_KAFKA"`
	HttpServer `yaml:"http_server"`
	Storage    `yaml:"storage"`
	Secure     `yaml:"secure"`
	Kafka      `yaml:"kafka"`
	Prometheus `yaml:"prometheus"`
}

type Secure struct {
	Signature string `yaml:"secure_signature" env:"SECURE_SIGNATURE" env-required:"true"`
	Server    string `yaml:"secure_server" env:"SECURE_SERVER" env-required:"true"`
	// Проверка пароля занимает много времени. Не рекомендуется ставить таймаут меньше 1500ms
	RequestTimeout time.Duration `yaml:"secure_request_timeout" env:"SECURE_REQUEST_TIMEOUT" env-required:"true"`
	Attempts       int           `yaml:"secure_attempts" env:"SECURE_ATTEMPTS"`
	Protocol       string        `yaml:"secure_protocol" env:"SECURE_PROTOCOL"`
	Username       string        `yaml:"secure_username" env:"SECURE_USERNAME" env-required:"true"`
	// Пароль опасно хранить в открытом виде !!!
	Password string `yaml:"secure_password" env:"SECURE_PASSWORD" env-required:"true"`
}

type HttpServer struct {
	Address         string        `yaml:"address" env:"ADDRESS" env-required:"true"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-required:"true"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-required:"true"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-required:"true"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-required:"true"`
}

type Storage struct {
	DatabaseLogin              string `yaml:"database_login" env:"DATABASE_LOGIN" env-required:"true"`
	DatabasePassword           string `yaml:"database_password" env:"DATABASE_PASSWORD" env-required:"true"`
	DatabaseAddress            string `yaml:"database_address" env:"DATABASE_ADDRESS" env-required:"true"`
	DatabaseName               string `yaml:"database_name" env:"DATABASE_NAME" env-required:"true"`
	DatabaseMaxOpenConnections int    `yaml:"database_max_open_connections" env:"DATABASE_MAX_OPEN_CONNECTIONS" env-required:"true"`

	QueryTimeout time.Duration `yaml:"query_timeout" env:"QUERY_TIMEOUT" env-required:"true"`
}

type Kafka struct {
	Brokers          []string `yaml:"kafka_brokers" env:"KAFKA_BROKERS"`
	UpdatePriceTopic string   `yaml:"kafka_topic_update_price" env:"KAFKA_TOPIC_UPDATE_PRICE"`
}

type Prometheus struct {
	PrometheusPort       string `yaml:"prometheus_port" env:"PROMETHEUS_PORT"`
	PrometheusMetricsURL string `yaml:"prometheus_metrics_url" env:"PROMETHEUS_METRICS_URL"`
}

// MustLoad возвращает конфигурацию, считанную из файла, путь к которому передан из командной строки по флагу config или
// содержится в переменной окружения STORE_CONFIG_PATH. Переопределение конфигурационных значений, при необходимости,
// осуществляется посредством переменных окружения (описанных в структурах данных в этом файле).
func MustLoad() *Config {
	flag.Parse()

	var configPath = flag.String("config", "", "путь к файлу конфигурации")
	var cfg Config

	if *configPath == "" {
		*configPath = os.Getenv("STORE_CONFIG_PATH")
	}
	if *configPath == "" {
		log.Fatal("config path is not set")
	}

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", *configPath)
	}

	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
