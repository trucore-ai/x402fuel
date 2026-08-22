package events

import (
	"os"
	"strings"
	"testing"

	"github.com/trucore-ai/x402fuel/internal/types"
)

func TestEventLogger(t *testing.T) {
	path := t.TempDir() + "/events.jsonl"
	logger, err := NewLogger(path)
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()
	evt := types.Event{
		Timestamp: "2026-08-21T20:00:00Z",
		EventType: types.EventWalletCreated,
		WalletID:  "test-wallet-1",
	}
	if err := logger.Log(evt); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "wallet_created") {
		t.Error("event log file does not contain wallet_created event")
	}
}

func TestEventLoggerKeySafety(t *testing.T) {
	path := t.TempDir() + "/events.jsonl"
	logger, _ := NewLogger(path)
	defer logger.Close()
	evt := types.NewEvent(types.EventPaymentAttempted)
	evt.WalletID = "test-wallet"
	logger.Log(evt)
	data, _ := os.ReadFile(path)
	content := string(data)
	if strings.Contains(content, "private") || strings.Contains(content, "privkey") || strings.Contains(content, "key_path") {
		t.Error("event log contains key-related fields")
	}
}