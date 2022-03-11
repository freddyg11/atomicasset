
package eos_contract_api_client

import (
    "time"
)

// HTTP

type HTTPResponse struct {
    HTTPStatusCode int
}

type APIError struct {
    Success bool    `json:"success"`
    Message string  `json:"message"`
}

// Health

type ChainHealth struct {
    Status      string
    HeadBlock   int64
    HeadTime    time.Time
}

type RedisHealth struct {
    Status string `json:"status"`
}

type PostgresHealth struct {
    Status string                       `json:"status"`
    Readers []map[string]interface{}    `json:"readers"`
}

type HealthData struct {
    Version string          `json:"version"`
    Postgres PostgresHealth `json:"postgres"`
    Redis RedisHealth       `json:"redis"`
    Chain ChainHealth       `json:"chain"`
}

type Health struct {
    HTTPResponse
    Success bool
    Data HealthData
    QueryTime time.Time
}
