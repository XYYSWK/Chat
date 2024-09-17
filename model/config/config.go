package config

import (
	"time"
)

type Server struct {
	RunMode               string        `yaml:"RunMode"`               // gin 工作模式
	HttpPort              string        `yaml:"HttpPort"`              // 默认的 HTTP 监听端口号
	ReadTimeout           time.Duration `yaml:"ReadTimeout"`           // 允许读取的最大持续时间
	WriteTimeout          time.Duration `yaml:"WriteTimeout"`          // 允许写入的最大持续时间
	DefaultContextTimeout time.Duration `yaml:"DefaultContextTimeout"` // 默认上下文超时
}

type AppConfig struct {
	Name      string `yaml:"Name"`
	Version   string `yaml:"Version"`
	StartTime string `yaml:"StartTime"` // 启动时间
	MachineID int64  `yaml:"MachineID"` // 机器ID
}

type PublicConfig struct {
	Server Server      `yaml:"Server"`
	Log    LogConfig   `yaml:"Log"`
	App    AppConfig   `yaml:"App"`
	Page   PageConfig  `yaml:"Page"`
	Rules  RulesConfig `yaml:"Rules"`
	Auto   Auto        `yaml:"Auto"`
	Worker Worker      `yaml:"Worker"`
	Limit  Limit       `yaml:"Limit"`
}

type PrivateConfig struct {
	Postgresql PostgresqlConfig `yaml:"Postgresql"`
	Redis      RedisConfig      `yaml:"Redis"`
	Email      Email            `yaml:"Email"`
	Token      Token            `yaml:"Token"`
	HuaWeiOBS  HuaWeiOBS        `yaml:"HuaWeiOBS"`
	RocketMQ   RocketMQ         `yaml:"RocketMQ"`
}

type LogConfig struct {
	Level         string `yaml:"Level"`         // 日志级别
	LogSavePath   string `yaml:"LogSavePath"`   // 日志保存路径
	LowLevelFile  string `yaml:"LowLevelFile"`  // 低级别日志文件名
	LogFileExt    string `yaml:"LogFileExt"`    // 日志文件扩展名
	HighLevelFile string `yaml:"HighLevelFile"` // 高级别日志文件名
	MaxSize       int    `yaml:"MaxSize"`       // 每个日志文件的最大尺寸
	MaxAge        int    `yaml:"MaxAge"`        // 保留的最大天数
	MaxBackups    int    `yaml:"MaxBackups"`    // 保留的最大备份数量
	Compress      bool   `yaml:"Compress"`      // 是否压缩日志文件
}

type PageConfig struct {
	DefaultPageSize int32  `yaml:"DefaultPageSize"`
	MaxPageSize     int32  `yaml:"MaxPageSize"`
	PageKey         string `yaml:"PageKey"`     // URL 中 page 的关键字
	PageSizeKey     string `yaml:"PageSizeKey"` // URL 中 pagesize 的关键字
}

type RulesConfig struct {
	UsernameLenMax   int           `yaml:"UsernameLenMax"`
	UsernameLenMin   int           `yaml:"UsernameLenMin"`
	PasswordLenMax   int           `yaml:"PasswordLenMax"`
	PasswordLenMin   int           `yaml:"PasswordLenMin"`
	CodeLength       int           `yaml:"CodeLength"` // 验证码长度
	AccountNumMax    int32         `yaml:"AccountNumMax"`
	BiggestFileSize  int64         `yaml:"BiggestFileSize"`
	UserMarkDuration time.Duration `yaml:"UserMarkDuration"`
	CodeMarkDuration time.Duration `yaml:"CodeMarkDuration"`
	DefaultAvatarURL string        `yaml:"DefaultAvatarURL"`
}

// Limit 限流
type Limit struct {
	IPLimit  IPLimit  `json:"IPLimit"`
	APILimit APILimit `json:"APILimit"`
}

// IPLimit IP 限流
type IPLimit struct {
	Cap     int64 `yaml:"Cap"`     // 令牌桶容量
	GenNum  int64 `yaml:"GenNum"`  // 每次生成的令牌数量
	GenTime int64 `yaml:"GenTime"` // 生成令牌的时间间隔，即每个多长时间生成一次令牌
	Cost    int64 `yaml:"Cost"`    // 每次请求消耗的令牌数量
}

