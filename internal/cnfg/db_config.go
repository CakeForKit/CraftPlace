package cnfg

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type DatebaseConnConfig struct {
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

const (
	PostgresDbType = "postgres"
)

type PostgresCredentials struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	DbName   string `mapstructure:"POSTGRES_DB"`
	Port     int    `mapstructure:"POSTGRES_PORT"`
	Username string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	ReadOnly bool   `mapstructure:"DB_READONLY"`
}

// "./configs/", "app_config", "yaml"
func LoadDatebaseConnConfig(pathConfig, nameConfig, typeConfig string) (config *DatebaseConnConfig, err error) {
	config = &DatebaseConnConfig{}
	v := viper.New()
	v.AddConfigPath(pathConfig)
	v.SetConfigName(nameConfig)
	v.SetConfigType(typeConfig)
	if err = v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigRead, err)
	}
	if err = v.UnmarshalKey("datebase", config); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigRead, err)
	}

	return config, nil
}

// "./configs/", "db_config", "env"
func LoadPgCredentials(pathConfig, nameConfig, typeConfig string) (*PostgresCredentials, error) {
	viper.AddConfigPath(pathConfig)
	viper.SetConfigName(nameConfig)
	viper.SetConfigType(typeConfig)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	config := &PostgresCredentials{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	fmt.Printf("LoadPgCredentials before:\n%v\n\n", config)

	viper.AutomaticEnv()
	isReadOnly, err := strconv.ParseBool(viper.GetString("DB_READONLY"))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	pg_host := viper.GetString("POSTGRES_HOST")

	config.ReadOnly = isReadOnly
	config.Host = pg_host

	fmt.Printf("LoadPgCredentials after:\n%v\n\n", config)
	return config, nil
}
