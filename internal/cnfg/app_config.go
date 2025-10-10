package cnfg

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	TokenSymmetricKey   string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration time.Duration `mapstructure:"access_token_duration"`
	Port                int           `mapstructure:"port"`
}

// v.AddConfigPath("./configs/", "app_config", "yaml")
func LoadAppConfig(pathConfig, nameConfig, typeConfig string) (config *AppConfig, err error) {
	config = &AppConfig{}
	v := viper.New()
	v.AddConfigPath(pathConfig)
	v.SetConfigName(nameConfig)
	v.SetConfigType(typeConfig)
	if err = v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	if err = v.UnmarshalKey("app", config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}

	return config, nil
}
