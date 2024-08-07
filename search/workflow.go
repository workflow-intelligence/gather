package search

import (
	"context"
	"fmt"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/rs/zerolog/log"
)

func (c *Client) CreateWorkflowIndex(Organization string, Repository string, RunId int64) (string, error) {
	index := fmt.Sprintf("github_%v_%v_%v", Organization, Repository, RunId)
	ctx := context.Background()
	res, err := c.OpenSearch.Indices.Create(ctx, opensearchapi.IndicesCreateReq{
		Index: index,
	})
	if err != nil {
		log.Error().
			Str("id", "ERR00021001").
			Str("index", index).
			Err(err).
			Msg("Could not create index")
		return "", err
	}
	return res.Index, nil
}
