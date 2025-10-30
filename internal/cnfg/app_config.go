package cnfg

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	TokenSymmetricKey   string
	AccessTokenDuration time.Duration
	Port                int
	SwaggerPort         int
	ContextPath         string
}

type AppConfigFile struct {
	TokenSymmetricKey   string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration time.Duration `mapstructure:"access_token_duration"`
	SwaggerPort         int           `mapstructure:"swagger_port"`
}

// v.AddConfigPath("./configs/", "app_config", "yaml")
func LoadAppConfig(pathConfig, nameConfig, typeConfig string) (*AppConfig, error) {
	configRead := &AppConfigFile{}
	v := viper.New()
	v.AddConfigPath(pathConfig)
	v.SetConfigName(nameConfig)
	v.SetConfigType(typeConfig)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	if err := v.UnmarshalKey("app", configRead); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}

	viper.AutomaticEnv() // Автоматически считывать все переменные окружения
	appPort, err := strconv.Atoi(viper.GetString("APP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	ContextPath := viper.GetString("APP_CONTEXT_PATH")

	config := &AppConfig{
		TokenSymmetricKey:   configRead.TokenSymmetricKey,
		AccessTokenDuration: configRead.AccessTokenDuration,
		Port:                appPort,
		SwaggerPort:         configRead.SwaggerPort,
		ContextPath:         ContextPath,
	}
	fmt.Printf("LoadAppConfig:\n%v\n\n", config)
	return config, nil
}
