package config

type QCloud struct {
	Secret QCloudSecret `yaml:"QCloudSecret"`
	COS    QCloudCOS    `yaml:"COS"`
}

type QCloudSecret struct {
	ID  string `yaml:"ID"`
	Key string `yaml:"Key"`
}

type QCloudCOS struct {
	Buckets []*QCloudCosBucket `yaml:"Buckets"`
}

type QCloudCosBucket struct {
	ID     string `yaml:"ID"`
	Name   string `yaml:"Name"`
	Region string `yaml:"Region"`
	Secure bool   `yaml:"Secure"`
	Domain string `yaml:"Domain"`
}
