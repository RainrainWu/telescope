# Telescope

Telescope is a dependencies scanner that helps developers sort out outdated dependencies with desired scope easily and review the necessity to arrange an upgrade plan.

## Usage
```
$ docker run --rm docker.io/r41nwu/telescope:latest

Usage: telescope [-f file_path] [-s outdated_scope]
  -f string
        dependencies file path (default "go.mod")
  -s string
        desired outdated scope (default "major")
```

### Pull the docker image
```
$ docker pull r41nwu/telescope:latest
```

### Execute inside docker container

- With `go.mod` file
```
docker run \
    --rm \
    -v "$YOUR_GO_MOD_FILE:/go.mod" \
    docker.io/r41nwu/telescope:latest \
    telescope -f "go.mod" -s "major"
```

- With `poetry.lock` file
```
docker run \
    --rm \
    -v "$YOUR_POETRY_LOCK_FILE:/poetry.lock" \
    docker.io/r41nwu/telescope:latest \
    telescope -f "poetry.lock" -s "major"
```

## Development
- Lint code
```
$ make lint
```

- Build docker image
```
$ make build
```
