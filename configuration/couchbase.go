package configuration

import "github.com/spf13/viper"

type CouchbaseConfig struct {
	Connect string              `yaml:"connect"`
	Auth    CouchbaseConfigAuth `yaml:"auth"`
}

type CouchbaseConfigAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func ShouldParseCouchbaseConfig() *CouchbaseConfig {
	couchbaseConfig := &CouchbaseConfig{}
	err := viper.UnmarshalKey("couchbase", couchbaseConfig)
	if err != nil {
		panic(err)
	}
	return couchbaseConfig
}
