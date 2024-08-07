package search

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v4"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/rs/zerolog/log"
)

type Client struct {
	OpenSearch *opensearchapi.Client
}

func New(User string, Password string, URL []string) (*Client, error) {
	opensearch_client, err := opensearchapi.NewClient(
		opensearchapi.Config{
			Client: opensearch.Config{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For testing only. Use certificate for validation.
				},
				Addresses: URL,
				Username:  User,
				Password:  Password,
			},
		},
	)
	if err != nil {
		log.Error().
			Str("id", "ERR00020001").
			Err(err).
			Msg("Could not connect to Opensearch")
		return nil, err
	}
	c := new(Client)
	c.OpenSearch = opensearch_client
	err = c.ValidateSetup()
	if err != nil {
		return c, err
	}
	return c, nil
}

// Validate the index setup in the opensearch cluster
func (c Client) ValidateSetup() error {
	ctx := context.Background()
	res, err := c.OpenSearch.Indices.Exists(ctx, opensearchapi.IndicesExistsReq{
		Indices: []string{"wi_jobs"},
	})
	if err != nil {
		if res.StatusCode != 404 {
			log.Error().
				Str("id", "ERR00020001").
				Err(err).
				Msg("Could not get index wi_jobs")
			return err
		}
		return c.CreateJobsIndex()
	}
	return nil
}
