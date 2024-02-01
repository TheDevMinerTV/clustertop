package cache

import (
	"github.com/barweiss/go-tuple"
	"sync"
)

type Cache struct {
	m sync.RWMutex

	v map[string]Node
}

type Node struct {
	CPU             Value `json:"cpu"`
	Memory          Value `json:"memory"`
	NetworkTransmit Value `json:"network_transmit"`
	NetworkReceive  Value `json:"network_receive"`
}

type Value = tuple.T2[float64, float64]

func New() *Cache {
	return &Cache{
		m: sync.RWMutex{},
		v: make(map[string]Node),
	}
}

type Updater = func() (newNodes map[string]Node, err error)

func (c *Cache) Update(fetch func() (newNodes map[string]Node, err error)) error {
	c.m.Lock()
	defer c.m.Unlock()

	newNodes, err := fetch()
	if err != nil {
		return err
	}

	c.v = newNodes

	return nil
}

func (c *Cache) Get() map[string]Node {
	c.m.RLock()
	defer c.m.RUnlock()

	n := make(map[string]Node)
	for k, v := range c.v {
		n[k] = v
	}

	return n
}
