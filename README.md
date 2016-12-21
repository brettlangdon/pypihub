PyPIHub
=======

PyPI server for serving Python packages out of GitHub.

## Installing

```bash
go get github.com/brettlangdon/pypihub
```

## Running

```bash
pypihub -h
usage: pypihub --username USERNAME --access-token ACCESS-TOKEN [--bind BIND] [REPONAMES [REPONAMES ...]]

positional arguments:
  reponames              list of '<username>/<repo>' repos to proxy for (env: PYPIHUB_REPOS)

options:
  --username USERNAME, -u USERNAME
                         Username of GitHub user to login as (env: PYPIHUB_USERNAME)
  --access-token ACCESS-TOKEN, -a ACCESS-TOKEN
                         GitHub personal access token to use for authenticating (env: PYPIHUB_ACCESS_TOKEN)
  --bind BIND, -b BIND   [<address>]:<port> to bind the server to (default: ':8287') (env: PYPIHUB_BIND) [default: :8287]
  --help, -h             display this help and exit
```

### Example

```bash
pypihub -u "<username>" -a "<github-access-token>" "brettlangdon/flask-env" "brettlangdon/flask-defer" [... <owner>/<repo>]
```

```bash
export PYPIHUB_USERNAME="<username>"
export PYPIHUB_ACCESS_TOKEN="<github-access-token>""
export PYPIHUB_REPOS="brettlangdon/flask-env brettlangdon/flask-defer [... <owner>/<repo>]"
pypihub
```

## Docker

```bash
docker run --rm -it -p "8287:8287" -e PYPIHUB_USERNAME="<username>" -e PYPIHUB_ACCESS_TOKEN="<github-acess-token>" -e PYPIHUB_REPOS="<owner>/<repo> ..." brettlangdon/pypihub:latest
```

### Using an env file

```
PYPIHUB_USERNAME=<username>
PYPIHUB_ACCESS_TOKEN=<github-access-token>
PYPIHUB_REPOS=<owner>/<repo> ...
```

```bash
docker run --rm -it -p "8287:8287" --env-file ./.env brettlangdon/pypihub:latest
```

## Endpoints

* `/` - Page containing all links for all projects/assets
  * This endpoint can be used with `--find-links` to make all projects accessible
  * e.g. `pip install --find-links http://localhost:8287/`
* `/<owner>` - Page containing all links for a given GitHub repo owner
  * This endpoint can be used with `--find-links` to make all projects for a given GitHub owner accessible
  * e.g. `pip install --find-links http://localhost:8287/brettlangdon`
* `/<owner>/<repo>` - Page containing all links for a specific GitHub repo
  * This endpoint can be used with `--find-links` to make all releases for a specific GitHub repo accessible
  * e.g. `pip install --find-links http://localhost:8287/brettlangdon/flask-env`
* `/simple` - PyPI simple index page
  * This page lists all of the project names available
  * This endpoint can be used with `--index-url` or `--extra-index-url`
  * e.g. `pip install --extra-index-url http://localhost:8287/simple`
* `/simple/<repo>` - PyPI simple index project links page
  * This page contains the links for the given project name
  * This endpoint can be used with `--find-links`, but is typically used by `pip` when using `--extra-index-url`
  * See `/simple` example above for usage

## Usage with pip

### Simple index
```bash
pip install --index-url http://localhost:8287/simple <project>
pip install --extra-index-url http://localhost:8287/simple <project>
```

### Find links

```bash
pip install --find-links http://localhost:8287/ <project>
pip install --find-links http://localhost:8287/<owner> project
pip install --find-links http://localhost:8287/<owner>/<project> project
```

### requirements.txt

```
--find-links http://localhost:8287/
<project>
```

```bash
pip install -r requirements.txt
```
