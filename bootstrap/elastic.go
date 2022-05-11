package bootstrap

import (
	"elastic_paragraph/pkg/config"
	es "elastic_paragraph/pkg/elastic"
	"github.com/olivere/elastic/v7"
	"time"
)

// SetupElastic 初始化 Elastic.
func SetupElastic() {
	es.Options = []elastic.ClientOptionFunc{
		elastic.SetURL(config.GetString("elastic.host")),
		elastic.SetBasicAuth(
			config.GetString("elastic.username"),
			config.GetString("elastic.password"),
		),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(5 * time.Second),
	}
}
