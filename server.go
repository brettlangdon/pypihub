package pypihub

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	client *Client
	config Config
	assets []Asset
	timer  *time.Timer
}

func NewServer(cfg Config) *Server {
	return &Server{
		client: NewClient(cfg),
		config: cfg,
		assets: make([]Asset, 0),
	}
}

func (s *Server) findAsset(owner string, repo string, id int) *Asset {
	for _, a := range s.assets {
		if a.Owner == owner && a.Repo == repo && a.ID == id {
			return &a
		}
	}
	return nil
}

func (s *Server) refetchAssets() {
	var err error
	s.assets, err = s.client.GetAllAssets()
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) startTimer() {
	if s.timer != nil {
		s.timer.Stop()
	}
	s.timer = time.AfterFunc(time.Duration(10*time.Minute), func() {
		go s.refetchAssets()
	})
}

func (s *Server) listAssets(w http.ResponseWriter, r *http.Request) {
	var repo = strings.Trim(r.URL.Path, "/")

	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><title>Links for %s</title><meta name=\"api-version\" content=\"2\" /><body>", repo)
	fmt.Fprintf(w, "<h1>Links for %s</h1>", repo)
	for _, a := range s.assets {
		if a.Repo == repo {
			fmt.Fprintf(w, "<a href=\"%s\">%s</a> ", a.URL(), a.Name)
		}
	}
	fmt.Fprintf(w, "</body></html>")
}

func (s *Server) listRepoAssets(w http.ResponseWriter, r *http.Request) {
	var parts = strings.SplitN(strings.Trim(r.URL.Path, "/"), "/", 2)
	var owner = parts[0]
	var repo = parts[1]

	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><title>Links for %s</title><meta name=\"api-version\" content=\"2\" /><body>", repo)
	fmt.Fprintf(w, "<h1>Links for %s</h1>", repo)
	for _, a := range s.assets {
		if a.Owner == owner && a.Repo == repo {
			fmt.Fprintf(w, "<a href=\"%s\">%s</a> ", a.URL(), a.Name)
		}
	}
	fmt.Fprintf(w, "</body></html>")
}

func (s *Server) listAllAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><title>All asset links</title><meta name=\"api-version\" content=\"2\" /><body>")
	fmt.Fprintf(w, "<h1>All asset links</h1>")
	for _, a := range s.assets {
		fmt.Fprintf(w, "<a href=\"%s\">%s</a> ", a.URL(), a.Name)
	}
	fmt.Fprintf(w, "</body></html>")
}

func (s *Server) listRepos(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><title>Simple index</title><meta name=\"api-version\" content=\"2\" /><body>")
	for _, r := range s.config.RepoNames {
		var parts = strings.SplitN(r, "/", 2)
		fmt.Fprintf(w, "<a href=\"/%s\">%s</a> ", parts[1], parts[1])
	}
	fmt.Fprintf(w, "</body></html>")
}

func (s *Server) fetchAsset(w http.ResponseWriter, r *http.Request) {
	var url = strings.Trim(r.URL.Path, "/")
	var parts = strings.SplitN(url, "/", 4)

	if len(parts) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var asset *Asset
	var id int64
	var err error
	id, err = strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	asset = s.findAsset(parts[0], parts[1], int(id))
	if asset == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var rc io.ReadCloser
	rc, err = s.client.DownloadAsset(*asset)
	if err != nil || rc == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer rc.Close()
	io.Copy(w, rc)
}

func (s *Server) serveFavicon(w http.ResponseWriter, r *http.Request) {
	var decoded []byte
	var err error
	decoded, err = base64.StdEncoding.DecodeString("AAABAAEAEBAAAAEACABoBQAAFgAAACgAAAAQAAAAIAAAAAEACAAAAAAAAAEAAAAAAAAAAAAAAAEAAAAAAAAAAAAAPs7/AJ9vNwBO3f8APMz/AEzb/wCdbjYAStn/AJZqNgBI1/8ARtX/AETT/wBA0f8AqnY3AD7P/wCjcjcApHI3ADzN/wBQ3v8AOsv/AE7c/wA4yf8Amm02ALF6OACYazYARtb/AETU/wCveTcAQtL/AKh1NwCtdzcAQND/AKFxNwCmczcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAcDA4RExUAAAAAAAAAAAAZGhwfDhEAFQAAAAAAAAAACRkaHB8OABMAAAAAAAAAAAcJGRocHw4RAAAAAAAgAgYFBwkZGhwfARETFQAhDyACAwUHCQoaHB8BERMVHSEPIBIUBQcJChocHwEREw0dIQ8AEhQFBwkKGhwfAQQeDR0hDyACBhYYCAALHB8BGx4NHSEPIAIGFhgICgscHxcbHg0dIRAgAgYWGAkKCxwAFxseDR0hECACBhYHCQoAAAAAAB4NHSEQIAIGAAAAAAAAAAAbAA0dIRAgAgAAAAAAAAAAFwAeDR0hECAAAAAAAAAAAAAXGx4NHSEAAAAAAPgfAADwLwAA8C8AAPAPAACAAQAAAAAAAAAAAAAIAAAAABAAAAAAAAAAAAAAgAEAAPAPAAD0DwAA9A8AAPgfAAA=")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s", decoded)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var parts = strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	parts = removeEmpty(parts)

	fmt.Println(r.URL.Path)

	switch len(parts) {
	case 0:
		s.listAllAssets(w, r)
	case 1:
		if parts[0] == "favicon.ico" {
			s.serveFavicon(w, r)
		} else {
			s.listAssets(w, r)
		}
	case 2:
		s.listRepoAssets(w, r)
	default:
		s.fetchAsset(w, r)
	}
}

func (s *Server) ListenAndServe() error {
	s.refetchAssets()
	s.startTimer()

	http.Handle("/", s)
	fmt.Println("Server listening at", s.config.Bind)
	return http.ListenAndServe(s.config.Bind, nil)
}
