package config

import (
	"bytes"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const Namespace = "Virtual-Clients"

const Default = `
log:
  level: 0

server:
  addr: "0.0.0.0:8690"

redis:
  ip: "localhost"
  port: 6379
  flush: false

client:
  numqouteperminute: 20        # Number of client requests 20
  qouteperminute: 1                   # The time range of the number of quote requests 1 minutes
  clientblockminute: 1         # Minutes to Block  1 minute
  volumeqoute: 1024            # volume of qoute 1024 KB
  amountofdailyvolume: 30      # Days of volume 30 Days
  amountofvolumeblocking: 10   # Days of locking 10 Days

observability:
  addr: "0.0.0.0:6831"
  prometheus: false
  jaeger: false
  `

type (
	Config struct {
		Server        Server        `mapstructure:"server" validate:"required"`
		Log           Log           `mapstructure:"log"`
		Redis         Redis         `mapstructure:"redis" validate:"required"`
		Client        Client        `mapstructure:"client" validate:"required"`
		Observability Observability `mapstructure:"observability"`
	}

	Server struct {
		Addr string `mapstructure:"addr" validate:"required"`
	}

	Log struct {
		Level int `mapstructure:"level"`
	}

	Redis struct {
		Ip    string `mapstructure:"ip" validate:"required"`
		Port  int    `mapstructure:"port" validate:"required"`
		Flush bool   `mapstructure:"flush"`
	}

	Client struct {
		NumQoutePerMinute      int `mapstructure:"numqouteperminute" validate:"required"`
		QoutePerMinute         int `mapstructure:"qouteperminute" validate:"required"`
		ClientBlockMinute      int `mapstructure:"clientblockminute" validate:"required"`
		VolumeQoute            int `mapstructure:"volumeqoute" validate:"required"`
		AmountOfDailyVolume    int `mapstructure:"amountofdailyvolume" validate:"required"`
		AmountOfVolumeBlocking int `mapstructure:"amountofvolumeblocking" validate:"required"`
	}

	Observability struct {
		Addr       string `mapstructure:"addr""`
		Prometheus bool   `mapstructure:"prometheus"`
		Jaeger     bool   `mapstructure:"jaeger" `
	}
)

var (
	config Config
	logger = log.With().Str("Service", "Visrtual Clients").Logger()
)

func Get() *Config {
	return &config
}

func (c Config) Validate() error {
	return validator.New().Struct(c)
}

func Load(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetEnvPrefix(Namespace)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadConfig(bytes.NewReader([]byte(Default))); err != nil {
		logger.Error().Err(err).Msgf("error loading default configs")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Info().Msgf("Config file changed %s", file)
		reload(e.Name)
	})

	return reload(file)
}

func reload(file string) bool {
	err := viper.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Error().Err(err).Msgf("config file not found %s", file)
		} else {
			logger.Error().Err(err).Msgf("config file read failed %s", file)
		}
		return false
	}

	err = viper.GetViper().UnmarshalExact(&config)
	if err != nil {
		logger.Error().Err(err).Msgf("config file loaded failed %s", file)
		return false
	}

	if err = config.Validate(); err != nil {
		logger.Error().Err(err).Msgf("invalid configuration %s", file)
	}

	logger.Info().Msgf("Config file loaded %s", file)
	return true
}
