
package eos_contract_api_client

import (
    "time"
    "testing"
    "net/http"
    "net/http/httptest"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

var asset1 = Asset{
    ID: "1099667509880",
    Contract: "atomicassets",
    Owner: "farmersworld",
    Name: "Silver Member",
    IsTransferable: true,
    IsBurnable: true,
    TemplateMint: "4433",
    Collection: Collection{
        CollectionName: "farmersworld",
        Name: "Farmers World",
        Author: ".jieg.wam",
        AllowNotify: true,
        AuthorizedAccounts: []string{
            ".jieg.wam",
            "farmersworld",
            "atomicdropsx",
            "atomicpacksx",
            "neftyblocksd",
        },
        NotifyAccounts: []string{
            "atomicdropsx",
        },
        MarketFee: 0.05,
        CreatedAtBlock: "123762633",
        CreatedAtTime: "1623323058000",
    },
    Schema: Schema{
        Name: "memberships",
        Contract: "",
        Format: []SchemaFormat{
            SchemaFormat{
              Name: "name",
              Type: "string",
            },
            SchemaFormat{
              Name: "img",
              Type: "image",
            },
            SchemaFormat{
              Name: "description",
              Type: "string",
            },
            SchemaFormat{
              Name: "type",
              Type: "string",
            },
            SchemaFormat{
              Name: "rarity",
              Type: "string",
            },
            SchemaFormat{
              Name: "level",
              Type: "uint8",
            },
        },
        CreatedAtBlock: "136880914",
        CreatedAtTime: "1629887699000",
    },
    Template: Template{
        ID: "260629",
        MaxSupply: "0",
        IsTransferable: true,
        IsBurnable: true,
        IssuedSupply: "112195",
        ImmutableData: map[string]interface{}{
            "img": "QmZWg1mP2UNcSwhrYNVqjk16BnhcWCz3oAva8BfiTNB3J4",
            "name": "Silver Member",
            "type": "Wood",
            "level": float64(2),
            "rarity": "Uncommon",
            "description": "This is a member card powered by Wood. When used by the farmer, it will increase the power and luck of the wood mining tools, and can mine the Farmer Coin that has been lost since ancient times.",
        },
        CreatedAtBlock: "136882467",
        CreatedAtTime: "1629888476000",
    },
    ImmutableData: map[string]interface{}{
        "asdx": "4321",
    },
    MutableData: map[string]interface{}{
        "asdf": "1234",
    },
    UpdatedAtBlock: "171080009",
    UpdatedAtTime: "1646996870500",
    TransferedAtBlock: "171080009",
    TransferedAtTime: "1646996870500",
    MintedAtBlock: "171080009",
    MintedAtTime: "1646996870500",
    BackedTokens: []Token{},
}

func TestGetHealth(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/health" {
            payload := `{
                "success":true,
                "data":{
                    "version":"1.0.0",
                    "postgres":{
                        "status":"OK",
                        "readers":[
                            {
                                "block_num":"167836036"
                            },
                            {
                                "block_num":"167836034"
                            }
                        ]
                    },
                    "redis":{
                        "status":"OK"
                    },
                    "chain":{
                        "status":"OK",
                        "head_block":167836035,
                        "head_time":1645374771500
                    }
                },
                "query_time":1645374772067
            }`

            res.Header().Add("Content-type", "application/json; charset=utf-8")
            res.Write([]byte(payload))
        }
    }))

    client := New(srv.URL)

    h, err := client.GetHealth()

    require.NoError(t, err)
    assert.Equal(t, 200, h.HTTPStatusCode)

    assert.True(t, h.Success)
    assert.Equal(t, time.Time(time.Date(2022, time.February, 20, 16, 32, 52, 67, time.UTC)), h.QueryTime)

    // Data
    assert.Equal(t, "1.0.0", h.Data.Version)

    // Postgres
    assert.Equal(t, "OK", h.Data.Postgres.Status)

    // Redis
    assert.Equal(t, "OK", h.Data.Redis.Status)

    // Chain
    assert.Equal(t, "OK", h.Data.Chain.Status)
    assert.Equal(t, int64(167836035), h.Data.Chain.HeadBlock)
    assert.Equal(t, time.Time(time.Date(2022, time.February, 20, 16, 32, 51, 500, time.UTC)), h.Data.Chain.HeadTime)
}

