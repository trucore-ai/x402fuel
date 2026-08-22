package chain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const USDCAddress = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"

var erc20ABI = `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"}]`

type Client struct {
	rpc    *ethclient.Client
	usdc   common.Address
	parsed abi.ABI
}

func NewClient(rpcURL string) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial RPC: %w", err)
	}
	parsed, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return nil, fmt.Errorf("parse ABI: %w", err)
	}
	return &Client{rpc: client, usdc: common.HexToAddress(USDCAddress), parsed: parsed}, nil
}

func (c *Client) Balance(ctx context.Context, address common.Address) (*big.Int, error) {
	data, err := c.parsed.Pack("balanceOf", address)
	if err != nil {
		return nil, fmt.Errorf("pack balanceOf: %w", err)
	}
	result, err := c.rpc.CallContract(ctx, ethereum.CallMsg{To: &c.usdc, Data: data}, nil)
	if err != nil {
		return nil, fmt.Errorf("call balanceOf: %w", err)
	}
	balance := new(big.Int)
	balance.SetBytes(result)
	return balance, nil
}

// SignEIP3009 creates an EIP-3009 authorization for gasless USDC transfer.
func SignEIP3009(priv *ecdsa.PrivateKey, from, to common.Address, value *big.Int, validAfter, validBefore *big.Int, nonce [32]byte) ([]byte, error) {
	authorizationTypeHash := crypto.Keccak256Hash([]byte(
		"TransferWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)",
	))
	domainSeparator := computeDomainSeparator(8453, USDCAddress)
	encoded := crypto.Keccak256(
		padBytes32(authorizationTypeHash.Bytes()),
		padAddress(from),
		padAddress(to),
		padUint256(value),
		padUint256(validAfter),
		padUint256(validBefore),
		nonce[:],
	)
	dataHash := crypto.Keccak256([]byte("\x19\x01"), domainSeparator[:], encoded)
	sig, err := crypto.Sign(dataHash, priv)
	if err != nil {
		return nil, fmt.Errorf("sign EIP-3009: %w", err)
	}
	sig[64] += 27
	return sig, nil
}

// FormatUSDC converts a string amount like "1.50" to USDC's 6-decimal big.Int.
// Uses string-based parsing to avoid floating-point precision loss.
func FormatUSDC(amount string) (*big.Int, error) {
	parts := strings.SplitN(amount, ".", 2)
	whole := new(big.Int)
	if _, ok := whole.SetString(parts[0], 10); !ok {
		return nil, fmt.Errorf("parse USDC amount %q: invalid whole part", amount)
	}
	whole.Mul(whole, big.NewInt(1_000_000))

	if len(parts) == 1 {
		return whole, nil
	}

	fracStr := parts[1]
	if len(fracStr) > 6 {
		fracStr = fracStr[:6]
	}
	for len(fracStr) < 6 {
		fracStr += "0"
	}
	frac := new(big.Int)
	if _, ok := frac.SetString(fracStr, 10); !ok {
		return nil, fmt.Errorf("parse USDC amount %q: invalid fractional part", amount)
	}
	return whole.Add(whole, frac), nil
}

// FormatUSDCFromInt converts a big.Int USDC amount (6 decimals) to human-readable string.
func FormatUSDCFromInt(amount *big.Int) string {
	whole := new(big.Int).Div(amount, big.NewInt(1_000_000))
	frac := new(big.Int).Mod(amount, big.NewInt(1_000_000))
	return fmt.Sprintf("%s.%06s", whole.String(), frac.String())
}

func padAddress(addr common.Address) []byte {
	return common.LeftPadBytes(addr.Bytes(), 32)
}

func padUint256(val *big.Int) []byte {
	return common.LeftPadBytes(val.Bytes(), 32)
}

func padBytes32(b []byte) []byte {
	return common.LeftPadBytes(b, 32)
}

func computeDomainSeparator(chainID int64, verifyingContract string) common.Hash {
	domainTypeHash := crypto.Keccak256Hash([]byte(
		"EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)",
	))
	nameHash := crypto.Keccak256Hash([]byte("USD Coin"))
	versionHash := crypto.Keccak256Hash([]byte("2"))
	result := crypto.Keccak256(
		domainTypeHash.Bytes(),
		padBytes32(nameHash.Bytes()),
		padBytes32(versionHash.Bytes()),
		padUint256(big.NewInt(chainID)),
		padAddress(common.HexToAddress(verifyingContract)),
	)
	return common.BytesToHash(result)
}