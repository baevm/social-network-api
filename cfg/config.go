package cfg

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_DSN string `mapstructure:"DB_DSN"`

	DB_CFG struct {
		MaxOpenConn     int32         `mapstructure:"maxOpenConn"`
		MaxIdleConn     int32         `mapstructure:"maxIdleConn"`
		ConnMaxLifeTime time.Duration `mapstructure:"connMaxLifeTime"`
	} `mapstructure:"postgres"`

	Redis struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		Pass string `mapstructure:"pass"`
	} `mapstructure:"redis"`

	RabbitMQ struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"rabbitmq"`

	Cloud struct {
		Name   string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
		Key    string `mapstructure:"CLOUDINARY_API_KEY"`
		Secret string `mapstructure:"CLOUDINARY_API_SECRET"`
		Folder string `mapstructure:"CLOUDINARY_UPLOAD_FOLDER"`
	} `mapstructure:"cloudinary"`

	Email struct {
		Name     string `mapstructure:"name"`
		Address  string `mapstructure:"address"`
		Password string `mapstructure:"password"`
	} `mapstructure:"email"`
}

var cfg Config

func Load(path string) error {
	readFromConfigFile()
	readFromEnvFile(path)

	viper.AutomaticEnv()

	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return nil
}

func Get() *Config {
	return &cfg
}

func readFromConfigFile() {
	viper.SetConfigFile("./cfg/cfg.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func readFromEnvFile(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("Error reading env file, %s", err)
	}
}
