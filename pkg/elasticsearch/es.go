package elasticsearch

import (
	"fmt"

	"github.com/olivere/elastic"
)

type ElasticSearch struct {
	Name    string
	Mapping map[string]interface{}
	client  *elastic.Client
}

type ElasticSearchClient struct {
	Name    string
	Mapping map[string]interface{}
	client  *elastic.Client
}

type ElasticSearchClientReq struct {
	Host          string
	Port          string
	SecondaryPort string
	IndexName     string
	Mapping       map[string]interface{}
}

func New(params *ElasticSearchClientReq) (es *ElasticSearch, err error) {
	client, err := elastic.NewClient(
		elastic.SetURL(
			fmt.Sprintf("http://%s:%s", params.Host, params.Port),
			fmt.Sprintf("http://%s:%s", params.Host, params.SecondaryPort),
		),
		elastic.SetSniff(false),
	)
	if err != nil {
		return
	}

	es = &ElasticSearch{
		client:  client,
		Name:    params.IndexName,
		Mapping: params.Mapping,
	}
	return
}
