[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/OpenFogStack/smart-factory-fog-example)

# Smart Factory Fog Example

This repository contains several services that form the smart factory example application shown in the figure:

![](./application.png)

Please note that the related research is still pending publication.

A full list of our [publications](https://www.mcc.tu-berlin.de/menue/forschung/publikationen/parameter/en/) and [prototypes](https://www.mcc.tu-berlin.de/menue/forschung/prototypes/parameter/en/) is available on our group website.

## License

The code in this repository is licensed under the terms of the [MIT](./LICENSE) license.

## Instructions

Use the `Makefile` to build and push services as Docker containers.
Adapt the `REPO` variable to your Docker repo.
Then use `make all -B` or `make [service name] -B` to build and push all services or a specific service.
[If your system supports it](https://www.gnu.org/software/make/manual/html_node/Parallel.html), use `make all -B -j` to build and push all services in parallel.

Alternatively, you can also build binaries directly by going into the directory of a service and running `go build .` (Go in version >= 1.13 is required).
