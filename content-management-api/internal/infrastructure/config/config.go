package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

type AppConfig struct {
	Database struct {
		Host          string `mapstructure:"host"`
		Port          int    `mapstructure:"port"`
		User          string `mapstructure:"user"`
		Password      string `mapstructure:"password"`
		Name          string `mapstructure:"name"`
		MaxConnection int32  `mapstructure:"max_connections"`
		MinConnection int32  `mapstructure:"min_connections"`
		MaxIdleTime   int32  `mapstructure:"max_idle_time"`
		MaxLifeTime   int32  `mapstructure:"max_conn_lifetime"`
	} `mapstructure:"database"`

	App struct {
		Port        string `mapstructure:"port"`
		MaxFileSize int64  `mapstructure:"max_file_size"`
	} `mapstructure:"app"`

	Aws struct {
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
		Region          string `mapstructure:"region"`
		S3              struct {
			Bucket     string `mapstructure:"bucket"`
			BaseFolder string `mapstructure:"base_folder"`
		} `mapstructure:"s3"`
	} `mapstructure:"aws"`

	Kafka struct {
		BootstrapServers   string `mapstructure:"bootstrap_servers"`
		GroupId            string `mapstructure:"group_id"`
		VideoUploadedTopic string `mapstructure:"video_uploaded_topic"`
	} `mapstructure:"kafka"`
}

func LoadConfig() AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	for _, key := range viper.AllKeys() {
		value := viper.GetString(key)
		viper.Set(key, os.ExpandEnv(value))
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	return config
}
