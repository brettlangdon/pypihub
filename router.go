package pypihub

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Router struct {
	config Config
	client *Client
	assets []Asset
	timer  *time.Timer
}

func NewRouter(config Config) *Router {
	return &Router{
		config: config,
		client: NewClient(config),
		assets: make([]Asset, 0),
	}
}

func (r *Router) refetchAssets() {
	go r.startAssetsTimer()

	log.Printf("refetching assets for %d repos", len(r.config.RepoNames))
	var err error
	r.assets, err = r.client.GetAllAssets()
	if err != nil {
		log.Println(err)
	}
	log.Printf("found %d assets for %d repos", len(r.assets), len(r.config.RepoNames))
}

func (r *Router) startAssetsTimer() {
	if r.timer != nil {
		r.timer.Stop()
	}
	r.timer = time.AfterFunc(5*time.Minute, r.refetchAssets)
}

func (r *Router) handleSimple(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<html><title>Simple index</title><meta name=\"api-version\" value=\"2\" /><body>")
	var projects = make(map[string]bool)
	for _, a := range r.assets {
		projects[strings.ToLower(a.Repo)] = true
	}

	for project := range projects {
		fmt.Fprintf(w, "<a href=\"/simple/%s\">%s</a> ", project, project)
	}
	fmt.Fprintf(w, "</body></html>")
}

func (r *Router) handleSimpleProject(w http.ResponseWriter, req *http.Request) {
	var vars map[string]string
	vars = mux.Vars(req)
	var repo = strings.ToLower(vars["repo"])

	fmt.Fprintf(w, "<html><title>Links for %s</title><meta name=\"api-version\" value=\"2\" /><body>", repo)
	fmt.Fprintf(w, "<h1>Links for all %s</h1>", repo)
	for _, a := range r.assets {
		if strings.ToLower(a.Repo) == repo {
			fmt.Fprintf(w, "<a href=\"%s\">%s</a> ", a.URL(), a.Name)
		}
	}
	fmt.Fprintf(w, "</body></html>")
}

func (r *Router) handleIndex(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<html><title>Links for all projects</title><body>")
	fmt.Fprintf(w, "<h1>Links for all projects</h1>")
	for _, a := range r.assets {
		fmt.Fprintf(w, "<a href=\"%s\">%s</a> ", a.URL(), a.Name)
	}
	fmt.Fprintf(w, "</body></html>")
}

func (r *Router) handleFavicon(w http.ResponseWriter, req *http.Request) {
	var decoded []byte
	var err error
	decoded, err = base64.StdEncoding.DecodeString("AAABAAEAEBAAAAEACABoBQAAFgAAACgAAAAQAAAAIAAAAAEACAAAAAAAAAEAAAAAAAAAAAAAAAEAAAAAAAAAAAAAPs7/AJ9vNwBO3f8APMz/AEzb/wCdbjYAStn/AJZqNgBI1/8ARtX/AETT/wBA0f8AqnY3AD7P/wCjcjcApHI3ADzN/wBQ3v8AOsv/AE7c/wA4yf8Amm02ALF6OACYazYARtb/AETU/wCveTcAQtL/AKh1NwCtdzcAQND/AKFxNwCmczcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAcDA4RExUAAAAAAAAAAAAZGhwfDhEAFQAAAAAAAAAACRkaHB8OABMAAAAAAAAAAAcJGRocHw4RAAAAAAAgAgYFBwkZGhwfARETFQAhDyACAwUHCQoaHB8BERMVHSEPIBIUBQcJChocHwEREw0dIQ8AEhQFBwkKGhwfAQQeDR0hDyACBhYYCAALHB8BGx4NHSEPIAIGFhgICgscHxcbHg0dIRAgAgYWGAkKCxwAFxseDR0hECACBhYHCQoAAAAAAB4NHSEQIAIGAAAAAAAAAAAbAA0dIRAgAgAAAAAAAAAAFwAeDR0hECAAAAAAAAAAAAAXGx4NHSEAAAAAAPgfAADwLwAA8C8AAPAPAACAAQAAAAAAAAAAAAAIAAAAABAAAAAAAAAAAAAAgAEAAPAPAAD0DwAA9A8AAPgfAAA=")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s", decoded)
}

