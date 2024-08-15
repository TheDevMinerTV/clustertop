package internal

var (
	NodeMapping = map[string]string{
		"195.90.212.50:9100": "node2",
		"178.254.36.101:9100": "node3",
		"195.90.221.208:9100": "node4",
		"195.90.223.155:9100": "node5",
		"178.254.37.105:9100": "node6",
	}
)

func NodeFromInstance(ip string) string {
	return NodeMapping[ip]
}
