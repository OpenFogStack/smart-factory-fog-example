# Smart Factory Fog Example

This repository contains several services that form a smart factory example application that can be deployed to the fog.

Please note that the related research is still pending publication.

A full list of our [publications](https://www.mcc.tu-berlin.de/menue/forschung/publikationen/parameter/en/) and [prototypes](https://www.mcc.tu-berlin.de/menue/forschung/prototypes/parameter/en/) is available on our group website.

## License

The code in this repository is licensed under the terms of the [MIT](./LICENSE) license.

## Instructions

Use the `Makefile` to build and push services as Docker containers.
Adapt the `REPO` variable to your Docker repo.
Then use `make all` or `make [service name]` to build and push all services or a specific service.
