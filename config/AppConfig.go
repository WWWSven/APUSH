package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// AppConfig 通知中心相关配置
type AppConfig struct {
	App      App      `yaml:"App"`      // 基本配置
	Alter    Alter    `yaml:"Alter"`    // 告警源
	Notifier Notifier `yaml:"Notifier"` // 通知目的地
}

type App struct {
	Port string `yaml:"port"`
}

type Alter struct {
	Grafana      Grafana      `yaml:"Grafana"`
	AlertManager AlertManager `yaml:"AlertManager"`
}
type Grafana struct {
	HookUri string `yaml:"hook_uri"`
}
type AlertManager struct {
	HookUri string `yaml:"hook_uri"`
}

type Notifier struct {
	DingTalk DingTalk `yaml:"DingTalk"`
	ShowDoc  ShowDoc  `yaml:"ShowDoc"`
}
type ShowDoc struct {
	HookUri string   `yaml:"hook_uri"`
	Tokens  []Tokens `yaml:"tokens"`
}
type Tokens struct {
	Name  string `yaml:"name"`
	Token string `yaml:"token"`
}
type DingTalk struct {
	HookUri string   `yaml:"hook_uri"`
	Groups  []Groups `yaml:"groups"`
}
type Groups struct {
	Name   string `yaml:"name"`
	Token  string `yaml:"token"`
	Secret string `yaml:"secret"`
}

var Config AppConfig

func Load(path string) AppConfig {
	// 加载配置
	file, _ := os.ReadFile("app.yaml")
	_ = yaml.Unmarshal(file, &Config)
	return Config
}

func GetConfig() AppConfig {
	return Config
}
