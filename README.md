<p align="center">
  <a href="https://blockwallet.io">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://user-images.githubusercontent.com/11839151/188500975-8cd95d07-c419-48aa-bb85-4200a6526f68.svg" />
      <source media="(prefers-color-scheme: light)" srcset="https://blockwallet.io/static/images/logo-blockwallet-black.svg" />
      <img src="[https://blockwallet.io/static/images/logo-medium.svg](https://user-images.githubusercontent.com/11839151/188500975-8cd95d07-c419-48aa-bb85-4200a6526f68.svg)" width="300" />
    </picture>
  </a>
</p>

<hr />

# BlockWallet - Campaigns service

This a service that stores the BlockWallet's incentive driven campaigns along with their participants, rewards and winners.

## Development

### Service

This is a `REST API` with two interfaces: `HTTP` and `gRPC`. It includes metrics, interceptors, logs and uses PostgreSQL to store the campaigns information.

## Getting started

1. Install the generate dependencies `buf` using the package manager of your preference.

If Homebrew package manager, here it is the installation command:

```
brew install buf
```

Otherwise, you can pick the most convinient installation for your environment here: https://buf.build/docs/installation.

2. Make sure your project builds successfully by running `make build`.

## Database configuration

If you're using a local databse, make sure to configure the connection properly by using the Program env variable `SQL_CONNECTION`. If you would like to go with the default configuration, make sure you have a PostgreSQL db running in `localhost:5432` with `user=postgres` and `password=admin`. This is the full default connection string:
`postgresql://localhost:5432/postgres?user=postgres&password=admin&sslmode=disable`

If you don't know how to initialize a database, you can run `make db/up` to have a proper db that fullfils the default connection string configuration. You can stop it without lossing all your's database configuration, in order to stop it, run `make db/down`.

### Migrations

This project runs the migrations automatically when the server starts. If would like to change this behavior indicate `SKIP_MIGRATIONS=true` when running this program.

## Building and running locally

On the base path of this project you can run `make build`. After that, running `./campaignsservice serve` will run both the
gRPC endpoint as well as the REST API in port 8080.

Also, you can can run `make run` and the result will be the same.

You can then test that the server is working issuing `curl localhost:8080/ready`. That must return `YES`.

Alternatively it is possible to generate a `Docker image` and run the service there by running this command:

`docker build --pull --rm -f "Dockerfile" -t campaignsservice:latest "."`

And then this other command:

`docker run --rm -it -p 8080:8080/tcp -p 8443:8443/tcp -p 9008:9008/tcp campaignsservice:latest`

### Program arguments

The server provides a way to modify some default parameters by these env variables: (this is not mandatory, please check
the file `cmd/server/init.go` to check the default values.)

```
LOG_LEVEL -> string [debug|info|error|warning|fatal|panic] (debug)

PORT -> int (8080)
METRICS_PORT -> int (9008)

DB_TYPE -> string [PostgreSQL|SQLite] (PostgreSQL)
SQL_CONNECTION -> string (postgresql://localhost:5432/postgres?user=postgres&password=admin&sslmode=disables)
ADMIN_USERNAME -> string (blockwallet)
ADMIN_PASSWORD -> string (password123)
```

### Public Endpoints

#### Campaigns

This endpoint list all the campaigns applying the indicated filters.

##### Request

`curl http://localhost:8080/v1/api/campaigns`

Note that this endpoint returns ACTIVE campaigns by default. In order to change that behavior, you need to specify the statuses you want to include.

##### Filters:

- `filters.statuses`-> string [CAMPAIGN_STATUS_PENDING | CAMPAIGN_STATUS_ACTIVE | CAMPAIGN_STATUS_FINISHED | CAMPAIGN_STATUS_CANCELLED]
- `filters.tags`-> string
- `filters.fromDate` -> string date [Format: `2006-01-02T15:04:05Z07:00`]
- `filters.toDate` -> string date [Format: `2006-01-02T15:04:05Z07:00`]
- `filters.chain_ids` -> int

##### Response

```
{
    "campaigns": [
        {
            "id": "7fefd5a3-c808-4353-b5a9-98686dfc7fb0",
            "supportedChains": [
                1,
                137
            ],
            "name": "Campaign 2",
            "description": "This is the second campaign for a PostgreSQL db",
            "status": "CAMPAIGN_STATUS_ACTIVE",
            "startDate": "2023-04-01T00:00:00Z",
            "endDate": "2023-06-01T00:00:00Z",
            "rewards": {
                "token": {
                    "id": "",
                    "name": "GoBlank",
                    "decimals": 18,
                    "symbol": "BLANK",
                    "contractAddresses": {}
                },
                "amounts": [
                    "40000000",
                    "30000000",
                    "20000000"
                ],
                "type": "REWARD_TYPE_PODIUM"
            },
            "accounts": [],
            "winners": [],
            "tags": [
                "BLANK",
                "staking2"
            ],
            "enrollMessage": "FYI: This is my cusotm enroll message"
        }
    ]
}
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