func TestGetHealthFailed(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/health" {
            payload := `{
                "success":true,
                "data":{
                    "version":"1.0.0",
                    "postgres":{
                        "status":"ERROR",
                        "readers":[]
                    },
                    "redis":{
                        "status":"ERROR"
                    },
                    "chain":{
                        "status":"ERROR",
                        "head_block":0,
                        "head_time":0
                    }
                },
                "query_time":1645374772067
            }`

            res.Header().Add("Content-type", "application/json")
            res.Write([]byte(payload))
        }
    }))

    client := New(srv.URL)

    h, err := client.GetHealth()

    require.NoError(t, err)
    assert.Equal(t, 200, h.HTTPStatusCode)

    assert.True(t, h.Success)
    assert.Equal(t, time.Time(time.Date(2022, time.February, 20, 16, 32, 52, 67, time.UTC)), h.QueryTime)

    // Data
    assert.Equal(t, "1.0.0", h.Data.Version)

    // Postgres
    assert.Equal(t, "ERROR", h.Data.Postgres.Status)

    // Redis
    assert.Equal(t, "ERROR", h.Data.Redis.Status)

    // Chain
    assert.Equal(t, "ERROR", h.Data.Chain.Status)
    assert.Equal(t, int64(0), h.Data.Chain.HeadBlock)

    assert.Equal(t, time.Unix(0, 0).UTC(), h.Data.Chain.HeadTime)
}

func TestAPIError(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        payload := `{
          "success": false,
          "message": "Some internal error"
        }`

        res.Header().Add("Content-type", "application/json")
        res.WriteHeader(500)
        res.Write([]byte(payload))
    }))

    client := New(srv.URL)

    _, err := client.GetHealth()

    assert.Error(t, err)
}


func TestHostHeader(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        assert.Equal(t, "my-custom-host", req.Host)
    }))

    client := New(srv.URL)
    client.Host = "my-custom-host"

    client.send("GET", "/", nil)
}

func TestInvalidContentType(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        res.Header().Add("Content-type", "some-type")
    }))

    client := New(srv.URL)

    _, err := client.send("GET", "/", nil)

    assert.EqualError(t, err, "Invalid content-type 'some-type', expected 'application/json'")
}

func TestGetAsset(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/attomicaassets/v1/assets/1099667509880" {
            payload := `{
              "success": true,
              "data": {
                "contract": "atomicassets",
                "asset_id": "1099667509880",
                "owner": "farmersworld",
                "is_transferable": true,
                "is_burnable": true,
                "collection": {
                  "collection_name": "farmersworld",
                  "name": "Farmers World",
                  "img": "QmX79zrJsk4DbWQ3krgu41pX3fdvEvWjkMXiNCKpxFXSgj",
                  "author": ".jieg.wam",
                  "allow_notify": true,
                  "authorized_accounts": [
                    ".jieg.wam",
                    "farmersworld",
                    "atomicdropsx",
                    "atomicpacksx",
                    "neftyblocksd"
                  ],
                  "notify_accounts": [
                    "atomicdropsx"
                  ],
                  "market_fee": 0.05,
                  "created_at_block": "123762633",
                  "created_at_time": "1623323058000"
                },
                "schema": {
                  "schema_name": "memberships",
                  "format": [
                    {
                      "name": "name",
                      "type": "string"
                    },
                    {
                      "name": "img",
                      "type": "image"
                    },
                    {
                      "name": "description",
                      "type": "string"
                    },
                    {
                      "name": "type",
                      "type": "string"
                    },
                    {
                      "name": "rarity",
                      "type": "string"
                    },
                    {
                      "name": "level",
                      "type": "uint8"
                    }
                  ],
                  "created_at_block": "136880914",
                  "created_at_time": "1629887699000"
                },
                "template": {
                  "template_id": "260629",
                  "max_supply": "0",
                  "is_transferable": true,
                  "is_burnable": true,
                  "issued_supply": "112195",
                  "immutable_data": {
                    "img": "QmZWg1mP2UNcSwhrYNVqjk16BnhcWCz3oAva8BfiTNB3J4",
                    "name": "Silver Member",
                    "type": "Wood",
                    "level": 2,
                    "rarity": "Uncommon",
                    "description": "This is a member card powered by Wood. When used by the farmer, it will increase the power and luck of the wood mining tools, and can mine the Farmer Coin that has been lost since ancient times."
                  },
                  "created_at_time": "1629888476000",
                  "created_at_block": "136882467"
                },
                "mutable_data": {
                    "asdf": "1234"
                },
                "immutable_data": {
                    "asdx": "4321"
                },
                "template_mint": "4433",
                "backed_tokens": [],
                "burned_by_account": null,
                "burned_at_block": null,
                "burned_at_time": null,
                "updated_at_block": "171080009",
                "updated_at_time": "1646996870500",
                "transferred_at_block": "171080009",
                "transferred_at_time": "1646996870500",
                "minted_at_block": "171080009",
                "minted_at_time": "1646996870500",
                "data": {
                  "img": "QmZWg1mP2UNcSwhrYNVqjk16BnhcWCz3oAva8BfiTNB3J4",
                  "name": "Silver Member",
                  "type": "Wood",
                  "level": 2,
                  "rarity": "Uncommon",
                  "description": "This is a member card powered by Wood. When used by the farmer, it will increase the power and luck of the wood mining tools, and can mine the Farmer Coin that has been lost since ancient times."
                },
                "name": "Silver Member"
              },
              "query_time": 1647016614598
            }`

            res.Header().Add("Content-type", "application/json; charset=utf-8")
            res.Write([]byte(payload))
        }
    }))

    client := New(srv.URL)

    a, err := client.GetAsset("1099667509880")

    require.NoError(t, err)
    assert.Equal(t, 200, a.HTTPStatusCode)
    assert.True(t, a.Success)
    assert.Equal(t, time.Time(time.Date(2022, time.March, 11, 16, 36, 54, 598, time.UTC)), a.QueryTime)
    assert.Equal(t, asset1, a.Data)
}

