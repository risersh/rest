package monitoring

import (
	"crypto/tls"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/risersh/rest/conf"
)

func Setup() {
	multilog.RegisterLogger(multilog.LogMethod("console"), multilog.NewConsoleLogger(&multilog.NewConsoleLoggerArgs{
		Level:  multilog.DEBUG,
		Format: multilog.FormatText,
		FilterDropPatterns: []*string{
			multilog.PtrString("producer"), // Drop rabbitmq producer logs.
		},
	}))

	multilog.RegisterLogger(multilog.LogMethod("elasticsearch"), multilog.NewElasticsearchLogger(&multilog.NewElasticsearchLoggerArgs{
		Level: multilog.DEBUG,
		Config: elasticsearch.Config{
			Addresses: []string{conf.Config.Elasticsearch.URL},
			Username:  conf.Config.Elasticsearch.Username,
			Password:  conf.Config.Elasticsearch.Password,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		Index: "logs-rest",
		Mapping: multilog.PtrString(`
		{
			"mappings": {
				"properties": {
					"time": { "type": "date" },
					"level": { "type": "keyword" },
					"group": { "type": "keyword" },
					"message": { "type": "text" },
					"data": { "type": "object" }
				}
			}
		}`),
		FilterDropPatterns: []*string{
			multilog.PtrString("producer"), // Drop rabbitmq producer logs.
		},
	}))
}
