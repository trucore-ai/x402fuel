package keystore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/pbkdf2"
)

const (
	keyDirPerm  = 0700
	keyFilePerm = 0600
	pbkdf2Iter  = 600000
	keyLen      = 32
	saltLen     = 16
)

type encryptedKey struct {
	Address string `json:"address"`
	Cipher  string `json:"cipher"`
	Nonce   string `json:"nonce"`
	Salt    string `json:"salt"`
}

type KeyStore struct {
	dir string
}

func New(dir string) (*KeyStore, error) {
	if err := os.MkdirAll(dir, keyDirPerm); err != nil {
		return nil, fmt.Errorf("create key dir: %w", err)
	}
	return &KeyStore{dir: dir}, nil
}

func (ks *KeyStore) Generate(passphrase string) (address string, err error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("generate key: %w", err)
	}
	address = crypto.PubkeyToAddress(priv.PublicKey).Hex()
	if err := ks.encryptAndSave(address, priv, passphrase); err != nil {
		return "", err
	}
	priv.D.SetInt64(0)
	return address, nil
}

func (ks *KeyStore) Load(address, passphrase string) (*ecdsa.PrivateKey, error) {
	path := ks.keyPath(address)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read key file: %w", err)
	}
	var ek encryptedKey
	if err := json.Unmarshal(data, &ek); err != nil {
		return nil, fmt.Errorf("parse key file: %w", err)
	}
	return ks.decrypt(ek, passphrase)
}

func (ks *KeyStore) List() ([]string, error) {
	entries, err := os.ReadDir(ks.dir)
	if err != nil {
		return nil, fmt.Errorf("read key dir: %w", err)
	}
	var addrs []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".key" {
			addrs = append(addrs, "0x"+e.Name()[:len(e.Name())-4])
		}
	}
	return addrs, nil
}

func (ks *KeyStore) encryptAndSave(address string, priv *ecdsa.PrivateKey, passphrase string) error {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("generate salt: %w", err)
	}
	key := pbkdf2.Key([]byte(passphrase), salt, pbkdf2Iter, keyLen, sha256.New)
	privBytes := crypto.FromECDSA(priv)
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nil, nonce, privBytes, nil)
	ek := encryptedKey{
		Address: address,
		Cipher:  hex.EncodeToString(ciphertext),
		Nonce:   hex.EncodeToString(nonce),
		Salt:    hex.EncodeToString(salt),
	}
	data, err := json.Marshal(ek)
	if err != nil {
		return fmt.Errorf("marshal key: %w", err)
	}
	return os.WriteFile(ks.keyPath(address), data, keyFilePerm)
}

func (ks *KeyStore) decrypt(ek encryptedKey, passphrase string) (*ecdsa.PrivateKey, error) {
	salt, _ := hex.DecodeString(ek.Salt)
	nonce, _ := hex.DecodeString(ek.Nonce)
	ciphertext, _ := hex.DecodeString(ek.Cipher)
	key := pbkdf2.Key([]byte(passphrase), salt, pbkdf2Iter, keyLen, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: wrong passphrase or corrupted key")
	}
	priv, err := crypto.ToECDSA(plaintext)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}
	return priv, nil
}

func (ks *KeyStore) keyPath(address string) string {
	addr := address
	if len(addr) > 2 && addr[:2] == "0x" {
		addr = addr[2:]
	}
	return filepath.Join(kes.dir, addr+".key")
}
