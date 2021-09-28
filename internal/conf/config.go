package conf

type Config struct {
	ServerHost              string   `yaml:"server_host"`
	ServerPort              string   `yaml:"server_port"`
	LogLevel                string   `yaml:"log_level"`
	LogEncoding             string   `yaml:"log_encoding"`
	LoggerColor             bool     `yaml:"logger_color"`
	LoggerDisableStacktrace bool     `yaml:"logger_disable_stacktrace"`
	LoggerDevMode           bool     `yaml:"logger_dev_mode"`
	LoggerDisableCaller     bool     `yaml:"logger_disable_caller"`
	LoggerDisabledHttp      []string `yaml:"log_disabled_http"`
	Database
}

type Database struct {
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	Hostname       string `yaml:"hostname"`
	Database       string `yaml:"database"`
	Port           int    `yaml:"port"`
	MaxConnections int    `yaml:"max_connections"`
	LogQueries     bool   `yaml:"log_queries"`
	Retries			int `yaml:"retries"`
	SleepBetweenRetries string `yaml:"sleep_between_retries"`
}
