# Golang service template

## Summary

This is a template for golang services that need to implement `eth` features.

### Features

- Async subscription to blockchain events
- Support for KV databases: Local db (debug) or Redis Sentinel (production)
- Support for external API calls
- Stability: Logger, metrics and interceptors
- Clustering: Docker & Docker compose

## Service

This is a `REST API` with two interfaces: `HTTP` and `gRPC`. It includes metrics, interceptors, logs, a local kv cache,
a redis implementation and unit tests.

## Getting started

1. Install the generate dependencies with `make deps`. Also, excute the command `make install` in order to
   install `protoc-gen-go`, `protoc-gen-grpc-gateway`, `protoc-gen-doc` which are necessary for us.

2. Make sure your project builds successfully by running `make build`.

## Building and running locally

On the base path of this project you can run `make build`. After that, running `./ethservice serve` will run both the
gRPC endpoint as well as the REST API in port 8080.

Also, you can can run `make run` and the result will be the same.

You can then test that the server is working issuing `curl localhost:8080/ready`. That must return `YES`.

Alternatively it is possible to generate a `Docker image` and run the service there by running this command:

`docker build --pull --rm -f "Dockerfile" -t ethservice:latest "."`

And then this other command:

`docker run --rm -it -p 8080:8080/tcp -p 8443:8443/tcp -p 9008:9008/tcp ethservice:latest`

### Program arguments

The server provides a way to modify some default parameters by these env variables: (this is not mandatory, please check
the file `cmd/server/init.go` to check the default values.)

```
LOG_LEVEL -> string [debug|info|error|warning|fatal|panic]

PORT -> int [8080]
METRICS_PORT -> int [9008]

KV_TYPE -> string [local|redis]
REDIS_SENTINEL_HOSTS -> string (comma separated if multiple values)
REDIS_SENTINEL_MASTER_NAME -> string
REDIS_PASSWORD -> string
REDIS_DB -> int [0|1|2|...]

SOME_HTTP_ENDPOINT -> string
SOME_HTTP_PROTOCOL -> string [http|https|wss]
SOME_HTTP_TIMEOUT -> int

LOCAL_CACHE_EXPIRATION -> int
LOCAL_CACHE_CLEAN_UP_INTERVAL -> int

ETH_ENDPOINT -> string
```

### Endpoints

#### Chains

This is an example about how to consume an external API and return something starting from that.

##### Request

`curl http://localhost:8080/chains`

##### Response

```
{
    "chains": [
        {
            "name": "Ethereum Mainnet",
            "chain": "ETH",
            "network": "mainnet",
            "icon": "ethereum",
            "rpc": [
                "https://mainnet.infura.io/v3/${INFURA_API_KEY}",
                "wss://mainnet.infura.io/ws/v3/${INFURA_API_KEY}",
                "https://api.mycryptoapi.com/eth",
                "https://cloudflare-eth.com"
            ],
            "nativeCurrency": {
                "name": "Ether",
                "symbol": "ETH",
                "decimals": "18"
            },
            "infoURL": "https://ethereum.org",
            "shortName": "eth",
            "chainId": "1",
            "networkId": "1",
            "ens": {
                "registry": "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"
            },
            "explorers": [
                {
                    "name": "etherscan",
                    "url": "https://etherscan.io",
                    "standard": "EIP3091"
                }
            ]
        },
        ...
```

#### Events

This template is subscripted to the Tornado events (deposits and withdrawals) in Polygon network. After every event the
service populates all the necessary data and store it in the KV database initialized.

##### Request

`curl http://localhost:8080/events?pair=ETH`

##### Response

```
{
    "events": [{
        "blockNumber": 4455249
        "commitment": "0x1b53a696c0ddce74074b6aba83495e94983cfd5fbe31e14576236dd051f3963d"
        "leafIndex": 0
        "timestamp": "1615964413"
        "transactionHash": "0xa54ac86c58db7ada9347b635943ff11bd175e50bbef59f38113c6c58ff2bc958"
    }]
}
```

## Develop

- After made changes run `make fmt` and `make lint` to lint your code
- If any proto is modified run `make generate` to generate the expected interface. Then you'll need to implement and
  register the handlers.

### Testing

To run the test excecute `make test`

### Metrics

The `REST` server runs on port `8080` but, in the port `9008` another server runs that includes usage metrics.

You can retrieve the metrics by running `curl localhost:9008/metrics`

### Code format

Before upload any change please run `make fmt` and `make lint` for code formatting and linting


### Redis integration

Redis Sentinel is in charge of monitoring the redis instances and return the address of the current master node and also select a new one if the existing one is not available. 

![Redis Sentinel](https://miro.medium.com/max/855/1*gszoEBW0lupbMDDGGgYOPA.png)

### Production

You should set up a password following this guide https://github.com/spotahome/redis-operator#enabling-redis-auth. Once the secret is created add it to the (redis file)[infra/manifests/redis-cluster.yaml] and (service file)[infra/manifests/main.yaml]. Read `Deployment` below.

### Local development

To use the service with Redis Sentinel you have to use Docker Compose:

```docker
docker-compose up -d --build
```

The previous command will re/build the service image and start the service, redis master and redis sentinel containers. Starting only the redis containers won't work because the sentinel will return the local IP (docker network) of the master node and the service won't be able to connect to it. 

### Deployment

Currently the deployment is manually. TBD: use Flux.

### Prerequisites

1. Install AWS v2 CLI https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

2. Ask for your credentials (with access to the ECRs) in the channel #devops

3. Configure your cli with `aws configure`

4. Install kubectl https://kubernetes.io/docs/tasks/tools/

### Build and push

#### CI/CD Pipeline

- Pushing to main using a PR will build and push a new Docker image using the first 5 characters of the commit as a version.
- Creating a new tag will build and push a new Docker image using the the tag you've just created as a version. ** Use [semver 2.0](https://semver.org/) starting with 'v'**

#### Manual

```
make docker-login
make docker-publish TAG=<semver>
```

### Deploy

Update (at least) the image version in the `dev/main.yaml` or `prod/main.yaml` in [k8s manifest](k8s/) and then:

```
make deploy ENV=[dev|prod]
```


