package search

import (
	"context"
	"encoding/json"
	"fmt"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func (c *Client) CreateJobsIndex() error {

	ctx := context.Background()
	settings := strings.NewReader(`{
		"settings": {
			"index": {
				"number_of_shards": 1,
				"number_of_replicas": 1
			}
		},		
		"mappings": {
			"properties": {
				"@timestamp": {
					"type": "date",
					"format": "rfc3339_lenient"
				},
				"ci": {
					"type": "text"
				},
				"organization": {
					"type": "text"
				},
				"repository": {
					"type": "text"
				},
				"runid": {
					"type": "integer"
				}
			}
		}
	}`)
	_, err := c.OpenSearch.Indices.Create(ctx, opensearchapi.IndicesCreateReq{
		Index: "wi_jobs",
		Body:  settings,
	})
	if err != nil {
		log.Error().
			Str("id", "ERR00032000").
			Str("index", "wi_jobs").
			Err(err).
			Msg("Could not create index")
	}
	return err
}

func (c *Client) AddJob(CI string, Organization string, Repository string, RunId int64) (string, error) {
	ctx := context.Background()
	jobson := `{
		"@timestamp": "%v",
		"ci": %v,
		"organization": "%v",
		"repository": "%v",
		"runid": %v
		}`
	schedule := time.Now().Add(time.Minute * time.Duration(5)).Format(time.RFC3339)
	json := fmt.Sprintf(jobson, schedule, Organization, Repository, RunId)
	document := strings.NewReader(json)
	log.Debug().Str("json", json).Msg("Create job")
	id := fmt.Sprintf("%v_%v_%v_%v", CI, Organization, Repository, RunId)
	_, err := c.OpenSearch.Document.Create(ctx, opensearchapi.DocumentCreateReq{
		Index:      "wi_jobs",
		DocumentID: id,
		Body:       document,
	})
	if err != nil {
		log.Error().
			Str("id", "ERR00032010").
			Str("index", "wi_jobs").
			Str("ci", CI).
			Str("organization", Organization).
			Str("repositiory", Repository).
			Int64("runid", RunId).
			Err(err).
			Msg("Could not add job")
	}
	return id, err
}

func (c *Client) DeleteJob(DocumentId string) error {
	ctx := context.Background()
	_, err := c.OpenSearch.Document.Delete(ctx, opensearchapi.DocumentDeleteReq{
		Index:      "wi_jobs",
		DocumentID: DocumentId,
	})
	if err != nil {
		log.Error().
			Str("index", "wi_jobs").
			Str("DocumentID", DocumentId).
			Msg("Could not delete Job")
	}
	return err
}

type Job struct {
	Timestamp    time.Time `json:"@timestamp"`
	CI           string    `json:"ci"`
	Organization string    `json:"organization"`
	Repository   string    `json:"repository"`
	RunId        int64     `json:"runid"`
}

type JobList map[string]Job

func (c *Client) PendingJobs() (JobList, error) {
	list := make(JobList)
	ctx := context.Background()
	search := strings.NewReader(fmt.Sprintf(`{
		"query": {
			"range": {
				"@timestamp": {
					"lte": "%v"
				}
			}
		}
	}`, time.Now().Format(time.RFC3339)))
	res, err := c.OpenSearch.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{"wi_jobs"},
		Body:    search,
	})
	if err != nil {
		log.Error().
			Str("id", "ERR00032020").
			Err(err).
			Msg("Could not search for jobs")
		return nil, err
	}
	for _, hit := range res.Hits.Hits {
		var job Job
		err := json.Unmarshal(hit.Source, &job)
		if err != nil {
			log.Error().
				Str("id", "ERR00032021").
				Str("json", string(hit.Source)).
				Err(err).
				Msg("Could not unmarshal job")
			return nil, err
		}
		key := fmt.Sprintf("%v_%v_%v_%v", job.CI, job.Organization, job.Repository, job.RunId)
		list[key] = job
	}

	return list, nil
}
