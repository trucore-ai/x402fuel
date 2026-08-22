package types

type PaymentRequired struct {
	X402Version int                    `json:"x402Version"`
	Error       string                 `json:"error,omitempty"`
	Resource    ResourceInfo           `json:"resource"`
	Accepts     []PaymentRequirements  `json:"accepts"`
	Extensions  map[string]interface{} `json:"extensions,omitempty"`
}

type ResourceInfo struct {
	URL         string   `json:"url"`
	Description string   `json:"description,omitempty"`
	MimeType    string   `json:"mimeType,omitempty"`
	ServiceName string   `json:"serviceName,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	IconURL     string   `json:"iconUrl,omitempty"`
}

type PaymentRequirements struct {
	Scheme            string                 `json:"scheme"`
	Network           string                 `json:"network"`
	Asset             string                 `json:"asset"`
	Amount            string                 `json:"amount"`
	PayTo             string                 `json:"payTo"`
	MaxTimeoutSeconds int                    `json:"maxTimeoutSeconds"`
	Extra             map[string]interface{} `json:"extra"`
}

type PaymentPayload struct {
	X402Version int                    `json:"x402Version"`
	Resource    ResourceInfo           `json:"resource,omitempty"`
	Accepted    PaymentRequirements    `json:"accepted"`
	Payload     map[string]interface{} `json:"payload"`
	Extensions  map[string]interface{} `json:"extensions,omitempty"`
}

type Transaction struct {
	ID        string `json:"id"`
	WalletID  string `json:"wallet_id"`
	Host      string `json:"host"`
	Amount    string `json:"amount"`
	Asset     string `json:"asset"`
	TxHash    string `json:"tx_hash,omitempty"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	SettledAt string `json:"settled_at,omitempty"`
}