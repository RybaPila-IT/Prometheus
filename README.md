# Prometheus example project

In this project we present a simple Golang application which demonstrates the
usage of 3 different Prometheus metric types:

* counter
* gauge
* histogram

The repository contains 3 clients used for testing each of the metrics. Clients
are placed in (surprise!) clients folder.

## How to run it

1. Go to the [Prometheus download page](https://prometheus.io/download/) on the Prometheus website.
2. Scroll down to the "Binary Releases" section and select your operating system from the list.
3. Choose the latest version of Prometheus by clicking on the link next to the version number.
4. Once the download is complete, extract the files from the archive.
5. Provide the path to Prometheus binary for the `start.sh` script.
6. Visit `http:localhot9090` and query Prometheus to see metric values!