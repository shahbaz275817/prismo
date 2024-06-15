
# Prismo Transaction Service
This project is a Go-based API for managing cardholder accounts and their transactions. It allows the creation of accounts and transactions, and retrieval of account information. The transactions can be of various types, such as normal purchase, withdrawal, credit voucher, or purchase with installments.


## Setup

This service runs on go. Follow instructions on https://golang.org/doc/install for setting up Go environment.



```
GOPATH=<workspace_dir>
export GOPATH
PATH="${PATH}:${GOPATH}/bin"
export PATH

# Use Makefile commands to setup and run tasks
make setup-local
make clean build
make clean test
```

Checkout the code, install the dependencies and build the project:

```
$ mkdir $GOPATH/bin
$ go get -u https://github.com/shahbaz275817/prismo
$ cd $GOPATH/src/https://github.com/shahbaz275817/prismo
$ make build
```

For running the project:

```
$ make run
```
The API will be available at http://localhost:8080/prismo/v1/



## Requirements

- Go 1.16+
- Make

## Data Structures

### Accounts

| Column          | Type   | Description                 |
|-----------------|--------|-----------------------------|
| `account_id`    | int    | Unique identifier for the account |
| `document_number` | string | Document number for the account holder |

### OperationTypes

| Column            | Type   | Description                  |
|-------------------|--------|------------------------------|
| `operation_type_id` | int    | Unique identifier for the operation type |
| `description`     | string | Description of the operation type |

### Transactions

| Column          | Type   | Description                               |
|-----------------|--------|-------------------------------------------|
| `transaction_id` | int    | Unique identifier for the transaction     |
| `account_id`    | int    | Identifier for the associated account      |
| `operation_type_id` | int    | Identifier for the type of operation     |
| `amount`        | float  | Amount of the transaction (negative for purchases and withdrawals, positive for credit vouchers) |
| `event_date`    | string | Date and time when the transaction occurred |

## API Endpoints

### Create an Account

**Endpoint:** `POST /accounts`

Curl:
```curl
curl --location 'http://localhost:8080/prismo/v1/accounts' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic bG1zX3BvcnRhbDptdFVOVzAzdkdjT0lUVGN6bHlmcld3PT0=' \
--data '{
    "document_number":"abc123"
}'
```

### Get Account

**Endpoint:** `GET /accounts`

Curl:
```curl
curl --location 'http://localhost:8080/prismo/v1/accounts/2' \
--header 'Authorization: Basic bG1zX3BvcnRhbDptdFVOVzAzdkdjT0lUVGN6bHlmcld3PT0='
```

### Create a Transaction

**Endpoint:** `POST /transactions`

Curl:
```curl
curl --location 'http://localhost:8080/prismo/v1/transactions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic bG1zX3BvcnRhbDptdFVOVzAzdkdjT0lUVGN6bHlmcld3PT0=' \
--data '{
"account_id": 1,
"operation_type_id": 4,
"amount": 213.45
}'
```