package types

import "time"

type SpendPolicy struct {
	MaxPerTxn   string    `json:"max_per_txn"`
	DailyCap    string    `json:"daily_cap"`
	Allowlist   []string  `json:"allowlist"`
	DailySpent  string    `json:"daily_spent"`
	DayStart    time.Time `json:"day_start"`
	Paused      bool      `json:"paused"`
}

type PolicyDecision struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

func (p *SpendPolicy) CheckPaused() PolicyDecision {
	if p.Paused {
		return PolicyDecision{Allowed: false, Reason: "wallet_paused"}
	}
	return PolicyDecision{Allowed: true}
}