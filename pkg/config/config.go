package config

import (
	"log"
	"sync"

	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Path 常用路径
type Path struct {
	VideoPath  string `mapstructure:"video_path"`
	AvatarPath string `mapstructure:"avatar_path"`
	EnvPath    string `mapstructure:"env_path"`
}
type MCPConfig struct {
	Clients           []McpClientConfig `mapstructure:"clients"`
	ToolManagerConfig ToolManagerConfig `mapstructure:"tool_manager"`
}
type ToolManagerConfig struct {
	MaxTime  int64 `mapstructure:"tool_execution_timeout"`
	MaxDepth int64 `mapstructure:"max_agent_depth"`
}
type McpClientConfig struct {
	ID                 string   `mapstructure:"id"`
	Name               string   `mapstructure:"name"`
	ConnectionType     string   `mapstructure:"connection_type"`
	Command            string   `mapstructure:"command"`
	Args               []string `mapstructure:"args"`
	ToolsToExecute     []string `mapstructure:"tools_to_execute"`
	ToolsToAutoExecute []string `mapstructure:"tools_to_auto_execute"`
}

// MySQLConfig Mysql配置
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

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

// JwtConfig JWT密钥
type JwtConfig struct {
	AccessSecret  string `mapstructure:"access_secret"`
	RefreshSecret string `mapstructure:"refresh_secret"`
}

// APIConfig a
type APIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
	MapAPI  string `mapstructure:"map_api"`
}

// Config 总配置
type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	Jwt   JwtConfig   `mapstructure:"jwt"`
	API   APIConfig   `mapstructure:"api"`
	Path  Path        `mapstructure:"filepath"`
	Mcp   MCPConfig   `mapstructure:"mcp"`
}

// Cfg 调用配置
var Cfg *Config
var lock sync.RWMutex

// Load 加载配置
func Load(confPath []string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	for _, p := range confPath {
		v.AddConfigPath(p)
	}
	v.AutomaticEnv()
	err := v.BindEnv("mysql.password", "MYSQL_PASSWORD")
	if err != nil {
		return nil, errors.Wrap(err, "mysql password bind env error")
	}
	err = v.BindEnv("redis.password", "REDIS_PASSWORD")
	if err != nil {
		return nil, errors.Wrap(err, "redis password bind env error")
	}
	err = v.BindEnv("api.api_key", "OPENAI_API_KEY")
	if err != nil {
		return nil, errors.Wrap(err, "openai_api_key bind env error")
	}
	err = v.BindEnv("jwt.access_secret", "JWT_ACCESS_SECRET")
	if err != nil {
		return nil, errors.Wrap(err, "jwt_access_secret bind env error")
	}
	err = v.BindEnv("jwt.refresh_secret", "JWT_REFRESH_SECRET")
	if err != nil {
		return nil, errors.Wrap(err, "jwt_refresh_secret bind env error")
	}

	v.SetDefault("mysql.host", "localhost")
	v.SetDefault("mysql.port", 3306)
	v.SetDefault("mysql.user", "test")
	v.SetDefault("mysql.password", "123456")
	v.SetDefault("re.host", "localhost")
	v.SetDefault("re.port", 6379)
	v.SetDefault("re.password", "123456")

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}
	lock.Lock()
	Cfg = &cfg
	lock.Unlock()

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		var newCfg Config
		if err := v.Unmarshal(&newCfg); err != nil {
			log.Println("failed to unmarshal config")
			return
		}
		lock.Lock()
		Cfg = &newCfg
		lock.Unlock()
		log.Println("config changed successfully")
	})
	log.Println("config init successfully")
	return &cfg, nil
}

type FlowRuleConfig struct {
	Resource          string `mapstructure:"resource"`
	Threshold         int    `mapstructure:"threshold"`
	ControlBehavior   int    `mapstructure:"controlBehavior"`
	MaxQueueingTimeMs uint32 `mapstructure:"maxQueueingTimeMs"`
	WarmUpPeriodSec   uint32 `mapstructure:"warmUpPeriodSec"`
}

