package config

import (
	"flag"
	"os"
)

type Config struct {
	GRPCPort  string
	NATSURL   string
	StoreName string
	DBPath    string
	SyncDir   string
}

func LoadConfig() *Config {
	c := &Config{}

	flag.StringVar(&c.GRPCPort, "grpc-port", getEnv("GRPC_PORT", "50051"), "gRPC server port")
	flag.StringVar(&c.NATSURL, "nats-url", getEnv("NATS_URL", "nats://localhost:4222"), "NATS server URL")
	flag.StringVar(&c.StoreName, "store-name", getEnv("STORE_NAME", "fileSearchStores/agentmeshresearchcore-1jsf1t5e0494"), "Gemini File Search Store name")
	flag.StringVar(&c.DBPath, "db-path", getEnv("DB_PATH", "/home/groovy-byte/agent_mesh.db"), "Path to SQLite database")
	flag.StringVar(&c.SyncDir, "sync-dir", getEnv("SYNC_DIR", "/home/groovy-byte/agent-mesh-core/tmp_sync"), "Directory for sync files")

	flag.Parse()
	return c
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
