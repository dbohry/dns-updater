# DNS Updater
This is a simple Go program that updates DNS records periodically.

### Overview
The program fetches the current public IP address from the specified IP provider and updates the DNS record using the specified target URL. It runs in a loop, updating the DNS record every 10 minutes.

### Usage

#### Environment Variables
- `TARGET_URL`: The URL where the DNS record should be updated.
- `IP_PROVIDER`: (optional) The URL of the service providing the public IP address. Default is https://api.ipify.org.

### Dependencies
- Go 1.22 or higher

### Docker

https://hub.docker.com/repository/docker/dbohry/dns-updater/general

#### Compiling
``
CGO_ENABLED=0 GOOS=linux go build -a -o app .
``

#### Build local docker image
``
docker build -t dns-updater:local .
``

#### Run
``
docker run --rm -e TARGET_URL={url} dns-updater:local
``