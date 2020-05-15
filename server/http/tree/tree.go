package tree

import (
	"encoding/json"
	"log"
	"net/http"

	"ritchie-server/server"
	"ritchie-server/server/tm"
)

type Handler struct {
	Config        server.Config
	Authorization server.Constraints
}

const (
	repoNameHeader      = "x-repo-name"
	authorizationHeader = "Authorization"
)

func NewConfigHandler(config server.Config, auth server.Constraints) server.DefaultHandler {
	return Handler{Config: config, Authorization: auth}
}

func (lh Handler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var method = r.Method
		if http.MethodGet == method {
			lh.processGet(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}

func (lh Handler) processGet(w http.ResponseWriter, r *http.Request) {
	org := r.Header.Get(server.OrganizationHeader)
	repos, err := lh.Config.ReadRepositoryConfig(org)
	if err != nil {
		log.Printf("Error while processing %v's repository configuration: %v", org, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(repos) == 0 {
		log.Println("No repository config found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	repoName := r.Header.Get(repoNameHeader)
	repo, err := tm.FindRepo(repos, repoName)
	if err != nil {
		log.Printf("no repo for org %s, with name %s, error: %v", org, repoName, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bt := r.Header.Get(authorizationHeader)
	finalTree, err := tm.TreeRemoteAllow(lh.Authorization, bt, org, r.URL.Path, repo)
	if err != nil {
		log.Printf("Error load final tree. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(finalTree); err != nil {
		log.Printf("Error encode finalTree: %v", finalTree)
	}
}
