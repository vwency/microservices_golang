package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func Init(env, servicePath string, cfg any) {
	viper.SetConfigName("config." + env)
	viper.SetConfigType("yaml")

	// Пути поиска
	viper.AddConfigPath(fmt.Sprintf("./config/%s", servicePath)) // ./config/service1/
	viper.AddConfigPath(".")                                     // fallback в текущую директорию

	// Чтение конфига
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Чтение переменных окружения
	viper.AutomaticEnv()

	// Декодирование в структуру
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	fmt.Printf("[CONFIG] Loaded config: %s\n", viper.ConfigFileUsed())
}
