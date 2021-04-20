// The config package describes the application configuration settings
package config

type Config struct {
	Addr string `config:"addr,description=Proxy address if a form of <protocol>://<host> (where protocol one of tcp or unix)"`

	LogLevel  string `config:"log_level"`
	LogFile   string `config:"log_file,description=Path to a log file to write logs (prints to stdout if not specified)"`
	LogFormat string `config:"log_format, description=Structured text or json ('text' and 'json' values respectively)"`
}

func New() Config {
	return Config{
		Addr:      "tcp://127.0.0.1:0",
		LogLevel:  "debug",
		LogFormat: "text",
	}
}
