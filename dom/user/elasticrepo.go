package user

import (
	"encoding/json"

	"github.com/apex/log"
	"github.com/jackmcguire1/UserService/pkg/elasticsearch"
	"github.com/jackmcguire1/UserService/pkg/utils"
	"github.com/olivere/elastic/v7"
)

type ElasticSearchRepository struct {
	BaseRepository

	Host          string
	Port          string
	SecondPort    string
	UserIndexName string
}

type ElasticSearchParams struct {
	Host          string
	Port          string
	SecondPort    string
	UserIndexName string
}

func NewElasticRepo(params *ElasticSearchParams) *ElasticSearchRepository {
	es := &ElasticSearchRepository{
		Host:          params.Host,
		Port:          params.Port,
		SecondPort:    params.SecondPort,
		UserIndexName: params.UserIndexName,
	}

	client, err := es.getEsClient()
	if err != nil {
		log.
			WithField("elastic-host", params.Host).
			WithField("elastic-port", params.Port).
			WithField("elastic-second-port", params.SecondPort).
			WithError(err).
			Fatal("failed to init elastic search repo")
	}

	err = client.CreateIndex()
	if err != nil {
		log.
			WithError(err).
			WithField("index-name", es.UserIndexName).
			Fatal("failed to create elastic search repo")
	}

	return es
}

func (repo *ElasticSearchRepository) GetUser(userId string) (u *User, err error) {
	es, err := repo.getEsClient()
	if err != nil {
		return nil, err
	}

	err = es.GetDoc(userId, &u)
	return
}

func (repo *ElasticSearchRepository) GetUsersByCountry(cc string) (users []*User, err error) {
	es, err := repo.getEsClient()
	if err != nil {
		return nil, err
	}

	q := elastic.NewMatchQuery("CountryCode", cc)
	res, err := es.Query(q)
	if err != nil {
		return
	}

	users = []*User{}
	for _, hit := range res {
		if hit.Source == nil {
			log.
				WithField("hit", utils.ToJSON(hit)).
				Error("document does not have source")

			continue
		}

		var user *User
		err = json.Unmarshal(*hit.Source, &user)
		if err != nil {
			return
		}

		users = append(users, user)
	}

	return users, err
}

func (repo *ElasticSearchRepository) PutUser(u *User) error {
	es, err := repo.getEsClient()
	if err != nil {
		return err
	}

	return es.PutDoc(u.ID, u)
}

func (repo *ElasticSearchRepository) DeleteUser(id string) error {
	client, err := repo.getEsClient()
	if err != nil {
		return err
	}

	return client.DeleteDoc(id)
}

func (repo *ElasticSearchRepository) getEsClient() (*elasticsearch.ElasticSearch, error) {
	return elasticsearch.New(&elasticsearch.ElasticSearchClientReq{
		Host:          repo.Host,
		Port:          repo.Port,
		SecondaryPort: repo.SecondPort,
		IndexName:     repo.UserIndexName,
		Mapping:       nil,
	})
}
