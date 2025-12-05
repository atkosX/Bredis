package config

var (
	Host string
	Port int
)

func init() {
	Host = "0.0.0.0"
	Port = 6379
}

var KeyLimit int = 10000
