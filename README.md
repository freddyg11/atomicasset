# EOSIO Contract API Client

This package aims to implement a client for Pinknetwork's [eosio contract api](https://github.com/pinknetworkx/eosio-contract-api) in go.

### Install package

``` bash
go get -u github.com/eosswedenorg-go/eosio-contract-api-client@latest
```

### Types

#### API Client struct

```go
type Client struct {
    Url string
}
```

#### Health struct

```go
type Health struct {
    HTTPStatusCode int
    Success bool
    Data HealthData
    QueryTime time.Time
}

```

### Functions

```go
func New(url string) *Client
```

Construct a new API Client

```go
func (c Client) GetHealth(params ReqParams) (Health, error)
```

Call `/health` and return the results.

### Author

Henrik Hautakoski - [Sw/eden](https://eossweden.org/) - [henrik@eossweden.org](mailto:henrik@eossweden.org)
