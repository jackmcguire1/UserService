package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/jackmcguire1/UserService/pkg/utils"
	"github.com/olivere/elastic"
)

func (es *ElasticSearch) CreateIndex() (err error) {
	ctx := context.Background()
	exists, err := es.client.IndexExists(es.Name).Do(ctx)
	if err != nil || exists {
		return
	}

	_, err = es.client.CreateIndex(es.Name).Do(ctx)
	if err != nil {
		return
	}

	if es.Mapping != nil {
		mapping := map[string]interface{}{
			"properties": es.Mapping,
		}

		_, err = es.client.PutMapping().Index(es.Name).Type("_doc").IncludeTypeName(true).BodyJson(mapping).Do(ctx)
	}

	return
}

func (es *ElasticSearch) PutDoc(docID string, doc interface{}) (err error) {
	ctx := context.Background()

	idx := es.client.Index().Index(es.Name).Type("_doc").BodyJson(doc)
	if docID != "" {
		idx = idx.Id(docID)
	}
	_, err = idx.Do(ctx)

	return
}

func (es *ElasticSearch) GetDoc(docID string, v interface{}) (err error) {
	ctx := context.Background()

	idx := es.client.Get().Index(es.Name).Id(docID).Pretty(true)
	if docID != "" {
		idx = idx.Id(docID)
	}
	result, err := idx.Do(ctx)
	if err != nil {
		var elasticErr *elastic.Error
		if errors.As(err, &elasticErr) {
			if elasticErr.Status == http.StatusNotFound {
				err = fmt.Errorf("failed to find doc %s %w", err, utils.ErrNotFound)
			}
		}

		return
	}
	if result.Error != nil {
		err = fmt.Errorf("%s", utils.ToJSON(result.Error))
		return
	}

	if result.Source == nil {
		err = fmt.Errorf("source is invalid")
		return
	}
	log.WithField("raw-data", string(*result.Source)).Info("es doc")

	err = json.Unmarshal(*result.Source, &v)

	return
}

func (es *ElasticSearch) DeleteDoc(docID string) (err error) {
	ctx := context.Background()

	idx := es.client.Delete().Index(es.Name).Type("_doc").Id(docID)
	if docID != "" {
		idx = idx.Id(docID)
	}
	result, err := idx.Do(ctx)
	if err != nil {
		var elasticErr *elastic.Error
		if errors.As(err, &elasticErr) {
			if elasticErr.Status == http.StatusNotFound {
				err = fmt.Errorf("failed to find doc %s %w", err, utils.ErrNotFound)
			}
		}

		return
	}

	if result.Result != "deleted" {
		err = fmt.Errorf("failed to deleted doc %q", utils.ToJSON(result))
	}

	return
}

func (es *ElasticSearch) Query(query elastic.Query) ([]*elastic.SearchHit, error) {
	ctx := context.Background()
	idx := es.client.Search(es.Name).Query(query)

	res, err := idx.Do(ctx)
	if err != nil {
		return nil, err
	}

	return res.Hits.Hits, nil
}

func (es *ElasticSearch) DeleteIndex() (err error) {
	ctx := context.Background()
	_, err = es.client.DeleteIndex(es.Name).Do(ctx)

	return
}

func (es *ElasticSearch) PutMapping() (err error) {
	ctx := context.Background()
	mapping := map[string]interface{}{
		"properties": es.Mapping,
	}

	_, err = es.client.PutMapping().Index(es.Name).Type("_doc").IncludeTypeName(true).BodyJson(mapping).Do(ctx)
	return
}
