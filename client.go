package pypihub

import (
	"io"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
)

type Client struct {
	config Config
	client *github.Client
	repos  []string
}

func NewClient(cfg Config) *Client {
	var t = github.BasicAuthTransport{
		Username: cfg.Username,
		Password: cfg.AccessToken,
	}
	return &Client{
		config: cfg,
		client: github.NewClient(t.Client()),
		repos:  cfg.RepoNames,
	}
}

func (c *Client) splitRepoName(r string) (string, string) {
	var p = strings.SplitN(r, "/", 2)
	switch len(p) {
	case 2:
		return p[0], p[1]
	case 1:
		return c.config.Username, p[0]
	default:
		return "", ""
	}
}

func (c *Client) GetRepoAssets(r string) ([]Asset, error) {
	var owner, repo string
	owner, repo = c.splitRepoName(r)

	var releases []*github.RepositoryRelease
	var err error
	releases, _, err = c.client.Repositories.ListReleases(owner, repo, nil)
	if err != nil {
		return nil, err
	}

	var allAssets = make([]Asset, 0)
	for _, rel := range releases {
		var assets []*github.ReleaseAsset
		assets, _, err = c.client.Repositories.ListReleaseAssets(owner, repo, *rel.ID, nil)
		if err != nil {
			return nil, err
		}
		for _, a := range assets {
			allAssets = append(allAssets, Asset{
				ID:    *a.ID,
				Name:  *a.Name,
				Owner: owner,
				Repo:  repo,
			})
		}
	}
	return allAssets, nil
}

func (c *Client) GetAllAssets() ([]Asset, error) {
	var allAssets = make([]Asset, 0)
	for _, r := range c.config.RepoNames {
		var repoAssets = make([]Asset, 0)
		var err error
		repoAssets, err = c.GetRepoAssets(r)
		if err != nil {
			return nil, err
		}
		allAssets = append(allAssets, repoAssets...)
	}

	return allAssets, nil
}

func (c *Client) DownloadAsset(a Asset) (io.ReadCloser, error) {
	var rc io.ReadCloser
	var redirect string
	var err error
	rc, redirect, err = c.client.Repositories.DownloadReleaseAsset(a.Owner, a.Repo, a.ID)
	if rc == nil {
		var resp *http.Response
		resp, err = http.Get(redirect)
		if err != nil {
			return nil, err
		}
		rc = resp.Body
	}
	return rc, err
}