// APILimit API 限流
type APILimit struct {
	Count    int           `yaml:"Count"`    // 令牌桶容量
	Duration time.Duration `yaml:"Duration"` // 填充令牌桶的时间间隔，即每隔多长时间会填充一次令牌
	Burst    int           `yaml:"Burst"`    // 令牌桶的最大容忍峰值，即在某个时间点可以容忍的最大请求数量
}

// Auto 自动任务配置
type Auto struct {
	Retry                     Retry         `yaml:"Retry"`
	DeleteExpiredFileDuration time.Duration `yaml:"DeleteExpiredFileDuration"` // 删除过期文件的时间
}

// Worker 工作池配置
type Worker struct {
	TaskChanCapacity   int `yaml:"TaskChanCapacity"`   // 任务队列容量
	WorkerChanCapacity int `yaml:"WorkerChanCapacity"` // 工作队列容量
	WorkerNum          int `yaml:"WorkerNum"`          // 工作池数
}

// Retry 重试
type Retry struct {
	Duration time.Duration `yaml:"Duration"` // 重试的时间间隔
	MaxTimes int           `yaml:"MaxTimes"` // 最大重试次数
}

type PostgresqlConfig struct {
	DriverName string `yaml:"DriverName"`
	SourceName string `yaml:"SourceName"`
}

type RedisConfig struct {
	Address   string        `yaml:"Address"`   // Redis 服务器地址
	Password  string        `yaml:"Password"`  // 认证密码
	DB        int           `yaml:"DB"`        // Redis 数据库索引
	PoolSize  int           `yaml:"PoolSize"`  // Redis 连接池大小
	CacheTime time.Duration `yaml:"CacheTime"` // 缓存时间
}

type Email struct {
	Username string   `yaml:"Username"` // 登录邮箱的用户名
	Password string   `yaml:"Password"`
	Host     string   `yaml:"Host"`  // 邮箱服务器的主机地址
	From     string   `yaml:"From"`  // 发件人邮箱
	To       []string `yaml:"To"`    // 收件人邮箱
	Port     int      `yaml:"Port"`  // 邮箱服务器的端口号
	IsSSL    bool     `yaml:"IsSSL"` // 是否使用 SSL 加密
}

type Token struct {
	Key                  string        `yaml:"Key"`                  // 生成 token 的密钥
	AccessTokenExpire    time.Duration `yaml:"AccessTokenExpire"`    // 用户 token 的访问令牌
	RefreshTokenExpire   time.Duration `yaml:"RefreshTokenExpire"`   // 用户 token 的刷新令牌
	AccountTokenDuration time.Duration `yaml:"AccountTokenDuration"` // 账户 token 的有效期限
	AuthorizationKey     string        `yaml:"AuthorizationKey"`     // 授权密钥，用于进行授权验证
	AuthorizationType    string        `yaml:"AuthorizationType"`    // 授权类型，指定授权的具体方式或策略
}

type HuaWeiOBS struct {
	// 前两个字段推荐从 系统环境变量中获取
	AccessKeyID      string // 访问 OBS 所需的密钥 ID
	SecretAccessKey  string // 访问 OBS 所需的密钥密钥
	BucketName       string `yaml:"BucketName"`       // 存储桶名称
	BucketUrl        string `yaml:"BucketUrl"`        // 存储桶 URL
	Location         string `yaml:"Location"`         // 存储桶所在区域，必须和传入 Endpoint 中 Region 保持一致
	Endpoint         string `yaml:"Endpoint"`         // OBS 服务的 Endpoint，用与访问 OBS 的 API
	BasePath         string `yaml:"BasePath"`         // 上传文件时，文件在存储桶中的基础路径
	AvatarType       string `yaml:"FileType"`         // 头像类型
	AccountAvatarUrl string `yaml:"AccountAvatarUrl"` // 账户头像 URL（存储桶中存储账户头像的一个特定路径）
	GroupAvatarUrl   string `yaml:"GroupAvatarUrl"`   // 群组头像 URL（存储桶中存储群组头像的一个特定路径）
}

type RocketMQ struct {
	Addr string `yaml:"Addr"` // RocketMQ 服务的地址
	Port int    `yaml:"Port"` // RocketMQ 服务的端口号
}
