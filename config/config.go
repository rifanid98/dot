package config

import (
	"context"
	"sync"

	"github.com/spf13/viper"

	"dot/pkg/helper"
	"dot/pkg/util"
)

var log = util.NewLogger()

type MongoDbConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Name       string
	Production bool
	Debug      bool
}

type MySQLConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Name       string
	Production bool
	Debug      bool
}

type RedisConfig struct {
	Host           string
	Port           int
	Password       string
	ExpirationTime int64
}

type AppConfig struct {
	BaseURL                                string
	Port                                   int32
	LoggerMode                             bool
	ApiKey                                 string
	MongoDb                                *MongoDbConfig
	MySQL                                  *MySQLConfig
	Redis                                  *RedisConfig
	Email                                  *EmailConfig
	JwtSecretKey                           string
	JwtSecretKeyUserManagement             string
	JwtSecretKeyUserManagementRefreshToken string
	JwtSecretKeyMayang                     string
	JwtSecretKeys                          []string
	JwtAccessTokenExpire                   int
	JwtRefreshTokenExpire                  int
	GcpPubsub                              *GcpPubsubConfig
	ApiCall                                *ApiCallConfig
	Db                                     int
	ConfigPath                             string
}

type EmailConfig struct {
	From string
}

type GcpPubsubConfig struct {
	ProjectId         string
	Subscriber        string
	OrderingKey       string
	OrderingKeyEnable bool
}

type (
	ApiCallConfig struct {
		Xendit XenditApiCall
	}

	XenditApiCall struct {
		BaseUrl            string
		ApiKey             string
		Timeout            int
		ApiVersion20220731 string
		PathQrCodes        string
	}
)

var appConfig *AppConfig
var lock = &sync.Mutex{}

func GetConfig() *AppConfig {
	if appConfig != nil {
		return appConfig
	}

	lock.Lock()
	defer lock.Unlock()
	if appConfig != nil {
		return appConfig
	}

	appConfig = SetConfig()
	return appConfig
}

func loadConfig() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../../../../")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Info(context.Background(), "failed to read env from config file or config path. using os env variable now", err)
	}
}

func SetConfig() *AppConfig {
	loadConfig()

	var appConfig = &AppConfig{
		BaseURL:      viper.GetString("BASE_URL"),
		Port:         viper.GetInt32("PORT"),
		LoggerMode:   viper.GetBool("LOGGER_MODE"),
		ApiKey:       viper.GetString("API_KEY"),
		MongoDb:      loadMongoDbConfig(),
		MySQL:        loadMySQLCOnfig(),
		Redis:        loadRedisConfig(),
		Email:        loadEmailConfig(),
		JwtSecretKey: viper.GetString("JWT_SECRET_KEY_1"),
		JwtSecretKeys: []string{
			viper.GetString("JWT_SECRET_KEY_1"),
			viper.GetString("JWT_SECRET_KEY_2"),
			viper.GetString("JWT_SECRET_KEY_3"),
		},
		JwtAccessTokenExpire:  viper.GetInt("JWT_ACCESS_TOKEN_EXPIRE"),
		JwtRefreshTokenExpire: viper.GetInt("JWT_REFRESH_TOKEN_EXPIRE"),
		GcpPubsub:             loadGcpPubsubConfig(),
		ApiCall:               loadApiCallConfig(),
		Db:                    viper.GetInt("DATABASE_USE"),
	}

	log.Info(context.Background(), "available env variables", helper.DataToString(appConfig))

	return appConfig
}

func loadMongoDbConfig() *MongoDbConfig {
	return &MongoDbConfig{
		Host:       viper.GetString("MONGO_DB_DATABASE_HOST"),
		Port:       viper.GetInt("MONGO_DB_DATABASE_PORT"),
		Username:   viper.GetString("MONGO_DB_DATABASE_USERNAME"),
		Password:   viper.GetString("MONGO_DB_DATABASE_PASSWORD"),
		Name:       viper.GetString("MONGO_DB_DATABASE_NAME"),
		Production: viper.GetBool("MONGO_DB_DATABASE_PRODUCTION"),
		Debug:      viper.GetBool("MONGO_DB_DEBUG"),
	}
}

func loadMySQLCOnfig() *MySQLConfig {
	return &MySQLConfig{
		Host:       viper.GetString("MYSQL_DATABASE_HOST"),
		Port:       viper.GetInt("MYSQL_DATABASE_PORT"),
		Username:   viper.GetString("MYSQL_DATABASE_USERNAME"),
		Password:   viper.GetString("MYSQL_DATABASE_PASSWORD"),
		Name:       viper.GetString("MYSQL_DATABASE_NAME"),
		Production: viper.GetBool("MYSQL_DATABASE_PRODUCTION"),
		Debug:      viper.GetBool("MYSQL_DEBUG"),
	}
}

func loadRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:           viper.GetString("REDIS_HOST"),
		Port:           viper.GetInt("REDIS_PORT"),
		Password:       viper.GetString("REDIS_PASSWORD"),
		ExpirationTime: viper.GetInt64("REDIS_EXPIRATION_TIME"),
	}
}

func loadEmailConfig() *EmailConfig {
	return &EmailConfig{
		From: viper.GetString("EMAIL_FROM"),
	}
}

func loadGcpPubsubConfig() *GcpPubsubConfig {
	return &GcpPubsubConfig{
		ProjectId:         viper.GetString("GOOGLE_APPLICATION_PUBSUB_PROJECT_ID"),
		Subscriber:        viper.GetString("GOOGLE_APPLICATION_PUBSUB_SUBSCRIBER"),
		OrderingKey:       viper.GetString("GOOGLE_APPLICATION_PUBSUB_ORDERING_KEY"),
		OrderingKeyEnable: viper.GetBool("GOOGLE_APPLICATION_PUBSUB_ORDERING_KEY_ENABLE"),
	}
}

func loadApiCallConfig() *ApiCallConfig {
	return &ApiCallConfig{
		Xendit: XenditApiCall{
			BaseUrl:            viper.GetString("API_XENDIT_BASE_URL"),
			ApiKey:             viper.GetString("API_XENDIT_API_KEY"),
			Timeout:            viper.GetInt("API_XENDIT_TIMEOUT"),
			ApiVersion20220731: "2022-07-31 ",
			PathQrCodes:        "/qr_codes",
		},
	}
}
