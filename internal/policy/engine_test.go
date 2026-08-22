package policy

import (
	"testing"

	"github.com/trucore-ai/x402fuel/internal/types"
)

func TestAllowPayment(t *testing.T) {
	policy := &types.SpendPolicy{MaxPerTxn: "10.00", DailyCap: "100.00", DailySpent: "0.000000"}
	eng := NewEngine(policy)
	result := eng.Evaluate("5.00", "api.example.com")
	if !result.Allowed {
		t.Errorf("expected allowed, got blocked: %s", result.Reason)
	}
}

func TestBlockExceedsMaxPerTxn(t *testing.T) {
	policy := &types.SpendPolicy{MaxPerTxn: "10.00", DailyCap: "100.00", DailySpent: "0.000000"}
	eng := NewEngine(policy)
	result := eng.Evaluate("15.00", "api.example.com")
	if result.Allowed {
		t.Error("expected blocked for exceeding max_per_txn")
	}
}

func TestBlockAllowlist(t *testing.T) {
	policy := &types.SpendPolicy{MaxPerTxn: "10.00", DailyCap: "100.00", Allowlist: []string{"allowed-api.example.com"}}
	eng := NewEngine(policy)
	if eng.Evaluate("1.00", "evil-api.example.com").Allowed {
		t.Error("expected blocked for non-allowlisted host")
	}
	if !eng.Evaluate("1.00", "allowed-api.example.com").Allowed {
		t.Error("expected allowed for allowlisted host")
	}
}

func TestBlockPaused(t *testing.T) {
	policy := &types.SpendPolicy{MaxPerTxn: "10.00", DailyCap: "100.00", Paused: true}
	eng := NewEngine(policy)
	if eng.Evaluate("1.00", "api.example.com").Allowed {
		t.Error("expected blocked when paused")
	}
}

func TestDailyCapEnforcement(t *testing.T) {
	policy := &types.SpendPolicy{MaxPerTxn: "50.00", DailyCap: "100.00", DailySpent: "95.000000"}
	eng := NewEngine(policy)
	result := eng.Evaluate("5.00", "api.example.com")
	if !result.Allowed {
		t.Errorf("expected allowed: %s", result.Reason)
	}
	eng.RecordSpend("5.00")
	result2 := eng.Evaluate("1.00", "api.example.com")
	if result2.Allowed {
		t.Error("expected blocked for exceeding daily cap")
	}
}