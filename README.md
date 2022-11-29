# Telescope

Telescope is a dependencies scanner that helps developers sort out outdated dependencies with desired scope easily and review the necessity to arrange an upgrade plan.

#### Currently Supported Dependencies Lock Files
- `go.mod`
- `poetry.lock`

## Usage
```
$ docker run --rm docker.io/r41nwu/telescope:latest

Usage: telescope [-f file_path] [-s outdated_scope] [-i ignored_dependency] [--skip-unknown] [ [--strict-semver]]
  -f string
        dependencies file path (default "go.mod")
  -i value
        ignore specific dependency
  -s string
        desired outdated scope (default "major")
  -skip-unknown
        skip dependencies with unknown versions
  -strict-semver
        parse dependencies file with strict SemVer format
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
    telescope -f "poetry.lock" -s "minor"
```

### Example Output
```
[MAJOR Version Outdated]====================
charset-normalizer                       2.1.1                3.0.1
cryptography                             37.0.4               38.0.4
executing                                0.10.0               1.2.0
moto                                     3.1.18               4.0.10
protobuf                                 3.20.3               4.21.9
pymongo                                  3.13.0               4.3.3
pytest-cov                               3.0.0                4.0.0
pytest-xdist                             2.5.0                3.0.2
rfc3986                                  1.5.0                2.0.0


[MINOR Version Outdated]====================
devtools                                 0.8.0                0.10.0
fastapi                                  0.78.0               0.88.0
googleapis-common-protos                 1.56.2               1.57.0
grpcio                                   1.50.0               1.51.0
h11                                      0.12.0               0.14.0
httpcore                                 0.15.0               0.16.2
hypothesis                               6.57.1               6.58.1
prometheus-client                        0.14.1               0.15.0
pytest-asyncio                           0.18.3               0.20.2
setuptools                               65.5.1               65.6.3
starlette                                0.19.1               0.22.0
uvicorn                                  0.17.6               0.20.0
virtualenv                               20.16.7              20.17.0
werkzeug                                 2.1.2                2.2.2
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
