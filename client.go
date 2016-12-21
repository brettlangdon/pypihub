package pypihub

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func (c *Client) getRepoTagAssets(owner string, repo string) ([]Asset, error) {
	var tags []*github.RepositoryTag
	var err error
	tags, _, err = c.client.Repositories.ListTags(owner, repo, nil)
	if err != nil {
		return nil, err
	}

	var allAssets = make([]Asset, 0)
	for _, tag := range tags {
		// Remove any `v` prefix, e.g. `v1.0.0` -> `1.0.0`
		var name = strings.Trim(*tag.Name, "v")
		name = fmt.Sprintf("%s-%s.tar.gz", repo, name)
		allAssets = append(allAssets, Asset{
			Name:   name,
			Owner:  owner,
			Repo:   repo,
			Ref:    *tag.Name,
			Format: "tarball",
		})
	}

	return allAssets, nil
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

	if len(releases) == 0 {
		return c.getRepoTagAssets(owner, repo)
	}

	var allAssets = make([]Asset, 0)
	for _, rel := range releases {
		var assets []*github.ReleaseAsset
		assets, _, err = c.client.Repositories.ListReleaseAssets(owner, repo, *rel.ID, nil)
		if err != nil {
			return nil, err
		}

		var hasTar = false
		for _, a := range assets {
			if strings.HasSuffix(*a.Name, ".tar.gz") {
				hasTar = true
			}
			allAssets = append(allAssets, Asset{
				ID:    *a.ID,
				Name:  *a.Name,
				Owner: owner,
				Repo:  repo,
			})
		}

		if hasTar == false {
			// Remove any `v` prefix, e.g. `v1.0.0` -> `1.0.0`
			var name = strings.Trim(*rel.Name, "v")
			name = fmt.Sprintf("%s-%s.tar.gz", repo, name)
			allAssets = append(allAssets, Asset{
				Name:   name,
				Owner:  owner,
				Repo:   repo,
				Ref:    *rel.TagName,
				Format: "tarball",
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

func (c *Client) DownloadArchive(a Asset) (io.ReadCloser, error) {
	var f = github.Tarball
	if a.Format == "zipball" {
		f = github.Zipball
	}

	var u *url.URL
	var err error
	u, _, err = c.client.Repositories.GetArchiveLink(a.Owner, a.Repo, f, nil)
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	resp, err = http.Get(u.String())
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
