package configs

import "github.com/spf13/viper"

type conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
}

func LoadConfig(path string) (*conf, error) {
	// Inicializar a estrutura de configuração
	cfg := &conf{
		// Valores padrão
		DBDriver:          "mysql",
		DBHost:            "mysql",
		DBPort:            "3306",
		DBUser:            "root",
		DBPassword:        "root",
		DBName:            "orders",
		WebServerPort:     "8080",
		GRPCServerPort:    "50051",
		GraphQLServerPort: "8081",
	}

	// Configurar o viper para ler variáveis de ambiente
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Tenta ler o arquivo .env, mas não falha se ele não existir
	_ = viper.ReadInConfig()

	// Sobrescrever os valores padrão com as variáveis de ambiente
	if envDBDriver := viper.GetString("DB_DRIVER"); envDBDriver != "" {
		cfg.DBDriver = envDBDriver
	}
	if envDBHost := viper.GetString("DB_HOST"); envDBHost != "" {
		cfg.DBHost = envDBHost
	}
	if envDBPort := viper.GetString("DB_PORT"); envDBPort != "" {
		cfg.DBPort = envDBPort
	}
	if envDBUser := viper.GetString("DB_USER"); envDBUser != "" {
		cfg.DBUser = envDBUser
	}
	if envDBPassword := viper.GetString("DB_PASSWORD"); envDBPassword != "" {
		cfg.DBPassword = envDBPassword
	}
	if envDBName := viper.GetString("DB_NAME"); envDBName != "" {
		cfg.DBName = envDBName
	}
	if envWebServerPort := viper.GetString("WEB_SERVER_PORT"); envWebServerPort != "" {
		cfg.WebServerPort = envWebServerPort
	}
	if envGRPCServerPort := viper.GetString("GRPC_SERVER_PORT"); envGRPCServerPort != "" {
		cfg.GRPCServerPort = envGRPCServerPort
	}
	if envGraphQLServerPort := viper.GetString("GRAPHQL_SERVER_PORT"); envGraphQLServerPort != "" {
		cfg.GraphQLServerPort = envGraphQLServerPort
	}

	return cfg, nil
}
