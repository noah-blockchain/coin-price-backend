package config

type Config struct {
	DbPort      int
	DbHost      string
	DbUser      string
	DbName      string
	DbPass      string
	ServicePort int
	Debug       bool
}
