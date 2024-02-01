package prometheus

import (
	tuple "github.com/barweiss/go-tuple"
	"stats.k8s.devminer.xyz/internal"
	"strconv"
)

type QueryResultType string

const (
	QueryResultTypeMatrix QueryResultType = "matrix"
	QueryResultTypeVector QueryResultType = "vector"
	QueryResultTypeScalar QueryResultType = "scalar"
	QueryResultTypeString QueryResultType = "string"
)

func (t QueryResultType) String() string {
	return string(t)
}

type InstantVector struct {
	Metric map[string]string `json:"metric"`
	// TODO: replace with `go-tuple`
	RawValue tuple.T2[int64, string] `json:"value"`
}

func (v InstantVector) Value() (float64, error) {
	value, err := strconv.ParseFloat(v.RawValue.V2, 64)

	return value, err
}

func (v InstantVector) MustValue() float64 {
	value, err := v.Value()
	if err != nil {
		panic(err)
	}

	return value
}

func (v InstantVector) Time() int64 {
	return v.RawValue.V1
}

func (v InstantVector) Node() string {
	return internal.NodeFromInstance(v.Metric["instance"])
}
