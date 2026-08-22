# Security Policy

## Reporting a Vulnerability

**Do not open a public issue.** Send details to:

- Email: security@trucore.xyz (or open a private security advisory on GitHub)

We aim to acknowledge reports within 48 hours and provide an initial assessment within 5 business days.

## Scope

The x402Fuel codebase, including:
- Key storage and encryption
- HTTP 402 payment signing
- Policy enforcement
- Event logging and telemetry
- The local daemon binary

## Out of Scope

- User misconfiguration (weak passphrases, exposed config files)
- Compromised host machines (if your machine is compromised, the attacker has your keys regardless)
- Base network or USDC contract vulnerabilities
- Third-party RPC providers

## Security Model

x402Fuel is **non-custodial by construction**:

1. **Keys never leave the user's machine.** Private keys are encrypted at rest with AES-256-GCM.
2. **No key material in logs or telemetry.** Enforced by automated tests.
3. **Hard budget caps.** Max-per-txn, daily cap, allowlist, and kill switch block payments before signing.
4. **EIP-3009 authorization signatures** — the agent wallet can pay without holding ETH for gas.
5. **Opt-in telemetry** sends aggregate counts only — no URLs, addresses, or amounts.

If you find a way to extract key material from logs, crash dumps, or error messages, that is a critical vulnerability.

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest main branch | ✅ |
| Alpha / pre-release | ✅ (report, but patches may wait for stable) |

## Disclosure Process

1. Report received → acknowledged within 48 hours
2. Triage + reproduction within 5 business days
3. Fix developed + tested
4. Release + advisory published
5. CVE requested if applicable