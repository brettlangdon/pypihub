package pypihub

import "fmt"

type Asset struct {
	ID    int
	Name  string
	Owner string
	Repo  string
}

func (a Asset) String() string {
	return a.Name
}

func (a Asset) URL() string {
	return fmt.Sprintf("/%s/%s/%d/%s", a.Owner, a.Repo, a.ID, a.Name)
}
