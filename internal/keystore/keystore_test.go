package keystore

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestGenerateAndLoad(t *testing.T) {
	dir := t.TempDir()
	ks, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	addr, err := ks.Generate("test-passphrase")
	if err != nil {
		t.Fatal(err)
	}
	if len(addr) < 42 {
		t.Errorf("invalid address: %s", addr)
	}
	priv, err := ks.Load(addr, "test-passphrase")
	if err != nil {
		t.Fatal(err)
	}
	if crypto.PubkeyToAddress(priv.PublicKey).Hex() != addr {
		t.Error("decrypted key does not match original address")
	}
}

func TestLoadWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	ks, _ := New(dir)
	addr, _ := ks.Generate("correct")
	_, err := ks.Load(addr, "wrong")
	if err == nil {
		t.Error("expected error with wrong passphrase")
	}
}

func TestList(t *testing.T) {
	dir := t.TempDir()
	ks, _ := New(dir)
	ks.Generate("pw1")
	ks.Generate("pw2")
	addrs, err := ks.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Errorf("expected 2 wallets, got %d", len(addrs))
	}
}

func TestTransientKeyUsage(t *testing.T) {
	dir := t.TempDir()
	ks, _ := New(dir)
	addr, _ := ks.Generate("pw")
	priv, _ := ks.Load(addr, "pw")
	_ = crypto.PubkeyToAddress(priv.PublicKey)
	priv.D.SetInt64(0)
	if priv.D.Sign() != 0 {
		t.Error("key D value should be zero after zeroing")
	}
}