# coin-price-backend

[![last commit](https://img.shields.io/github/last-commit/noah-blockchain/coin-price-backend.svg)]()
[![license](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://github.com/noah-blockchain/coin-price-backend/blob/master/LICENSE)
[![version](https://img.shields.io/github/tag/noah-blockchain/coin-price-backend.svg)](https://github.com/noah-blockchain/coin-price-backend/releases/latest)
[![](https://tokei.rs/b1/github/noah-blockchain/coin-price-backend?category=lines)](https://github.com/noah-blockchai/coin-price-backend)

### How To Run This Project
> Make Sure you have run the coins.sql in your postgres

Since the project already use Go Module, I recommend to put the source code in any folder but GOPATH.

#### Run the Testing

```bash
$ make test
```

#### Run the Applications
Here is the steps to run it with `docker-compose`

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone https://github.com/noah-blockchain/coin-price-backend.git

#move to project
$ cd coin-price-backend

# Build the docker image first
$ make docker

# Run the application
$ make run

# check if the containers are running
$ docker ps

# Execute the call
$ curl localhost:10500/price/NOAH

# Stop
$ make stop
```
