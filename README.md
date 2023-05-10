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

## Summary

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

### Using native OS env

On the base path of this project you can run `make build`. After that, running `./campaignsservice serve` will run both the
gRPC endpoint as well as the REST API in port 8080.

Also, you can can run `make run` and the result will be the same.

You can then test that the server is working issuing `curl localhost:8080/ready`. That must return `YES`.

### Using docker

You can use `docker-compose` to run the campaignsservice and the postgres database without any extra configuration. To do so, run:

- `make dev/up` -> Run postgres database and service. Run `make dev/down` to stop it.
- `make dev/run` -> Run postgres database, build and run the campaigns service using `latest` tag. Run `make dev/down` to stop it.
- `make service/up` -> Run just the campaignsservice. Run `make service/down` to stop it.
- `make db/up` -> Run just the postgres database. Run `make db/down` to stop it.

If you don't want to use `docker-compose` it is possible to generate a `Docker image` manually and run the service there by running this command:

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

<hr />

### Public Endpoints

These endpoints are accessible by everybody and they are meant to be used by dApps and the BlockWallet extension.

#### GetCampaigns

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

#### CampaignById

This template is subscripted to the Tornado events (deposits and withdrawals) in Polygon network. After every event the
service populates all the necessary data and store it in the KV database initialized.

##### Request

`curl http://localhost:8080/v1/api/campaigns/7fefd5a3-c808-4353-b5a9-98686dfc7fb0`

##### Response

```
{
    "campaign":  {
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
}
```

#### Enroll account in campaign

The accounts that wants to enroll in a campaign should sign the campaign `enrollMessage` and provide the signature hash in order to make the registration valid. You can directly use the below endpoint `GetCampaignEnrollMessage` to just get the message to sign.

##### Request

`curl -X POST --location 'http://localhost:8080/v1/api/campaigns/7fefd5a3-c808-4353-b5a9-98686dfc7fb0/enroll' \
--header 'Content-Type: application/json' \
--data '{
    "account_address":"0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
    "signature":"0xc7f3b054f0ed23f04d1214c3f35584c04994db5699da1f5e269b7304ee0efecc3bd219cb2e188ec4f353492418ea1fe4dc99efc3164fe44027dd7183405b93f01b"
}'`

#### GetCampaignEnrollMessage

To reduce the request response size, you can directly get the cmapaing's enroll message

##### Request

`curl http://localhost:8080/v1/api/campaigns/7fefd5a3-c808-4353-b5a9-98686dfc7fb0/enroll-message`

##### Response

```
{
    "message": "FYI: This is my cusotm enroll message"
}
```

#### GetCampaignAccounts

To reduce the request response size, you can directly get the campaigns participant

##### Request

`curl http://localhost:8080/v1/api/campaigns/7fefd5a3-c808-4353-b5a9-98686dfc7fb0/accounts`

##### Response

```
{
    "accounts": [
        "0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22"
    ]
}
```

#### GetTokens

Get all the tokens configured

##### Request

`curl http://localhost:8080/v1/api/tokens`

##### Response

```
{
    "tokens": [
        {
            "id": "99ae3179-a06f-4c0d-92f5-6e44f1d3308a",
            "name": "GoBlank",
            "decimals": 18,
            "symbol": "",
            "contractAddresses": {
                "1": "0x41A3Dba3D677E573636BA691a70ff2D606c29666",
                "137": "0xf4C83080E80AE530d6f8180572cBbf1Ac9D5d435"
            }
        }
    ]
}
```

#### GetTokenById

##### Request

`curl http://localhost:8080/v1/api/tokens/99ae3179-a06f-4c0d-92f5-6e44f1d3308a`

##### Response

```
{
    "token": [
        {
            "id": "99ae3179-a06f-4c0d-92f5-6e44f1d3308a",
            "name": "GoBlank",
            "decimals": 18,
            "symbol": "",
            "contractAddresses": {
                "1": "0x41A3Dba3D677E573636BA691a70ff2D606c29666",
                "137": "0xf4C83080E80AE530d6f8180572cBbf1Ac9D5d435"
            }
        }
    ]
}
```

