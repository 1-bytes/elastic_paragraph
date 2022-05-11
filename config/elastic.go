package config

import "elastic_paragraph/pkg/config"

func init() {
	config.Add("elastic", map[string]interface{}{
		"host":     config.Env("ELASTIC_HOST", "127.0.0.1:9200"),
		"username": config.Env("ELASTIC_USERNAME", ""),
		"password": config.Env("ELASTIC_PASSWORD", ""),
	})
}
