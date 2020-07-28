# srvxmplname

*Example usage for gosrvlib*

[![Build Status](https://travis-ci.com/nexmoinc/srvxmplname.svg?token=YHpDM41jM29w1XFFg2HR&branch=main)](https://travis-ci.com/nexmoinc/srvxmplname?token=YHpDM41jM29w1XFFg2HR&branch=main)

* **category**    Application
* **copyright**   2020 Vonage
* **license**     see [LICENSE](LICENSE)
* **link**        https://github.com/nexmoinc/srvxmplname

-----------------------------------------------------------------

## TOC

* [Description](#description)
* [Requirements](#requirements)
* [Quick Start](#quickstart)
* [Running all tests](#runtest)
* [Usage](#usage)
* [Configuration](#configuration)
* [Examples](#examples)
* [Logs](#logs)
* [Metrics](#metrics)
* [Profiling](#profiling)
* [OpenApi](#openapi)
* [Docker](#docker)
* [Links](#links)

-----------------------------------------------------------------

<a name="description"></a>
## Description

Project description.

-----------------------------------------------------------------

<a name="requirements"></a>
## Requirements

An additional Python program is used to check the validity of the JSON configuration files against a JSON schema:

```
sudo pip install jsonschema
```

-----------------------------------------------------------------

<a name="quickstart"></a>
## Quick Start

This project includes a Makefile that allows you to test and build the project in a Linux-compatible system with simple commands.  
All the artifacts and reports produced using this Makefile are stored in the *target* folder.  

All the packages listed in the *resources/DockerDev/Dockerfile* file are required in order to build and test all the library options in the current environment. Alternatively, everything can be built inside a [Docker](https://www.docker.com) container using the command "make dbuild".

To see all available options:
```
make help
```

To download the dependencies
```
make deps mod
```

To format the code (please use this command before submitting any pull request):
```
make format
```

To execute all the default test builds and generate reports in the current environment:
```
make qa
```

To build the executable file:
```
make build
```


-----------------------------------------------------------------

<a name="runtest"></a>
## Running all tests

Before committing the code, please check if it passes all tests using
```bash
make qa
```

Other make options are available install this library globally and build RPM and DEB packages.
Please check all the available options using `make help`.

-----------------------------------------------------------------

<a name="usage"></a>
## Usage

```bash
srvxmplname [flags]

Flags:

-c, --configDir  string  Configuration directory to be added on top of the search list
-f, --logFormat  string  Logging format: CONSOLE, JSON
-o, --loglevel   string  Log level: EMERGENCY, ALERT, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG
```

----------------------------------------------------------------

<a name="configuration"></a>
## Configuration

See [CONFIG.md](CONFIG.md).

-----------------------------------------------------------------

<a name="examples"></a>
## Examples

Once the application has being compiled with `make build`, it can be quickly tested:

```bash
target/usr/bin/srvxmplname -c resources/test/etc/srvxmplname
```

<a name="logs"></a>
## Logs

This program logs the log messages in JSON format:

```
{
    "URI":"/",
    "code":200,
    "datetime":"2020-10-06T14:56:48Z",
    "hostname":"myserver",
    "level":"info",
    "msg":"request",
    "program":"srvxmplname",
    "release":"1",
    "timestamp":1475765808084372773,
    "type":"GET",
    "version":"1.0.0"
}
```

Logs are sent to stderr by default.

The log level can be set either in the configuration or as command argument (`logLevel`).

-----------------------------------------------------------------

<a name="metrics"></a>
## Metrics

This service provides [Prometheus](https://prometheus.io/) metrics at the `/metrics` endpoint.

-----------------------------------------------------------------

<a name="profiling"></a>
## Profiling

This service provides [PPROF](https://github.com/google/pprof) profiling data at the `/pprof` endpoint.

The pprof data can be analyzed and displayed using the pprof tool:

```
go get github.com/google/pprof
```

Example:

```
pprof -seconds 10 -http=localhost:8182 http://INSTANCE_URL:PORT/pprof/profile
```

-----------------------------------------------------------------

<a name="openapi"></a>
## OpenAPI

The srvxmplname API is specified via the [OpenAPI 3](https://www.openapis.org/) file: `openapi.yaml`.

The openapi file can be edited using the Swagger Editor:

```
docker pull swaggerapi/swagger-editor
docker run -p 8056:8080 swaggerapi/swagger-editor
```

and pointing the Web browser to http://localhost:8056

-----------------------------------------------------------------

<a name="docker"></a>
## Docker

To build a Docker scratch container for the srvxmplname executable binary execute the following command:
```
make docker
```

To push the Docker container in our ECR repo execute:
```
make dockerpush
```
Note that this command will require to set the follwoing environmental variables or having an AWS profile installed:

* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`
* `AWS_DEFAULT_REGION`


### Useful Docker commands

To manually create the container you can execute:
```
docker build --tag="vonage/srvxmplnamedev" .
```

To log into the newly created container:
```
docker run -t -i vonage/srvxmplnamedev /bin/bash
```

To get the container ID:
```
CONTAINER_ID=`docker ps -a | grep vonage/srvxmplnamedev | cut -c1-12`
```

To delete the newly created docker container:
```
docker rm -f $CONTAINER_ID
```

To delete the docker image:
```
docker rmi -f vonage/srvxmplnamedev
```

To delete all containers
```
docker rm $(docker ps -a -q)
```

To delete all images
```
docker rmi $(docker images -q)
```

-----------------------------------------------------------------

