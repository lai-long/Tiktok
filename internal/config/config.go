package config

import (
	"Tiktok/pkg/logger"
	"os"
	"sync"

	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
	URL                string   `mapstructure:"url"`
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

// QiNiuConfig 七牛云配置
type QiNiuConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Domain    string `mapstructure:"domain"`
}

type APIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
	MapAPI  string `mapstructure:"map_api"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

// Config 总配置
type Config struct {
	EtcdAddr string      `mapstructure:"etcd_addr"`
	MySQL    MySQLConfig `mapstructure:"mysql"`
	Redis    RedisConfig `mapstructure:"redis"`
	Jwt      JwtConfig   `mapstructure:"jwt"`
	API      APIConfig   `mapstructure:"api"`
	QiNiu    QiNiuConfig `mapstructure:"qi"`
	Path     Path        `mapstructure:"filepath"`
	Mcp      MCPConfig   `mapstructure:"mcp"`
	Log      LogConfig   `mapstructure:"log"`
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
	err := v.BindEnv("mysql.host", "MYSQL_HOST")
	if err != nil {
		return nil, errors.Wrap(err, "mysql host bind env error")
	}
	err = v.BindEnv("mysql.user", "MYSQL_USER")
	if err != nil {
		return nil, errors.Wrap(err, "mysql user bind env error")
	}
	err = v.BindEnv("mysql.password", "MYSQL_PASSWORD")
	if err != nil {
		return nil, errors.Wrap(err, "mysql password bind env error")
	}
	err = v.BindEnv("redis.host", "REDIS_HOST")
	if err != nil {
		return nil, errors.Wrap(err, "redis host bind env error")
	}
	err = v.BindEnv("redis.password", "REDIS_PASSWORD")
	if err != nil {
		return nil, errors.Wrap(err, "redis password bind env error")
	}
	err = v.BindEnv("api.api_key", "OPENAI_API_KEY")
	if err != nil {
		return nil, errors.Wrap(err, "openai_api_key bind env error")
	}
	err = v.BindEnv("api.map_api", "AMAP_KEY")
	if err != nil {
		return nil, errors.Wrap(err, "amap_key bind env error")
	}
	err = v.BindEnv("jwt.access_secret", "JWT_ACCESS_SECRET")
	if err != nil {
		return nil, errors.Wrap(err, "jwt_access_secret bind env error")
	}
	err = v.BindEnv("jwt.refresh_secret", "JWT_REFRESH_SECRET")
	if err != nil {
		return nil, errors.Wrap(err, "jwt_refresh_secret bind env error")
	}
	err = v.BindEnv("qi.access_key", "QINIU_ACCESS_KEY")
	if err != nil {
		return nil, errors.Wrap(err, "qiniu_access_key bind env error")
	}
	err = v.BindEnv("qi.secret_key", "QINIU_SECRET_KEY")
	if err != nil {
		return nil, errors.Wrap(err, "qiniu_secret_key bind env error")
	}
	err = v.BindEnv("qi.bucket", "QINIU_BUCKET")
	if err != nil {
		return nil, errors.Wrap(err, "qiniu_bucket bind env error")
	}
	err = v.BindEnv("qi.domain", "QINIU_DOMAIN")
	if err != nil {
		return nil, errors.Wrap(err, "qiniu_domain bind env error")
	}
	err = v.BindEnv("etcd_addr", "ETCD_ADDR")
	if err != nil {
		return nil, errors.Wrap(err, "etcd_addr bind env error")
	}

	v.SetDefault("mysql.host", "localhost")
	v.SetDefault("mysql.port", 3306)
	v.SetDefault("mysql.user", "test")
	v.SetDefault("mysql.password", "123456")
	v.SetDefault("re.host", "localhost")
	v.SetDefault("re.port", 6379)
	v.SetDefault("re.password", "123456")

	v.SetDefault("etcd_addr", "127.0.0.1:2379")

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
			logger.Error("failed to unmarshal config", zap.Error(err))
			return
		}
		lock.Lock()
		Cfg = &newCfg
		lock.Unlock()
		logger.Info("config changed successfully")
	})
	logger.Info("config init successfully")
	return &cfg, nil
}

// LoadEnv 加载 .env 文件，优先使用 ENV_PATH 环境变量指定的路径
func LoadEnv() error {
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = ".env"
	}
	if err := godotenv.Load(envPath); err != nil {
		logger.Warn("load .env file error", zap.String("path", envPath), zap.Error(err))
	}
	return nil
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
	logger.Info("Sentinel rules loaded successfully")
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
	logger.Info("loaded flow rules", zap.Int("count", len(sentinelRules)))
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

	logger.Info("loaded circuit breaker rules", zap.Int("count", len(sentinelRules)))
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
