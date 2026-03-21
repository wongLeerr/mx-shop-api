package config

type UserInteractionConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      uint64 `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	Group     string `mapstructure:"group" json:"group"`
	Dataid    string `mapstructure:"dataid" json:"dataid"`
}

type ServerConfig struct {
	Name                string                `mapstructure:"name" json:"name"`
	Host                string                `mapstructure:"host" json:"host"`
	Port                int                   `mapstructure:"port" json:"port"`
	Tags                []string              `mapstructure:"tags" json:"tags"`
	UserInteractionConf UserInteractionConfig `mapstructure:"user_interaction_srv" json:"user_interaction_srv"`
	JWTInfo             JWTConfig             `mapstructure:"jwt" json:"jwt"`
	ConsulInfo          ConsulConfig          `mapstructure:"consul" json:"consul"`
}
