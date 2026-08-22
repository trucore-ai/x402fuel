package proxy

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/trucore-ai/x402fuel/internal/chain"
	"github.com/trucore-ai/x402fuel/internal/events"
	"github.com/trucore-ai/x402fuel/internal/keystore"
	"github.com/trucore-ai/x402fuel/internal/policy"
	"github.com/trucore-ai/x402fuel/internal/types"
)

// Proxy intercepts outgoing HTTP requests from an agent.
// When a 402 Payment Required response is received, it parses the payment
// requirements, checks policy, signs an EIP-3009 authorization, and returns
// the signed payload so the agent can retry with X-PAYMENT.
type Proxy struct {
	target      *url.URL
	reverse     *httputil.ReverseProxy
	policyEng   *policy.Engine
	keyStore    *keystore.KeyStore
	chainClient *chain.Client
	eventLog    *events.Logger
	walletAddr  string
	passphrase  string
}

func New(
	targetURL string,
	pol *policy.Engine,
	ks *keystore.KeyStore,
	cc *chain.Client,
	el *events.Logger,
	walletAddr, passphrase string,
) (*Proxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("parse target URL: %w", err)
	}
	p := &Proxy{
		target:      target,
		policyEng:   pol,
		keyStore:    ks,
		chainClient: cc,
		eventLog:    el,
		walletAddr:  walletAddr,
		passphrase:  passphrase,
	}
	p.reverse = &httputil.ReverseProxy{
		Director:       p.director,
		ModifyResponse: p.modifyResponse,
	}
	return p, nil
}

func (p *Proxy) director(req *http.Request) {
	req.URL.Scheme = p.target.Scheme
	req.URL.Host = p.target.Host
	req.Host = p.target.Host
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.reverse.ServeHTTP(w, r)
}

func (p *Proxy) modifyResponse(resp *http.Response) error {
	if resp.StatusCode != 402 {
		return nil
	}

	start := time.Now()
	host := resp.Request.URL.Host

	evt := types.NewEvent(types.Event402Encountered)
	evt.Host = host
	evt.WalletID = p.walletAddr
	p.eventLog.Log(evt)

	payReqHeader := resp.Header.Get("X-PAYMENT-REQUIRED")
	if payReqHeader == "" {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		payReqHeader = string(bodyBytes)
	}

	pr, err := parsePaymentRequired(payReqHeader)
	if err != nil {
		log.Printf("proxy: failed to parse 402: %v", err)
		return nil
	}

	var bestReq *types.PaymentRequirements
	for _, req := range pr.Accepts {
		if req.Network == "evm:8453" && req.Asset == "USDC" {
			r := req
			bestReq = &r
			break
		}
	}
	if bestReq == nil {
		log.Printf("proxy: no Base USDC payment option found in %d accepts", len(pr.Accepts))
		return nil
	}

	decision := p.policyEng.Evaluate(bestReq.Amount, host)
	if !decision.Allowed {
		blockedEvt := types.NewEvent(types.EventPaymentBlocked)
		blockedEvt.Host = host
		blockedEvt.Amount = bestReq.Amount
		blockedEvt.Decision = "blocked_by_policy"
		blockedEvt.Reason = decision.Reason
		blockedEvt.WalletID = p.walletAddr
		p.eventLog.Log(blockedEvt)

		pr.Error = fmt.Sprintf("payment blocked: %s", decision.Reason)
		return p.write402Response(resp, pr)
	}

	priv, err := p.keyStore.Load(p.walletAddr, p.passphrase)
	if err != nil {
		log.Printf("proxy: load key: %v", err)
		return nil
	}
	defer priv.D.SetInt64(0)

	amount, _ := chain.FormatUSDC(bestReq.Amount)
	from := common.HexToAddress(p.walletAddr)
	to := common.HexToAddress(bestReq.PayTo)
	validAfter := big.NewInt(time.Now().Unix())
	validBefore := big.NewInt(time.Now().Add(30 * time.Minute).Unix())
	var nonce [32]byte

	sig, err := chain.SignEIP3009(priv, from, to, amount, validAfter, validBefore, nonce)
	if err != nil {
		log.Printf("proxy: sign EIP-3009: %v", err)
		return nil
	}

	payload := types.PaymentPayload{
		X402Version: 1,
		Resource:    pr.Resource,
		Accepted:    *bestReq,
		Payload: map[string]interface{}{
			"signature":     fmt.Sprintf("0x%x", sig),
			"from":          p.walletAddr,
			"validAfter":    validAfter.String(),
			"validBefore":   validBefore.String(),
			"nonce":         fmt.Sprintf("0x%x", nonce),
			"authorization": "eip3009",
		},
	}

	payloadJSON, _ := json.Marshal(payload)
	payloadB64 := base64.StdEncoding.EncodeToString(payloadJSON)

	attemptEvt := types.NewEvent(types.EventPaymentAttempted)
	attemptEvt.Host = host
	attemptEvt.Amount = bestReq.Amount
	attemptEvt.WalletID = p.walletAddr
	attemptEvt.LatencyMs = time.Since(start).Milliseconds()
	p.eventLog.Log(attemptEvt)

	p.policyEng.RecordSpend(bestReq.Amount)

	modifiedBody := map[string]interface{}{
		"payment_payload": payload,
		"payment_header":  fmt.Sprintf("X-PAYMENT: %s", payloadB64),
	}
	bodyJSON, _ := json.Marshal(modifiedBody)
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyJSON))
	resp.ContentLength = int64(len(bodyJSON))
	resp.Header.Set("Content-Type", "application/json")
	resp.Header.Set("X-402FUEL-SIGNED", "true")
	return nil
}

func (p *Proxy) write402Response(resp *http.Response, pr *types.PaymentRequired) error {
	body, _ := json.Marshal(pr)
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	resp.ContentLength = int64(len(body))
	resp.Header.Set("Content-Type", "application/json")
	resp.StatusCode = 402
	return nil
}

func parsePaymentRequired(input string) (*types.PaymentRequired, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		decoded = []byte(input)
	}
	var pr types.PaymentRequired
	if err := json.Unmarshal(decoded, &pr); err != nil {
		return nil, fmt.Errorf("unmarshal PaymentRequired: %w", err)
	}
	return &pr, nil
}