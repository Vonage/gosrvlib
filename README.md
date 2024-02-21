<!-- Space: APIOSS -->
<!-- Parent: Projects -->
<!-- Title: gosrvlib -->

# gosrvlib

*Go Service Library*

This Open Source project contains a collection of high-quality [GO](https://go.dev/) (golang) packages.

Each package follows common conventions and they can be individually imported in any project.

This package collection forms the base structure for fully-fledged production-ready web-services.

A new Web service can be generated by using the command `make project CONFIG=project.cfg`.
The new generated project name, description, etc..., can be set in the file specified via the CONFIG parameter.

The packages documentation is available at: [https://pkg.go.dev/github.com/Vonage/gosrvlib/](https://pkg.go.dev/github.com/Vonage/gosrvlib)

[![Go Reference](https://pkg.go.dev/badge/github.com/Vonage/gosrvlib.svg)](https://pkg.go.dev/github.com/Vonage/gosrvlib)   
[![check](https://github.com/Vonage/gosrvlib/actions/workflows/check.yaml/badge.svg)](https://github.com/Vonage/gosrvlib/actions/workflows/check.yaml)
[![Coverage Status](https://coveralls.io/repos/github/Vonage/gosrvlib/badge.svg?branch=main)](https://coveralls.io/github/Vonage/gosrvlib?branch=main)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=Vonage_gosrvlib&metric=coverage)](https://sonarcloud.io/summary/new_code?id=Vonage_gosrvlib)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=Vonage_gosrvlib&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=Vonage_gosrvlib)  
[![Go Report Card](https://goreportcard.com/badge/github.com/Vonage/gosrvlib)](https://goreportcard.com/report/github.com/Vonage/gosrvlib)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Vonage_gosrvlib&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=Vonage_gosrvlib)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=Vonage_gosrvlib&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=Vonage_gosrvlib)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=Vonage_gosrvlib&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=Vonage_gosrvlib)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=Vonage_gosrvlib&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=Vonage_gosrvlib)  
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=Vonage_gosrvlib&metric=bugs)](https://sonarcloud.io/summary/new_code?id=Vonage_gosrvlib)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=Vonage_gosrvlib&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=Vonage_gosrvlib)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=Vonage_gosrvlib&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=Vonage_gosrvlib)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=Vonage_gosrvlib&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=Vonage_gosrvlib)


* **category**    Library
* **license**     [MIT](https://github.com/Vonage/gosrvlib/blob/main/LICENSE)
* **link**        https://github.com/Vonage/gosrvlib

-----------------------------------------------------------------

## TOC

* [Quick Start](#quickstart)
* [Running all tests](#runtest)
* [Examples](#examples)

-----------------------------------------------------------------

<a name="quickstart"></a>
## Developers' Quick Start

To quickly get started with this project, follow these steps:

1. Ensure you have installed the latest Go version and Python3 for some extra tests.
1. Clone the repository: `git clone https://github.com/Vonage/gosrvlib.git`.
2. Change into the project directory: `cd gosrvlib`.
3. Install the required dependencies and test everything: `DEVMODE=LOCAL make x`.

Now you are ready to start developing with gosrvlib!


This project includes a *Makefile* that allows you to test and build the project in a Linux-compatible system with simple commands.  
All the artifacts and reports produced using this *Makefile* are stored in the *target* folder.  

Alternatively, everything can be built inside a [Docker](https://www.docker.com) container using the command `make dbuild` that uses the environment defined at `resources/docker/Dockerfile.dev`.

To see all available options:
```bash
make help
```

-----------------------------------------------------------------

<a name="runtest"></a>
## Running all tests

Before committing the code, please format it and check if it passes all tests using
```bash
DEVMODE=LOCAL make x
```

-----------------------------------------------------------------

<a name="examples"></a>
## Examples

Please check the `examples/service` folder for an example of a service based on this library.

The following command generates a new project from the example using the data set in the `project.cfg` file:

```bash
make project CONFIG=project.cfg
```

-----------------------------------------------------------------
