
package eos_contract_api_client

import (
    "fmt"
    "strings"
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

func isContentType(t string, expected string) bool {

    p := strings.IndexByte(t, ';')
    if p >= 0 {
        t = t[:p]
    }
    return t == expected
}

func (c *Client) send(method string, path string) (*req.Resp, error) {
    hdr := make(http.Header)
    r := req.New()

    if len(c.Host) > 0 {
        hdr.Add("Host", c.Host)
    }

    resp, err := r.Do(method, c.Url + path, hdr)
    if err != nil {
        return nil, err
    }

    t := resp.Response().Header.Get("Content-type")
    if ! isContentType(t, "application/json") {
        return nil, fmt.Errorf("Invalid content-type '%s', expected 'application/json'", t)
    }

    return resp, err
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