### Private Endpoints

There are endpoints that are restricted to admin access. We use a Basic autentication credentials set on the request headers.

`Authorization: Basic {credentials}`

The credentials are of the shape: `username:password` based-64 encoded.

#### CreateCampaign

Admins may use this enpoint to create campaigns. Note that depending on the information the admin has, the body of this request may change.

For instance, in order to create a campaign, you should specify the reward token you want to use. If the token has already been used in campaign and the admin has its `id` (can be grabbed from the `GetTokens` public endpoint), he can just specify the `id` in the request, otherwise the admin should specify the whole information in order to create the token in the same flow.

Also, there are some restrictions regarding the campaign status, whether it is active by default or not. Here are the things admins should pay attention:

- You can create finished campaings (end_date after than today).
- You can activate a campaign that hasn't started yet (start_date after than today).

Last but not least, the campaign's enroll message can be specified in the request, otherwise a message will be auto-generated by the service using a prefix + the campaigns name. For instance, if the prefix is `Please sign this message in order to enroll in the` and the campaigns name is `Staking campaign`, the final enroll message will be: `Please sign this message in order to enroll in the Staking campaign`.

##### Request

`curl -X POST --location 'http://localhost:8080/v1/admin/campaigns'` -d /

```
{
    "campaign": {
        "name":"Campaign 2",
        "description":"This is an active campaign",
        "is_active":true,
        "start_date":"2023-04-01T00:00:00Z",
        "end_date":"2023-06-01T00:00:00Z",
        "rewards":{
            "amounts":["40000000","30000000","20000000"],
            "type":"REWARD_TYPE_PODIUM",
            "token":{
                "create": {
                    "name":"GoBlank",
                    "symbol":"BLANK",
                    "decimals":18,
                    "contract_addresses":{
                        "1":"0x41A3Dba3D677E573636BA691a70ff2D606c29666",
                        "137":"0xf4C83080E80AE530d6f8180572cBbf1Ac9D5d435"
                    }
                }
            }
        },
        "tags":["BLANK","staking1"],
        "supported_chains": [1,137]
    }
}
```

If you want to specify the `token_id` you should only remove the `rewards.create` and add a new `rewards.id` property with the desired `token_id`.

##### Response

```
{
    "campaign": {
        "id": "f3fa1d90-362a-4674-8515-25d5c8b50aef",
        "supportedChains": [
            1,
            137
        ],
        "name": "Campaign 2",
        "description": "This is an active campaign",
        "status": "CAMPAIGN_STATUS_ACTIVE",
        "startDate": "2023-04-01T00:00:00Z",
        "endDate": "2023-06-01T00:00:00Z",
        "rewards": {
            "token": {
                "id": "",
                "name": "GoBlank",
                "decimals": 18,
                "symbol": "BLANK",
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
        "enrollMessage": "Sign this message to enroll in Campaign 2"
    }
}
```

#### UpdateCampaign

Admins can only update campaigns status and in case of updating the status to `FINISHED` they can also specify the winners. The winners may change depending on the campaign's reward type, where those campaigns which rewards are of they type `PODIUM` the amount of winners should match the quantity of `amounts`.

Possible transitions:

- `PENDING` -> `ACTIVE` (Campaign's start_date should be after now and end_date should be before the current datetime)
- `PENDING` -> `CANCELLED`
- `ACTIVE` -> `CANCELLED`
- `ACTIVE` -> `FINISHED` (must specify winners)
- `FINISHED` -> `FINISHED` (must specify winners again)

If some of these transitions cannot fulfill what you want to do with the campaign, you should `CANCEL` it and create a new one.

##### Request

`curl -X PATCH --location 'http://localhost:8080/v1/admin/campaigns/f3fa1d90-362a-4674-8515-25d5c8b50aef' -d /`

```
{
    "stauts": "CAMPAIGN_STATUS_FINISHED",
    "winners":["0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22","0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C21","0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C32"]
}
```

<hr />

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

### Production

### Deployment

This service is deployed automatically when a new version is pushed to the aws ECR. Please refer to https://github.com/block-wallet/block-devops repository to see the k8s configuration.

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
