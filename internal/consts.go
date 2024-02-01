package internal

var (
	NodeMapping = map[string]string{
		"178.254.36.101:9100": "node3",
		"195.90.221.208:9100": "node4",
		"195.90.223.155:9100": "node5",
		"10.0.99.6:9100":      "node6",
		"10.0.99.7:9100":      "node7",
	}
)

func NodeFromInstance(ip string) string {
	return NodeMapping[ip]
}
