
package eos_contract_api_client

import (
    "encoding/json"
    "io/ioutil"
    "github.com/imroc/req"
)

type Client struct {
    Url string
}

func New(url string) (*Client) {
    return &Client{
        Url: url,
    }
}

func (c *Client) send(method string, path string) (*req.Resp, error) {
    r := req.New()
    return r.Do(method, c.Url + path)
}

//  GetHealth - Fetches "/health" from API
// ---------------------------------------------------------
func (c *Client) GetHealth() (Health, error) {

    var health Health

    r, err := c.send("GET", "/health")
    if err == nil {
        resp := r.Response()
        body, _ := ioutil.ReadAll(resp.Body)

        // Set HTTPStatusCode
        health.HTTPStatusCode = resp.StatusCode

        // Parse json
        err = json.Unmarshal(body, &health)
    }
    return health, err
}
