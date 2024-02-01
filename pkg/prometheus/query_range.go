package prometheus

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type QueryRangeResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType QueryResultType `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func (c Client) QueryRange(query string) (*QueryRangeResult, error) {
	now := time.Now()

	u := c.baseUrl
	u.Path = "/api/v1/query_range"

	q := u.Query()
	q.Set("query", query)
	q.Set("start", now.Add(-(5 * time.Minute)).Format(time.RFC3339))
	q.Set("end", now.Format(time.RFC3339))
	q.Set("step", "30s")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to query prometheus: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := QueryRangeResult{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, err
	}

	return &result, nil
}
