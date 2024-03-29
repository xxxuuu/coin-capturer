package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"strings"
)

type Config struct {
	NodeAddress     string   `koanf:"node_address"`
	MonitoredWallet []string `koanf:"monitored_wallet"`
	DingtalkToken   string   `koanf:"dingtalk_token"`
	OklinkToken     string   `koanf:"oklink_token"`
	LowerLimitValue string   `koanf:"lower_limit_value"`
}

func InitConfig() (*Config, error) {
	k := koanf.New(".")
	_ = k.Load(file.Provider("config.yaml"), yaml.Parser())
	_ = k.Load(file.Provider("conf/config.yaml"), yaml.Parser())
	_ = k.Load(env.ProviderWithValue("", ".", func(s string, v string) (string, interface{}) {
		key := strings.ToLower(s)
		if strings.Contains(v, ",") {
			return key, strings.Split(v, ",")
		}
		return key, v
	}), nil)
	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
