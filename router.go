package pypihub

import (
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

	var err error
	r.assets, err = r.client.GetAllAssets()
	if err != nil {
		log.Println(err)
	}
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
	h = mux.NewRouter()
	h.HandleFunc("/", r.handleIndex).Methods("GET")
	h.HandleFunc("/favicon.ico", r.handleFavicon).Methods("GET")
	h.HandleFunc("/simple", r.handleSimple).Methods("GET")
	h.HandleFunc("/{owner}", r.handleOwnerIndex).Methods("GET")
	h.HandleFunc("/{owner}/{repo}", r.handleRepoIndex).Methods("GET")
	h.HandleFunc("/{owner}/{repo}", r.handleRepoIndex).Methods("GET")
	h.HandleFunc("/{owner}/{repo}/{asset}", r.handleFetchAsset).Methods("GET")
	return r.logRequests(h)
}

func (r *Router) Start() error {
	r.refetchAssets()
	http.Handle("/", r.Handler())
	return http.ListenAndServe(r.config.Bind, nil)
}
