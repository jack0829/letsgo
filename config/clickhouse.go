package config

type ClickHouse struct {
	Debug    bool     `yaml:"Debug"`
	Addr     []string `yaml:"Addr"`
	Database string   `yaml:"Database"`
	Username string   `yaml:"Username"`
	Password string   `yaml:"Password"`
}
