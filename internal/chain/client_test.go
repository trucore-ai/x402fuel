package chain

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestFormatUSDC(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.50", "1500000"},
		{"0.01", "10000"},
		{"100", "100000000"},
		{"0.000001", "1"},
	}
	for _, tt := range tests {
		result, err := FormatUSDC(tt.input)
		if err != nil {
			t.Errorf("FormatUSDC(%q): unexpected error: %v", tt.input, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("FormatUSDC(%q): expected %s, got %s", tt.input, tt.expected, result.String())
		}
	}
}

func TestFormatUSDCFromInt(t *testing.T) {
	result := FormatUSDCFromInt(big.NewInt(1500000))
	if result != "1.500000" {
		t.Errorf("expected 1.500000, got %s", result)
	}
}

func TestEIP3009Signing(t *testing.T) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	from := crypto.PubkeyToAddress(priv.PublicKey)
	to := common.HexToAddress("0x1111111111111111111111111111111111111111")
	value := big.NewInt(1000000)
	validAfter := big.NewInt(0)
	validBefore := big.NewInt(9999999999)
	var nonce [32]byte
	sig, err := SignEIP3009(priv, from, to, value, validAfter, validBefore, nonce)
	if err != nil {
		t.Fatal(err)
	}
	if len(sig) != 65 {
		t.Errorf("expected 65-byte signature, got %d", len(sig))
	}
}