func (r *Router) handleOwnerIndex(w http.ResponseWriter, req *http.Request) {
	var vars map[string]string
	vars = mux.Vars(req)
	var owner = strings.ToLower(vars["owner"])

	fmt.Fprintf(w, "<html><title>Packages for %s</title><body>", owner)
	fmt.Fprintf(w, "<h1>Links for %s projects</h1>", owner)
	for _, a := range r.assets {
		if strings.ToLower(a.Owner) == owner {
			fmt.Fprintf(w, "<a href=\"%s\">%s</a> ", a.URL(), a.Name)
		}
	}
	fmt.Fprintf(w, "</body></html>")
}

func (r *Router) handleRepoIndex(w http.ResponseWriter, req *http.Request) {
	var vars map[string]string
	vars = mux.Vars(req)
	var owner = strings.ToLower(vars["owner"])
	var repo = strings.ToLower(vars["repo"])

	fmt.Fprintf(w, "<html><title>Packages for %s/%s</title><body>", owner, repo)
	fmt.Fprintf(w, "<h1>Links for all %s/%s</h1>", owner, repo)
	for _, a := range r.assets {
		if strings.ToLower(a.Owner) == owner && strings.ToLower(a.Repo) == repo {
			fmt.Fprintf(w, "<a href=\"%s\">%s</a> ", a.URL(), a.Name)
		}
	}
	fmt.Fprintf(w, "</body></html>")
}

func (r *Router) handleFetchAsset(w http.ResponseWriter, req *http.Request) {
	var vars map[string]string
	vars = mux.Vars(req)
	var owner = strings.ToLower(vars["owner"])
	var repo = strings.ToLower(vars["repo"])
	var asset = vars["asset"]

	for _, a := range r.assets {
		if strings.ToLower(a.Owner) == owner && strings.ToLower(a.Repo) == repo && a.Name == asset {
			var rc io.ReadCloser
			var err error
			rc, err = a.Download(r.client)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			io.Copy(w, rc)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (r *Router) logRequests(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.URL.Path)
		h.ServeHTTP(w, req)
	})
}

func (r *Router) Handler() http.Handler {
	var h *mux.Router
	h = mux.NewRouter().StrictSlash(false)

	// Static favicon
	h.HandleFunc("/favicon.ico", r.handleFavicon).Methods("GET")

	// All links, useful for non-owner/repo specific --find-links
	h.HandleFunc("/", r.handleIndex).Methods("GET")

	// Simple index
	h.HandleFunc("/simple", r.handleSimple).Methods("GET")
	h.HandleFunc("/simple/", r.handleSimple).Methods("GET")
	h.HandleFunc("/simple/{repo}", r.handleSimpleProject).Methods("GET")
	h.HandleFunc("/simple/{repo}/", r.handleSimpleProject).Methods("GET")

	// Owner/repo specific find-links
	h.HandleFunc("/{owner}", r.handleOwnerIndex).Methods("GET")
	h.HandleFunc("/{owner}/", r.handleOwnerIndex).Methods("GET")
	h.HandleFunc("/{owner}/{repo}", r.handleRepoIndex).Methods("GET")
	h.HandleFunc("/{owner}/{repo}/", r.handleRepoIndex).Methods("GET")

	// Download asset
	h.HandleFunc("/{owner}/{repo}/{asset}", r.handleFetchAsset).Methods("GET")
	return r.logRequests(h)
}

func (r *Router) Start() error {
	r.refetchAssets()
	http.Handle("/", r.Handler())
	return http.ListenAndServe(r.config.Bind, nil)
}