func TestGetAssets(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/attomicaassets/v1/assets?before=100&is_transferable=true&schema_name=test" {
            payload := `{
  "success": true,
  "data": [
    {
      "contract": "atomicassets",
      "asset_id": "1099667509880",
      "owner": "farmersworld",
      "is_transferable": true,
      "is_burnable": true,
      "collection": {
        "collection_name": "farmersworld",
        "name": "Farmers World",
        "img": "QmX79zrJsk4DbWQ3krgu41pX3fdvEvWjkMXiNCKpxFXSgj",
        "author": ".jieg.wam",
        "allow_notify": true,
        "authorized_accounts": [
          ".jieg.wam",
          "farmersworld",
          "atomicdropsx",
          "atomicpacksx",
          "neftyblocksd"
        ],
        "notify_accounts": [
          "atomicdropsx"
        ],
        "market_fee": 0.05,
        "created_at_block": "123762633",
        "created_at_time": "1623323058000"
      },
      "schema": {
        "schema_name": "memberships",
        "format": [
          {
            "name": "name",
            "type": "string"
          },
          {
            "name": "img",
            "type": "image"
          },
          {
            "name": "description",
            "type": "string"
          },
          {
            "name": "type",
            "type": "string"
          },
          {
            "name": "rarity",
            "type": "string"
          },
          {
            "name": "level",
            "type": "uint8"
          }
        ],
        "created_at_block": "136880914",
        "created_at_time": "1629887699000"
      },
      "template": {
        "template_id": "260629",
        "max_supply": "0",
        "is_transferable": true,
        "is_burnable": true,
        "issued_supply": "112195",
        "immutable_data": {
          "img": "QmZWg1mP2UNcSwhrYNVqjk16BnhcWCz3oAva8BfiTNB3J4",
          "name": "Silver Member",
          "type": "Wood",
          "level": 2,
          "rarity": "Uncommon",
          "description": "This is a member card powered by Wood. When used by the farmer, it will increase the power and luck of the wood mining tools, and can mine the Farmer Coin that has been lost since ancient times."
        },
        "created_at_time": "1629888476000",
        "created_at_block": "136882467"
      },
      "mutable_data": {
          "asdf": "1234"
      },
      "immutable_data": {
          "asdx": "4321"
      },
      "template_mint": "4433",
      "backed_tokens": [],
      "burned_by_account": null,
      "burned_at_block": null,
      "burned_at_time": null,
      "updated_at_block": "171080009",
      "updated_at_time": "1646996870500",
      "transferred_at_block": "171080009",
      "transferred_at_time": "1646996870500",
      "minted_at_block": "171080009",
      "minted_at_time": "1646996870500",
      "data": {
        "img": "QmZWg1mP2UNcSwhrYNVqjk16BnhcWCz3oAva8BfiTNB3J4",
        "name": "Silver Member",
        "type": "Wood",
        "level": 2,
        "rarity": "Uncommon",
        "description": "This is a member card powered by Wood. When used by the farmer, it will increase the power and luck of the wood mining tools, and can mine the Farmer Coin that has been lost since ancient times."
      },
      "name": "Silver Member"
    }],
    "query_time":1646996870918
    }`

            res.Header().Add("Content-type", "application/json; charset=utf-8")
            res.Write([]byte(payload))
        }
    }))

    client := New(srv.URL)

    a, err := client.GetAssets(AssetsRequestParams{Before: 100, SchemaName: "test", IsTransferable: true})

    require.NoError(t, err)
    assert.Equal(t, 200, a.HTTPStatusCode)
    assert.True(t, a.Success)
    assert.Equal(t, time.Time(time.Date(2022, time.March, 11, 11, 7, 50, 918, time.UTC)), a.QueryTime)

    expected := []Asset{asset1}

    assert.Equal(t, expected, a.Data)
}
