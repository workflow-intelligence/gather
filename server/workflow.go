package server

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CreateRequest struct {
	Organization string `json:"Organization"`
	Repository   string `json:"repository"`
	RunID        int64  `json:"run_id"`
}

type CreateResponse struct {
	Error  *ErrorMessage `json:"error"`
	Status string        `json:"status"`
	Index  string        `json:"index"`
}

func WorkflowCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var c CreateRequest
	err := decoder.Decode(&c)
	if err != nil {
		errorOut(w, 500, "ERR00010020", "Login error", err)
		return
	}

	log.Debug().Str("Organization", c.Organization).Str("repository", c.Repository).Int64("runid", c.RunID).Str("id", "DBG0001010").Msg("Create workflow")
	index, err := backend.CreateWorkflowIndex(c.Organization, c.Repository, c.RunID)
	if err != nil {
		errorOut(w, 500, "ERR00010021", "Could not create index", err)
		return
	}
	err = backend.AddJob(c.Organization, c.Repository, c.RunID)

	CR := CreateResponse{Status: "OK", Index: index}

	log.Debug().Str("status", CR.Status).Str("id", "DBG00010011").Msg("Workflow created")
	json, err := json.Marshal(CR)
	if err != nil {
		errorOut(w, 500, "ERR00010022", "Could marshal login info", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(json[:]))
}
