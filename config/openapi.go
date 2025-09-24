package config

type OpenAPI struct {
	Addr   string        `yaml:"Addr"`
	Client OpenAPIClient `yaml:"Client"`
	Debug  bool          `yaml:"Debug"`
}

type OpenAPIClient struct {
	ID     string `yaml:"ID"`
	Secret string `yaml:"Secret"`
}
