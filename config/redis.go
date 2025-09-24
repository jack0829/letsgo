package config

type Redis struct {
	Addr        string `yaml:"Addr"`
	User        string `yaml:"User"`
	Password    string `yaml:"Password"`
	DB          int    `yaml:"DB"`
	PoolSize    int    `yaml:"PoolSize"`
	MinIdleConn int    `yaml:"MinIdleConn"`
}
