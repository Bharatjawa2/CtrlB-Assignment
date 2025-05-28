package config

import (
	"flag"
	"log"
	"os"
	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct{
	Addr string `yaml:"address" env-required:"true"`
}

type AdminConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type Config struct{
	Env string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	JWTSecret string `yaml:"jwt_secret"`
	Admin AdminConfig  `yaml:"admin" env-required:"true"`
}

func MustLoad() *Config{
	var configPath string
	configPath=os.Getenv("CONFIG_PATH")
	if configPath==""{
		flags:=flag.String("config","","path to configuration file")
		flag.Parse()

		configPath=*flags

		if configPath==""{
			log.Fatal("Config Path is not Set")
		}
	}

	if _, err:=os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("config path does not exist: %s",configPath)
	}

	var cfg Config

	err:=cleanenv.ReadConfig(configPath,&cfg)
	if err!=nil{
		log.Fatalf("Cannot read config file %s",err.Error())
	}

	return &cfg
}