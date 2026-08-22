package policy

import (
	"fmt"
	"math/big"
	"time"

	"github.com/trucore-ai/x402fuel/internal/chain"
	"github.com/trucore-ai/x402fuel/internal/types"
)

type Engine struct {
	policy *types.SpendPolicy
}

func NewEngine(policy *types.SpendPolicy) *Engine {
	if policy.DayStart.IsZero() {
		policy.DayStart = time.Now().UTC().Truncate(24 * time.Hour)
	}
	if policy.DailySpent == "" {
		policy.DailySpent = "0.000000"
	}
	return &Engine{policy: policy}
}

func (e *Engine) Evaluate(amount string, host string) types.PolicyDecision {
	if e.policy.Paused {
		return types.PolicyDecision{Allowed: false, Reason: "kill_switch_active"}
	}
	if len(e.policy.Allowlist) > 0 {
		allowed := false
		for _, a := range e.policy.Allowlist {
			if a == host || a == "*" {
				allowed = true
				break
			}
		}
		if !allowed {
			return types.PolicyDecision{Allowed: false, Reason: fmt.Sprintf("host_not_allowlisted: %s", host)}
		}
	}
	maxPerTxn, err := chain.FormatUSDC(e.policy.MaxPerTxn)
	if err != nil {
		return types.PolicyDecision{Allowed: false, Reason: "invalid_max_per_txn_config"}
	}
	requestedAmt, err := chain.FormatUSDC(amount)
	if err != nil {
		return types.PolicyDecision{Allowed: false, Reason: fmt.Sprintf("invalid_amount: %s", amount)}
	}
	if requestedAmt.Cmp(maxPerTxn) > 0 {
		return types.PolicyDecision{Allowed: false, Reason: fmt.Sprintf("exceeds_max_per_txn (%s > %s)", amount, e.policy.MaxPerTxn)}
	}
	e.resetDailyIfNeeded()
	dailyCap, err := chain.FormatUSDC(e.policy.DailyCap)
	if err != nil {
		return types.PolicyDecision{Allowed: false, Reason: "invalid_daily_cap_config"}
	}
	dailySpent, err := chain.FormatUSDC(e.policy.DailySpent)
	if err != nil {
		dailySpent = big.NewInt(0)
	}
	newTotal := new(big.Int).Add(dailySpent, requestedAmt)
	if newTotal.Cmp(dailyCap) > 0 {
		remaining := new(big.Int).Sub(dailyCap, dailySpent)
		return types.PolicyDecision{
			Allowed: false,
			Reason:  fmt.Sprintf("daily_cap_exceeded: spent %s, cap %s, remaining %s", e.policy.DailySpent, e.policy.DailyCap, chain.FormatUSDCFromInt(remaining)),
		}
	}
	return types.PolicyDecision{Allowed: true}
}

func (e *Engine) RecordSpend(amount string) error {
	e.resetDailyIfNeeded()
	spent, err := chain.FormatUSDC(e.policy.DailySpent)
	if err != nil {
		spent = big.NewInt(0)
	}
	amt, err := chain.FormatUSDC(amount)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}
	e.policy.DailySpent = chain.FormatUSDCFromInt(new(big.Int).Add(spent, amt))
	return nil
}

func (e *Engine) resetDailyIfNeeded() {
	now := time.Now().UTC()
	if e.policy.DayStart.IsZero() || now.Sub(e.policy.DayStart) >= 24*time.Hour {
		e.policy.DayStart = now.Truncate(24 * time.Hour)
		e.policy.DailySpent = "0.000000"
	}
}

func (e *Engine) SetPaused(paused bool) { e.policy.Paused = paused }
func (e *Engine) IsPaused() bool         { return e.policy.Paused }
func (e *Engine) GetPolicy() *types.SpendPolicy { return e.policy }