package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Gateways Gateways       `mapstructure:"gateways"`
	Services Services       `mapstructure:"services"`
}

type DatabaseConfig struct {
	Pslq       PsqlConfig   `mapstructure:"psql"`
	Redis      RedisConfig  `mapstructure:"redis"`
	RedisCache RedisConfig  `mapstructure:"redis_cache"`
	Influx     InfluxConfig `mapstructure:"influx"`
}

type PsqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Ssl      string `mapstructure:"ssl"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
	Timeout  int    `mapstructure:"timeout"`
}

type InfluxConfig struct {
	Host   string `mapstructure:"host"`
	Port   string `mapstructure:"port"`
	Bucket string `mapstructure:"bucket"`
	Org    string `mapstructure:"org"`
	Token  string `mapstructure:"token"`
}

type Gateways struct {
	Idpay    IdpayGateway    `mapstructure:"idpay"`
	Zarinpal ZarinpalGateway `mapstructure:"zarinpal"`
}

type IdpayGateway struct {
	BaseUrl  string `mapstructure:"base_url"`
	ApiToken string `mapstructure:"api_token"`
}

type ZarinpalGateway struct {
	BaseUrl  string `mapstructure:"base_url"`
	ApiToken string `mapstructure:"api_token"`
}

type Services struct {
	Alert AlertService `mapstructure:"alert"`
}

type AlertService struct {
	BaseUrl string `mapstructure:"base_url"`
}

func ViperConfig() (*viper.Viper, Config, error) {
	v := viper.New()
	v.SetEnvPrefix("AT")
	v.AutomaticEnv()

	v.SetConfigName("config")
	v.AddConfigPath("./config")
	v.AddConfigPath("/config")
	v.AddConfigPath("/app/config")
	err := v.ReadInConfig()
	if err != nil {
		return nil, Config{}, err
	}
	var config Config
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, Config{}, err
	}

	//location, err := time.LoadLocation("Asia/Tehran")
	//if err != nil {
	//	panic(err)
	//}
	//time.Local = location
	//boil.SetLocation(location)
	return v, config, nil
}
