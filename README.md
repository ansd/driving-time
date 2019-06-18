# driving-time

`driving-time` is a CLI (and server) displaying driving time based on live traffic data.

## Requirement
A Google Cloud Platform API key (see [here](https://developers.google.com/maps/documentation/distance-matrix/get-api-key) how to get one)
since the CLI queries the Google Maps Distance Matrix API.

## Installation
Options:
1. Download the binary in https://github.com/ansd/driving-time/releases
1. Use Homebrew: see https://github.com/ansd/homebrew-tap
1. Build from source: `go get github.com/ansd/driving-time`

## Usage
Create a config file by copying [driving-time.yml.template](driving-time.yml.template) to `driving-time.yml`.

Start the server:
```
$ driving-time --config driving-time.yml serve
```
![serve.png](docs/serve.png)

Print live traffic data on the command line:
![print.png](docs/print.png)
