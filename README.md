# ⚡ x402Fuel

**Non-custodial HTTP 402 wallet daemon for AI agents.**

Give your AI agent a USDC wallet on Base that speaks [x402](https://github.com/x402-foundation/x402). Agents pay for APIs, data, and compute autonomously — with hard budget caps and keys that never leave your machine.

## How it works

```
Agent → API (402 Payment Required) → x402Fuel intercepts
  → Checks budget policy → Signs EIP-3009 → Retries with X-PAYMENT
  → API delivers resource → USDC settles on Base
```

## Quick start

### Install

```bash
go install github.com/trucore-ai/x402fuel@latest
```

### Create a wallet

```bash
x402fuel create --label "My Trading Agent"
# Enter passphrase when prompted
# → 0x... address on Base
```

### Fund it

Send USDC on Base to the wallet address. Base USDC: `0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913`

### Start the daemon

```bash
x402fuel serve
# Dashboard: http://localhost:8420
# API:       http://localhost:8420/api
```

### Point your agent at the proxy

Configure your agent to route HTTP requests through `http://localhost:8421`. When a 402 Payment Required response comes back, x402Fuel handles the payment.

## Configuration

Create `~/.x402fuel/config.yaml`:

```yaml
wallet:
  key_dir: ~/.x402fuel/keys
proxy:
  listen_addr: 127.0.0.1:8421
  max_per_txn: "10.00"    # Max USDC per payment
  daily_cap: "100.00"      # Max USDC per day
  allowlist: []             # Allowed API hosts (empty = all allowed)
chain:
  rpc_url: https://mainnet.base.org
  chain_id: 8453
events:
  log_path: ~/.x402fuel/events.jsonl
telemetry:
  enabled: false            # Opt-in aggregate metrics
```

## CLI

```bash
x402fuel create     # Create a new wallet
x402fuel serve      # Start the daemon
x402fuel pause      # Kill switch — block all payments
x402fuel resume     # Resume payments
x402fuel status     # Wallet summary
```

## API

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/wallets` | Create wallet (body: `{"label":"...", "passphrase":"..."}`) |
| GET | `/api/wallets` | List wallets |
| GET | `/api/wallets/:id` | Wallet details |
| POST | `/api/wallets/:id/pause` | Pause payments |
| POST | `/api/wallets/:id/resume` | Resume payments |
| GET | `/api/wallets/:id/transactions` | Transaction history |
| GET | `/api/policy` | Current spend policy |
| GET | `/health` | Health check |

## Architecture

- **Non-custodial by construction** — keys encrypted at rest (AES-256-GCM), zero key material in logs or telemetry
- **EIP-3009 gasless USDC transfers** — agent wallet doesn't need ETH for gas
- **Budget enforcement** — max-per-txn, daily cap, per-service allowlist, global kill switch
- **Single Go binary** — macOS (arm64) + Linux (amd64/arm64)
- **HTMX dashboard** — embedded, zero JS framework

## Safety

- Keys NEVER leave your machine
- Before any payment is signed: policy check → max per txn → daily cap → allowlist → kill switch
- In-flight signed payments still settle (cannot be revoked on-chain)
- Opt-in telemetry sends aggregate counts only — no URLs, no addresses, no amounts

## License

MIT — free for self-host use. Paid hosted control plane coming later.

## Links

- [x402 Protocol](https://github.com/x402-foundation/x402)
- [Base Network](https://base.org)
- [EIP-3009](https://eips.ethereum.org/EIPS/eip-3009)