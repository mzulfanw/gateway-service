package config

import "github.com/spf13/viper"

type Configuration struct {
	ProductServiceUrl string `mapstructure:"PRODUCT_SERVICE_URL"`
	OrderServiceUrl   string `mapstructure:"ORDER_SERVICE_URL"`
}

func LoadConfig(filename string) (*Configuration, error) {
	var envCfg Configuration

	viper.AddConfigPath(".")
	viper.SetConfigName(filename)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&envCfg); err != nil {
		return nil, err
	}

	return &envCfg, nil
}
