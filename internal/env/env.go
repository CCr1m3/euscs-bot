package env

import (
	"os"

	"github.com/sirupsen/logrus"
)

type dbtypeenv string

var DB struct {
	Type dbtypeenv
	Path string
}

const (
	SQLITE dbtypeenv = "sqlite"
	MYSQL  dbtypeenv = "mysql"
)

type modeenv string

var Mode modeenv // prod | test | dev

const (
	PROD modeenv = "prod"
	TEST modeenv = "test"
	DEV  modeenv = "dev"
)

type loglevelenv string

var LogLevel loglevelenv

const (
	DEBUG loglevelenv = "debug"
	INFO  loglevelenv = "info"
)

var SessionKey string

var Discord struct {
	Token             string
	GuildID           string
	OAuth2ID          string
	OAuth2Secret      string
	OAuth2RedirectURL string
}

var Twitch struct {
	OAuth2ID          string
	OAuth2Secret      string
	OAuth2RedirectURL string
}

func Init() {
	DB.Type = dbtypeenv(os.Getenv("dbtype"))
	if DB.Type == "" {
		DB.Type = SQLITE
	}
	DB.Path = os.Getenv("dbpath")
	if DB.Path == "" {
		DB.Path = "euscs.db"
	}
	Mode = modeenv(os.Getenv("mode"))
	if Mode == "" {
		Mode = PROD
	}
	LogLevel = loglevelenv(os.Getenv("loglevel"))
	if LogLevel == "" {
		LogLevel = INFO
	}
	SessionKey = os.Getenv("SESSION_KEY")
	if SessionKey == "" {
		SessionKey = "thisistotallyasupersecretsessionkey"
		logrus.Warning("a session key was not set, please use the env variable SESSION_KEY to ensure good cookie encryption")
	}
	Discord.Token = os.Getenv("discordtoken")
	Discord.GuildID = os.Getenv("discordguildid")
	Discord.OAuth2ID = os.Getenv("discordoauth2id")
	Discord.OAuth2Secret = os.Getenv("discordoauth2secret")
	Discord.OAuth2RedirectURL = os.Getenv("discordoauth2redirectURL")
	Twitch.OAuth2ID = os.Getenv("discordoauth2id")
	Twitch.OAuth2Secret = os.Getenv("twitchoauth2secret")
	Twitch.OAuth2RedirectURL = os.Getenv("twitchoauth2redirectURL")
}
