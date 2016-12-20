package pypihub

import (
	"os"
	"strings"

	arg "github.com/alexflint/go-arg"
)

type Config struct {
	Username    string   `arg:"-u,--username,env:PYPIHUB_USERNAME,required,help:Username of GitHub user to login as (env: PYPIHUB_USERNAME)"`
	AccessToken string   `arg:"-a,--access-token,env:PYPIHUB_ACCESS_TOKEN,required,help:GitHub personal access token to use for authenticating (env: PYPIHUB_ACCESS_TOKEN)"`
	RepoNames   []string `arg:"positional,help:list of '<username>/<repo>' repos to proxy for (env: PYPIHUB_REPOS)"`
	Bind        string   `arg:"-b,--bind,env:PYPIHUB_BIND,help:[<address>]:<port> to bind the server to (default: ':8287') (env: PYPIHUB_BIND)"`
}

func ParseConfig() Config {
	var config = Config{
		Bind:      ":8287",
		RepoNames: make([]string, 0),
	}

	arg.MustParse(&config)

	if val, ok := os.LookupEnv("PYPIHUB_REPOS"); ok {
		config.RepoNames = append(config.RepoNames, strings.Split(val, " ")...)
	}
	for i := range config.RepoNames {
		config.RepoNames[i] = strings.TrimSpace(config.RepoNames[i])
	}

	config.RepoNames = uniqueSlice(config.RepoNames)

	return config
}
