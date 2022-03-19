package pkg

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

//  Config stores server config params.
type Config struct {
	ServerAddress  string `env:"RUN_ADDRESS"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseDSN    string `env:"DATABASE_DSN"`
}

//  Default config params.
const (
	defDatabaseDSN    = ""
	defServerAddress  = ""
	defAccrualAddress = ""
)

//  NewConfig inits new config.
//  Reads flag params over default params, then redefines  with environment params.
func NewConfig() (*Config, error) {
	cfg := Config{}
	cfg.readFlagConfig()

	if err := cfg.readEnvConfig(); err != nil {
		return nil, fmt.Errorf("error config initializatioin: %w", err)
	}

	// TODO validation

	return &cfg, nil
}

//  readFlagConfig reads flag params over default params.
func (c *Config) readFlagConfig() {
	flag.StringVar(&c.DatabaseDSN, "d", defDatabaseDSN, "адрес подключения к базе данных")
	flag.StringVar(&c.ServerAddress, "a", defServerAddress, "адрес и порт запуска сервиса")
	flag.StringVar(&c.AccrualAddress, "r", defAccrualAddress, "адрес системы расчёта начислений")
	flag.Parse()
}

//  readEnvConfig redefines config params with environment params.
func (c *Config) readEnvConfig() error {
	envConfig := &Config{}

	if err := env.Parse(envConfig); err != nil {
		return fmt.Errorf("ошибка чтения переменных окружения:%w", err)
	}

	if envConfig.DatabaseDSN != "" {
		c.DatabaseDSN = envConfig.DatabaseDSN
	}
	if envConfig.ServerAddress != "" {
		c.ServerAddress = envConfig.ServerAddress
	}
	if envConfig.AccrualAddress != "" {
		c.DatabaseDSN = envConfig.AccrualAddress
	}

	return nil
}
