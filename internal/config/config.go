package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the top-level configuration for x402fuel.
type Config struct {
	Wallet    WalletConfig    `yaml:"wallet"`
	Proxy     ProxyConfig     `yaml:"proxy"`
	Server    ServerConfig    `yaml:"server"`
	Chain     ChainConfig     `yaml:"chain"`
	Events    EventsConfig    `yaml:"events"`
	Telemetry TelemetryConfig `yaml:"telemetry"`
}

// WalletConfig holds key storage settings.
type WalletConfig struct {
	KeyDir string `yaml:"key_dir"`
}

// ProxyConfig holds the HTTP 402 proxy gateway settings.
type ProxyConfig struct {
	ListenAddr string   `yaml:"listen_addr"`
	MaxPerTxn  string   `yaml:"max_per_txn"`
	DailyCap   string   `yaml:"daily_cap"`
	Allowlist  []string `yaml:"allowlist"`
}

// ServerConfig holds the daemon API server settings.
type ServerConfig struct {
	Port int `yaml:"port"`
}

// ChainConfig holds the blockchain connection settings.
type ChainConfig struct {
	RPCURL  string `yaml:"rpc_url"`
	ChainID int    `yaml:"chain_id"`
}

// EventsConfig holds event logging settings.
type EventsConfig struct {
	LogPath string `yaml:"log_path"`
}

// TelemetryConfig holds telemetry reporting settings.
type TelemetryConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return Config{
		Wallet: WalletConfig{
			KeyDir: filepath.Join(home, ".x402fuel", "keys"),
		},
		Proxy: ProxyConfig{
			ListenAddr: "127.0.0.1:8421",
			MaxPerTxn:  "10.00",
			DailyCap:   "100.00",
			Allowlist:  []string{},
		},
		Server: ServerConfig{
			Port: 8420,
		},
		Chain: ChainConfig{
			RPCURL:  "https://mainnet.base.org",
			ChainID: 8453,
		},
		Events: EventsConfig{
			LogPath: filepath.Join(home, ".x402fuel", "events.jsonl"),
		},
		Telemetry: TelemetryConfig{
			Enabled:  false,
			Endpoint: "https://telemetry.x402fuel.trucore.xyz/v1/ping",
		},
	}
}

// Load reads a YAML config file at path and merges it with defaults.
// Missing values in the YAML fall back to defaults. If path is empty or the
// file does not exist, pure defaults are returned.
func Load(path string) (Config, error) {
	cfg := Default()

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}