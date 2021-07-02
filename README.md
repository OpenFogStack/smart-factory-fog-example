# Smart Factory Fog Example

This repository contains several services that form a smart factory example application that can be deployed to the fog.

## Research

If you use this software in a publication, please cite it as:

### Text

T. Pfandzelter, J. Hasenburg, and D. Bermbach, **From Zero to Fog: Efficient Engineering of Fog-Based Internet of Things Applications**, Software: Practice and Experience, vol. 51, no. 8, pp. 1798–1821, Aug. 2021.

### BibTeX

```bibtex
@article{pfandzelter-zero2fog-wiley,
    author = "Pfandzelter, Tobias and Hasenburg, Jonathan and Bermbach, David",
    title = "From Zero to Fog: Efficient Engineering of Fog-Based Internet of Things Applications",
    journal = "Software: Practice and Experience",
    year = 2021,
    volume = 51,
    number = 8,
    pages = "1798--1821",
    publisher = "Wiley"
}
```

For a full list of publications, please see [our website](https://www.mcc.tu-berlin.de/menue/forschung/publikationen/parameter/en/).

## License

The code in this repository is licensed under the terms of the [MIT](./LICENSE) license.

## Instructions

Use the `Makefile` to build and push services as Docker containers.
Adapt the `REPO` variable to your Docker repo.
Then use `make all -B` or `make [service name] -B` to build and push all services or a specific service.
[If your system supports it](https://www.gnu.org/software/make/manual/html_node/Parallel.html), use `make all -B -j` to build and push all services in parallel.

Alternatively, you can also build binaries directly by going into the directory of a service and running `go build .` (Go in version >= 1.13 is required).
