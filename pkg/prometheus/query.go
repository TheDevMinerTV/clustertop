package prometheus

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type QueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType QueryResultType `json:"resultType"`
		Result     []InstantVector `json:"result"`
	} `json:"data"`
}

func (c Client) Query(query string) (*QueryResult, error) {
	u := c.baseUrl
	u.Path = "/api/v1/query"

	q := u.Query()
	q.Set("query", query)
	q.Set("time", time.Now().Format(time.RFC3339))
	q.Set("timeout", "5s")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to query Prometheus: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := QueryResult{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, err
	}

	return &result, nil
}
