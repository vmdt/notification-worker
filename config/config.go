package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	echoserver "github.com/vmdt/notification-worker/pkg/echo"
	mailer "github.com/vmdt/notification-worker/pkg/email"
	"github.com/vmdt/notification-worker/pkg/logger"
	"github.com/vmdt/notification-worker/pkg/mongodb"
	"github.com/vmdt/notification-worker/pkg/rabbitmq"
	redis2 "github.com/vmdt/notification-worker/pkg/redis"
)

var configPath string

type Config struct {
	Logger   *logger.LoggerConfig     `mapstructure:"logger"`
	MongoDb  *mongodb.MongoDbOptions  `mapstructure:"mongodb"`
	Rabbitmq *rabbitmq.RabbitMQConfig `mapstructure:"rabbitmq"`
	Echo     *echoserver.EchoConfig   `mapstructure:"echo"`
	Mailer   *mailer.MailerConfig     `mapstructure:"mailer"`
	Redis    *redis2.RedisOptions     `mapstructure:"redis"`
}

func InitConfig() (
	*Config,
	*logger.LoggerConfig,
	*mongodb.MongoDbOptions,
	*rabbitmq.RabbitMQConfig,
	*echoserver.EchoConfig,
	*mailer.MailerConfig,
	*redis2.RedisOptions,
	error,
) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	if configPath == "" {
		configPathFromEnv := os.Getenv("CONFIG_PATH")
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			//https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
			//https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
			d, err := dirname()
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, err
			}

			configPath = d
		}
	}

	cfg := &Config{}
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return cfg, cfg.Logger, cfg.MongoDb, cfg.Rabbitmq, cfg.Echo, cfg.Mailer, cfg.Redis, nil

}

func filename() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

func dirname() (string, error) {
	filename, err := filename()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filename), nil
}
