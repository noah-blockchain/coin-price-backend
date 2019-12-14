# coin-price-backend
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
