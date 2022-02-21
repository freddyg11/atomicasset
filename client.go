
package eos_contract_api_client

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "github.com/imroc/req"
)

type Client struct {
    Url string
    Host string
}

func New(url string) (*Client) {
    return &Client{
        Url: url,
    }
}

func (c *Client) send(method string, path string) (*req.Resp, error) {
    hdr := make(http.Header)
    r := req.New()

    if len(c.Host) > 0 {
        hdr.Add("Host", c.Host)
    }

    return r.Do(method, c.Url + path, hdr)
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
