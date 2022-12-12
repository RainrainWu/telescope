# Telescope

Telescope is a dependencies scanner that helps developers sort out outdated dependencies with desired scope easily and review the necessity to arrange an upgrade plan.

[![Go Report Card](https://goreportcard.com/badge/github.com/rainrainwu/telescope)](https://goreportcard.com/report/github.com/rainrainwu/telescope)
[![codecov](https://codecov.io/gh/RainrainWu/telescope/branch/master/graph/badge.svg?token=4HWXOZEIsG)](https://codecov.io/gh/RainrainWu/telescope)

#### Currently Supported Dependencies Lock Files
- `go.mod`
- `poetry.lock`
- `Pipfile.lock`

## Usage
```
$ docker run --rm docker.io/r41nwu/telescope:latest

Usage: telescope [-f file_path] [-s outdated_scope] [-i ignored_dependency] [-c critical_dependency] [--skip-unknown] [--strict-semver]
  -c value
        highlight critical dependencies with regular expression
  -f string
        dependencies file path (default "go.mod")
  -i value
        ignore specific dependencies with regular expression
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
    telescope -f "go.mod" -s "minor" -c "major:.*"
```

- With `poetry.lock` file
```
docker run \
    --rm \
    -v "$YOUR_POETRY_LOCK_FILE:/poetry.lock" \
    docker.io/r41nwu/telescope:latest \
    telescope -f "poetry.lock" -s "minor" -c "major:.*"
```

- With `Pipfile.lock` file
```
docker run \
    --rm \
    -v "$YOUR_PIPFILE_LOCK_FILE:/Pipfile.lock" \
    docker.io/r41nwu/telescope:latest \
    telescope -f "Pipfile.lock" -s "minor" -c "major:.*"
```

### Advanced Flags Usage

#### `-s` Desired Scope
Specify the scope you want while reporting outdated dependencies.
```
// report dependencies with outdated version on major scope only
telescope -s "major"

// report dependencies with outdated version on any of major, minor, or patch scope
telescope -s "patch"
```

#### `-c` Critical Dependencies
Return a non-zero exit code if any of the matched dependencies is outdated (dependencies will start with a `*` prefix).
```
// raise error if there is any dependency outdated on major version
telescope -c "major:.*"

// raise error if any dependency from `golang.org` or `k8s.io` has a outdated scope greater than minor version
telescope -c "minor:^golang.org/.*$" -c "minor:^k8s.io/.*$"
```

#### `-i` Ignored Dependencies
Ignored dependencies will not be taken into account during reporting.
```
// ignore all pytest-related packages
telescope -i "^pytest.*$"
```

#### `--skip-unknown` Skip Dependencies with Unknown Version
Skip dependency if its current version can not be parsed or unable to obtained the latest version from package index url.
```
// skip dependencies with unknown version
telescope --skip-unknown
```

#### `--strict-semver` Strict Semantic Version
By default telescope will tend to truncated useless information (e.g. alpha/beta release tag) and parse as many version expressions as possible, but you are still able to force apply strict semver format and the malformed expression will be treated as unknown one.
```
// force apply strict semantic version format
telescope --strict-semver
```

### Example Output
```
[ 6 MAJOR Version Outdated ]========================================

  github.com/Azure/azure-sdk-for-go                  65.0.0+incompatible  67.1.0+incompatible 
  github.com/docker/docker                           20.10.21+incompatible 23.0.0-beta.1+incompatible
  github.com/evanphx/json-patch                      4.12.0+incompatible  5.6.0+incompatible  
* github.com/hashicorp/go-hclog                      0.14.1               1.4.0               
* github.com/hashicorp/golang-lru                    0.5.4                1.0.1               
  k8s.io/client-go                                   0.25.3               11.0.0+incompatible 


[ 47 MINOR Version Outdated ]========================================

  cloud.google.com/go/compute                        1.12.1               1.14.0              
  github.com/Microsoft/go-winio                      0.5.1                0.6.0               
  github.com/PuerkitoBio/purell                      1.1.1                1.2.0               
  github.com/armon/go-metrics                        0.3.10               0.4.1               
  github.com/cenkalti/backoff/v4                     4.1.3                4.2.0               
  github.com/cespare/xxhash/v2                       2.1.2                2.2.0               
  github.com/coreos/go-systemd/v22                   22.4.0               22.5.0              
  github.com/digitalocean/godo                       1.89.0               1.91.1              
  github.com/docker/distribution                     2.7.1+incompatible   2.8.1+incompatible  
  github.com/emicklei/go-restful/v3                  3.8.0                3.10.1              
  github.com/envoyproxy/protoc-gen-validate          0.8.0                0.10.0-SNAPSHOT.0   
  github.com/go-kit/kit                              0.10.0               0.12.0              
  github.com/go-openapi/jsonreference                0.19.6               0.20.0              
  github.com/go-openapi/swag                         0.21.1               0.22.3              
  github.com/go-openapi/validate                     0.21.0               0.22.0
```

## Development
- Lint code
```
$ make lint
```

- Run Unit Tests
```
$ make test
```

- Build docker image
```
$ make build
```
