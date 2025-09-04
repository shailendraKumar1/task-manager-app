package config

import (
	"task-manager-app/constants"
	"task-manager-app/utils"
	"os"
	"strings"
)

type applicationConfig struct {
	AppName          string
	AppVersion       string
	AppHost          string
	AppPort          string
	KafkaHosts       []string
	KafkaGroupId     string
	UserAppBaseUri   string
	ErrorCodes       string
	KafkaRetryTopic  string
	RedisPort        string
	RedisPassword    string
	RedisUsername    string
	RedisDb          int
	RedisTimeout     int
	RedisPoolTimeout int
	RedisEndpoint    string
	KafkaUsername    string
	KafkaPassword    string
	PostgresAddress  string
	Username         string
	Password         string
	DbName           string
	PostgresHost     string
	PostgresPort     string
	MaxDbConnections int
}

var (
	ApplicationConfig = &applicationConfig{}
)

func (a applicationConfig) SetAwsSecretValues() {

	ApplicationConfig = &applicationConfig{
		AppName:          os.Getenv(constants.AppName),
		AppVersion:       os.Getenv(constants.AppVersion),
		AppHost:          os.Getenv(constants.HOST),
		AppPort:          os.Getenv(constants.PORT),
		KafkaHosts:       strings.Split(os.Getenv(constants.KafkaHosts), ","),
		KafkaGroupId:     os.Getenv(constants.KafkaGroupId),
		UserAppBaseUri:   os.Getenv(constants.UserAppBaseUri),
		KafkaRetryTopic:  os.Getenv(constants.KafkaRetryTopic),
		KafkaUsername:    os.Getenv(constants.KafkaUsername),
		KafkaPassword:    os.Getenv(constants.KafkaPassword),
		PostgresAddress:  os.Getenv(constants.PostgresAddress),
		Username:         os.Getenv(constants.PostgresUsername),
		Password:         os.Getenv(constants.PostgresPassword),
		DbName:           os.Getenv(constants.PostgresDbName),
		PostgresHost:     os.Getenv(constants.PostgresHost),
		PostgresPort:     os.Getenv(constants.PostgresPort),
		MaxDbConnections: utils.TaskManagerUtils.ParseStringToInt(os.Getenv(constants.MaxDbConnections)),
	}

}
