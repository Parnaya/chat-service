package database

import (
	"chat.service/configuration"
	"github.com/couchbase/gocb/v2"
)

func GetCluster(couchbaseConfig *configuration.CouchbaseConfig) (*gocb.Cluster, error) {
	return gocb.Connect(
		couchbaseConfig.ConnectString,
		gocb.ClusterOptions{
			Username: couchbaseConfig.Auth.Username,
			Password: couchbaseConfig.Auth.Password,
		})
}

func ShouldGetCluster(couchbaseConfig *configuration.CouchbaseConfig) *gocb.Cluster {
	cluster, err := GetCluster(couchbaseConfig)
	if err != nil {
		panic(err)
	}

	return cluster
}
