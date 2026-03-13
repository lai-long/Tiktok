package conf

import (
	"github.com/spf13/viper"
)

type MySQLConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Database  string `mapstructure:"database"`
	Charset   string `mapstructure:"charset"`
	ParseTime bool   `mapstructure:"parse_time"`
	Loc       string `mapstructure:"loc"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}
type JwtConfig struct {
	Secret string `mapstructure:"secret"`
}
type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"re"`
	Jwt   JwtConfig   `mapstructure:"jwt"`
}

var Cfg *Config

func Load(confPath []string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	for _, p := range confPath {
		v.AddConfigPath(p)
	}
	v.SetDefault("mysql.host", "localhost")
	v.SetDefault("mysql.port", 3306)
	v.SetDefault("mysql.user", "root")
	v.SetDefault("mysql.password", "root")
	v.SetDefault("re.host", "localhost")
	v.SetDefault("re.port", 6379)
	v.SetDefault("re.password", "")
	v.SetDefault("jwt.secret", "secret")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	Cfg = &cfg
	return &cfg, nil
}
