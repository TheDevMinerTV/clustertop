package scraper

import (
	"stats.k8s.devminer.xyz/internal/cache"
	"stats.k8s.devminer.xyz/pkg/prometheus"
)

var (
	cpuUsageQuery        = "sum(rate(node_cpu_seconds_total{mode!=\"idle\"}[1m])) by (instance)"
	cpuCoresQuery        = "count(node_cpu_seconds_total{mode=\"idle\"}) by (instance)"
	memoryUsageQuery     = "node_memory_MemTotal_bytes - (node_memory_MemFree_bytes + node_memory_Buffers_bytes + node_memory_Cached_bytes)"
	memoryTotalQuery     = "node_memory_MemTotal_bytes"
	maxReceiveSpeedMbits = 250.0 // Mbit/s
	maxTransmitRateMbits = 600.0 // Mbit/s
	networkRxQuery       = "rate(node_network_receive_bytes_total{device=\"ens3\"}[1m])"
	networkTxQuery       = "rate(node_network_transmit_bytes_total{device=\"ens3\"}[1m])"
)

type Scraper struct {
	client *prometheus.Client
}

func New(client *prometheus.Client) *Scraper {
	return &Scraper{client: client}
}

func (u *Scraper) Scrape() (map[string]cache.Node, error) {
	newNodes := make(map[string]cache.Node)

	cpuUsage, err := u.fetchCpuUsage()
	if err != nil {
		return nil, err
	}

	cpuCores, err := u.fetchCpuCores()
	if err != nil {
		return nil, err
	}

	memoryUsage, err := u.fetchMemoryUsage()
	if err != nil {
		return nil, err
	}

	memoryTotal, err := u.fetchMemoryTotal()
	if err != nil {
		return nil, err
	}

	networkRx, err := u.fetchNetworkRx()
	if err != nil {
		return nil, err
	}

	networkTx, err := u.fetchNetworkTx()
	if err != nil {
		return nil, err
	}

	for node, v := range cpuUsage {
		newNodes[node] = cache.Node{
			CPU: cache.Value{V1: v * 100},
		}
	}

	for node, v := range cpuCores {
		current, ok := newNodes[node]
		if !ok {
			continue
		}

		newNodes[node] = cache.Node{
			CPU: cache.Value{
				V1: current.CPU.V1,
				V2: v * 100,
			},
		}
	}

	for node, v := range memoryUsage {
		current, ok := newNodes[node]
		if !ok {
			continue
		}

		newNodes[node] = cache.Node{
			CPU:    current.CPU,
			Memory: cache.Value{V1: toGB(v)},
		}
	}

	for node, v := range memoryTotal {
		current, ok := newNodes[node]
		if !ok {
			continue
		}

		newNodes[node] = cache.Node{
			CPU:    current.CPU,
			Memory: cache.Value{V1: current.Memory.V1, V2: toGB(v)},
		}
	}

	for node, v := range networkRx {
		current, ok := newNodes[node]
		if !ok {
			continue
		}

		newNodes[node] = cache.Node{
			CPU:             current.CPU,
			Memory:          current.Memory,
			NetworkReceive:  cache.Value{V1: toMbits(v), V2: maxReceiveSpeedMbits},
			NetworkTransmit: current.NetworkTransmit,
		}
	}

	for node, v := range networkTx {
		current, ok := newNodes[node]
		if !ok {
			continue
		}

		newNodes[node] = cache.Node{
			CPU:             current.CPU,
			Memory:          current.Memory,
			NetworkReceive:  current.NetworkReceive,
			NetworkTransmit: cache.Value{V1: toMbits(v), V2: maxTransmitRateMbits},
		}
	}

	return newNodes, nil
}

func (u *Scraper) fetchCpuUsage() (map[string]float64, error) {
	res, err := u.client.Query(cpuUsageQuery)
	if err != nil {
		return nil, err
	}

	values := make(map[string]float64)

	for _, v := range res.Data.Result {
		values[v.Node()] = v.MustValue()
	}

	return values, nil
}

func (u *Scraper) fetchCpuCores() (map[string]float64, error) {
	res, err := u.client.Query(cpuCoresQuery)
	if err != nil {
		return nil, err
	}

	values := make(map[string]float64)

	for _, v := range res.Data.Result {
		values[v.Node()] = v.MustValue()
	}

	return values, nil
}

func (u *Scraper) fetchMemoryUsage() (map[string]float64, error) {
	res, err := u.client.Query(memoryUsageQuery)
	if err != nil {
		return nil, err
	}

	values := make(map[string]float64)

	for _, v := range res.Data.Result {
		values[v.Node()] = v.MustValue()
	}

	return values, nil
}

func (u *Scraper) fetchMemoryTotal() (map[string]float64, error) {
	res, err := u.client.Query(memoryTotalQuery)
	if err != nil {
		return nil, err
	}

	values := make(map[string]float64)

	for _, v := range res.Data.Result {
		values[v.Node()] = v.MustValue()
	}

	return values, nil
}

func (u *Scraper) fetchNetworkRx() (map[string]float64, error) {
	res, err := u.client.Query(networkRxQuery)
	if err != nil {
		return nil, err
	}

	values := make(map[string]float64)

	for _, v := range res.Data.Result {
		values[v.Node()] = v.MustValue()
	}

	return values, nil
}

func (u *Scraper) fetchNetworkTx() (map[string]float64, error) {
	res, err := u.client.Query(networkTxQuery)
	if err != nil {
		return nil, err
	}

	values := make(map[string]float64)

	for _, v := range res.Data.Result {
		values[v.Node()] = v.MustValue()
	}

	return values, nil
}

func toGB(v float64) float64 {
	return v / 1024 / 1024 / 1024
}

func toMbits(v float64) float64 {
	return v / 1024 / 1024 / 8
}
