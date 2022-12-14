package config

import (
	"encoding/json"
	"os"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"sigs.k8s.io/yaml"
)

var validate = validator.New()

func LoadConf(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	data = []byte(os.ExpandEnv(string(data)))
	js, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := json.Unmarshal(js, c); err != nil {
		return nil, err
	}
	if err := defaults.Set(c); err != nil {
		return nil, err
	}
	if err := validate.Struct(c); err != nil {
		return nil, err
	}
	return c, nil
}

type Config struct {
	Debug     DebugConfig     `json:"debug"`
	WsProxy   WsProxyConfig   `json:"wsProxy" validate:"required"`
	Metrics   MetricsConfig   `json:"metrics"`
	Dreamkast DreamkastConfig `json:"dreamkast" validate:"required"`
	Redis     RedisConfig     `json:"redis" validate:"required"`
	Tracks    []TrackConfig   `json:"tracks" validate:"required"`
	Watcher   WatcherConfig   `json:"watcher" validate:"required"`
}

type DebugConfig struct {
	Development       bool `json:"development" default:"true"`
	Debug             bool `json:"debug"`
	DisableObsWatcher bool `json:"disableObsWatcher"`
	DisableDkWatcher  bool `json:"disableDkWatcher"`
	DisableNotifier   bool `json:"disableNotifier"`
	DisableWsProxy    bool `json:"disableWsProxy"`
}

type MetricsConfig struct {
	BindAddr string `json:"bindAddr" default:":20081"`
}

type WsProxyConfig struct {
	BindAddr string `json:"bindAddr" default:":20080"`
}

type DreamkastConfig struct {
	EventAbbr           string `json:"eventAbbr" validate:"required"`
	EndpointUrl         string `json:"endpointUrl" validate:"required"`
	Auth0Domain         string `json:"auth0Domain"`
	Auth0ClientId       string `json:"auth0ClientId"`
	Auth0ClientSecret   string `json:"auth0ClientSecret"`
	Auth0ClientAudience string `json:"auth0ClientAudience" default:"https://event.cloudnativedays.jp/"`
}

type RedisConfig struct {
	Host string `json:"host" validate:"required"`
}

type TrackConfig struct {
	DkTrackId int32       `json:"dkTrackId" validate:"required"`
	Obs       ObsConfig   `json:"obs" validate:"required"`
	Slack     SlackConfig `json:"slack" validate:"required"`
}

type ObsConfig struct {
	Host     string `json:"host" validate:"required"`
	Password string `json:"password"`
}

type SlackConfig struct {
	BotToken  string `json:"botToken" validate:"required"`
	ChannelId string `json:"channelId" validate:"required"`
}

type WatcherConfig struct {
	ObsWatcher ObsWatcherConfig `json:"obsWatcher"`
	DkWatcher  DkWatcherConfig  `json:"dkWatcher"`
}

type ObsWatcherConfig struct {
	SyncPeriodSeconds             int `json:"syncPeriodSeconds" default:"5"`
	StartPreparationPeriodSeconds int `json:"startPreparationPeriodSeconds" default:"30"`
}

type DkWatcherConfig struct {
	SyncPeriodSeconds                int `json:"syncPeriodSeconds" default:"10"`
	HowManyMinutesBeforeNotification int `json:"howManyMinutesUntilNotify" default:"5"`
}
