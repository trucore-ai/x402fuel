package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	home, _ := os.UserHomeDir()

	if cfg.Wallet.KeyDir != filepath.Join(home, ".x402fuel", "keys") {
		t.Errorf("expected wallet key dir %s, got %s", filepath.Join(home, ".x402fuel", "keys"), cfg.Wallet.KeyDir)
	}
	if cfg.Proxy.ListenAddr != "127.0.0.1:8421" {
		t.Errorf("expected proxy listen addr 127.0.0.1:8421, got %s", cfg.Proxy.ListenAddr)
	}
	if cfg.Proxy.MaxPerTxn != "10.00" {
		t.Errorf("expected max per txn 10.00, got %s", cfg.Proxy.MaxPerTxn)
	}
	if cfg.Proxy.DailyCap != "100.00" {
		t.Errorf("expected daily cap 100.00, got %s", cfg.Proxy.DailyCap)
	}
	if cfg.Server.Port != 8420 {
		t.Errorf("expected server port 8420, got %d", cfg.Server.Port)
	}
	if cfg.Chain.RPCURL != "https://mainnet.base.org" {
		t.Errorf("expected chain rpc url https://mainnet.base.org, got %s", cfg.Chain.RPCURL)
	}
	if cfg.Chain.ChainID != 8453 {
		t.Errorf("expected chain id 8453, got %d", cfg.Chain.ChainID)
	}
	if cfg.Events.LogPath != filepath.Join(home, ".x402fuel", "events.jsonl") {
		t.Errorf("expected events log path %s, got %s", filepath.Join(home, ".x402fuel", "events.jsonl"), cfg.Events.LogPath)
	}
	if cfg.Telemetry.Enabled != false {
		t.Errorf("expected telemetry disabled, got %v", cfg.Telemetry.Enabled)
	}
	if cfg.Telemetry.Endpoint != "https://telemetry.x402fuel.trucore.xyz/v1/ping" {
		t.Errorf("unexpected telemetry endpoint: %s", cfg.Telemetry.Endpoint)
	}
}

func TestLoad_DefaultsWhenEmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load(\"\") returned error: %v", err)
	}

	if cfg.Server.Port != 8420 {
		t.Errorf("expected default port 8420, got %d", cfg.Server.Port)
	}
}

func TestLoad_DefaultsWhenFileMissing(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("Load(missing file) returned error: %v", err)
	}

	if cfg.Server.Port != 8420 {
		t.Errorf("expected default port 8420, got %d", cfg.Server.Port)
	}
}

func TestLoad_YAMLFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `
wallet:
  key_dir: /custom/keys
server:
  port: 9000
proxy:
  listen_addr: 0.0.0.0:9001
  max_per_txn: "50.00"
chain:
  chain_id: 84532
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load(%s) returned error: %v", configPath, err)
	}

	if cfg.Wallet.KeyDir != "/custom/keys" {
		t.Errorf("expected wallet key dir /custom/keys, got %s", cfg.Wallet.KeyDir)
	}
	if cfg.Server.Port != 9000 {
		t.Errorf("expected server port 9000, got %d", cfg.Server.Port)
	}
	if cfg.Proxy.ListenAddr != "0.0.0.0:9001" {
		t.Errorf("expected proxy listen addr 0.0.0.0:9001, got %s", cfg.Proxy.ListenAddr)
	}
	if cfg.Proxy.MaxPerTxn != "50.00" {
		t.Errorf("expected max per txn 50.00, got %s", cfg.Proxy.MaxPerTxn)
	}
	if cfg.Chain.ChainID != 84532 {
		t.Errorf("expected chain id 84532, got %d", cfg.Chain.ChainID)
	}

	// Unspecified fields should remain at defaults
	if cfg.Proxy.DailyCap != "100.00" {
		t.Errorf("expected daily cap default 100.00, got %s", cfg.Proxy.DailyCap)
	}
	if cfg.Chain.RPCURL != "https://mainnet.base.org" {
		t.Errorf("expected chain rpc url default, got %s", cfg.Chain.RPCURL)
	}
}