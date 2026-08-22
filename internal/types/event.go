package types

import "time"

type EventType string

const (
	EventWalletCreated    EventType = "wallet_created"
	Event402Encountered   EventType = "402_encountered"
	EventPaymentBlocked   EventType = "payment_blocked"
	EventPaymentAttempted EventType = "payment_attempted"
	EventPaymentSettled   EventType = "payment_settled"
	EventPaymentFailed    EventType = "payment_failed"
)

type Event struct {
	Timestamp string            `json:"timestamp"`
	EventType EventType         `json:"event_type"`
	WalletID  string            `json:"wallet_id,omitempty"`
	Host      string            `json:"host,omitempty"`
	Amount    string            `json:"amount,omitempty"`
	Asset     string            `json:"asset,omitempty"`
	TxHash    string            `json:"tx_hash,omitempty"`
	Decision  string            `json:"decision,omitempty"`
	Reason    string            `json:"reason,omitempty"`
	LatencyMs int64             `json:"latency_ms,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

func NewEvent(t EventType) Event {
	return Event{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		EventType: t,
	}
}