package config

type Config struct {
	DbPort        int
	DbHost        string
	DbUser        string
	DbName        string
	DbPass        string
	NatsClusterID string
	NatsAddr      string
	ServicePort   int
	Debug         bool
}
