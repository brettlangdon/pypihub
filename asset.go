package pypihub

import (
	"fmt"
	"io"
)

type Asset struct {
	ID     int
	Name   string
	Owner  string
	Repo   string
	Ref    string
	Format string
}

func (a Asset) String() string {
	return a.Name
}

func (a Asset) URL() string {
	return fmt.Sprintf("/%s/%s/%s", a.Owner, a.Repo, a.Name)
}

func (a Asset) Download(c *Client) (io.ReadCloser, error) {
	if a.Ref != "" && a.Format != "" {
		return c.DownloadArchive(a)
	}

	return c.DownloadAsset(a)
}