type CircuitBreakerConfig struct {
	Resource         string  `mapstructure:"resource"`
	Threshold        float64 `mapstructure:"threshold"`
	Strategy         string  `mapstructure:"strategy"`
	RetryTimeoutMs   uint32  `mapstructure:"retryTimeoutMs"`
	MinRequestAmount uint64  `mapstructure:"minRequestAmount"`
	MaxAllowedRtMs   uint64  `mapstructure:"maxAllowedRtMs"`
}

type RulesConfig struct {
	Flow           []FlowRuleConfig       `mapstructure:"flow"`
	CircuitBreaker []CircuitBreakerConfig `mapstructure:"circuitBreaker"`
}

func LoadRules(configPath []string) error {
	v := viper.New()
	v.SetConfigName("sentinel")
	v.SetConfigType("yaml")
	for _, p := range configPath {
		v.AddConfigPath(p)
	}
	if err := v.ReadInConfig(); err != nil {
		return errors.Wrap(err, "failed to read config")
	}
	var cfg RulesConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}
	if err := loadFlowRules(cfg.Flow); err != nil {
		return errors.Wrap(err, "failed to load flow rules")
	}
	if err := loadCircuitBreakerRules(cfg.CircuitBreaker); err != nil {
		return errors.Wrap(err, "failed to load circuit breaker rules")
	}
	log.Println("Sentinel rules loaded successfully")
	return nil
}

func loadFlowRules(flowRules []FlowRuleConfig) error {
	if len(flowRules) == 0 {
		return nil
	}
	sentinelRules := make([]*flow.Rule, 0, len(flowRules))
	for _, r := range flowRules {
		rule := &flow.Rule{
			Resource:        r.Resource,
			Threshold:       float64(r.Threshold),
			ControlBehavior: flow.ControlBehavior(r.ControlBehavior),
		}
		switch flow.ControlBehavior(r.ControlBehavior) {
		case flow.Throttling:
			rule.MaxQueueingTimeMs = r.MaxQueueingTimeMs
		case flow.Reject:
			if r.WarmUpPeriodSec > 0 {
				rule.TokenCalculateStrategy = flow.WarmUp
				rule.WarmUpPeriodSec = r.WarmUpPeriodSec
			}
		}
		sentinelRules = append(sentinelRules, rule)
	}
	_, err := flow.LoadRules(sentinelRules)
	if err != nil {
		return errors.Wrap(err, "failed to load flow rules")
	}
	log.Printf("loaded %d flow rules", len(sentinelRules))
	return nil
}

func loadCircuitBreakerRules(cbRules []CircuitBreakerConfig) error {
	if len(cbRules) == 0 {
		return nil
	}
	sentinelRules := make([]*circuitbreaker.Rule, 0, len(cbRules))
	for _, r := range cbRules {
		rule := &circuitbreaker.Rule{
			Resource:         r.Resource,
			Threshold:        r.Threshold,
			Strategy:         parseStrategy(r.Strategy),
			RetryTimeoutMs:   r.RetryTimeoutMs,
			MinRequestAmount: r.MinRequestAmount,
			MaxAllowedRtMs:   r.MaxAllowedRtMs,
		}
		sentinelRules = append(sentinelRules, rule)
	}
	_, err := circuitbreaker.LoadRules(sentinelRules)
	if err != nil {
		return errors.Wrap(err, "failed to load circuit breaker rules")
	}

	log.Printf("loaded %d circuit breaker rules", len(sentinelRules))
	return nil
}

func parseStrategy(s string) circuitbreaker.Strategy {
	switch s {
	case "slowRequestRatio":
		return circuitbreaker.SlowRequestRatio
	case "errorRatio":
		return circuitbreaker.ErrorRatio
	case "errorCount":
		return circuitbreaker.ErrorCount
	default:
		return circuitbreaker.SlowRequestRatio
	}
}
