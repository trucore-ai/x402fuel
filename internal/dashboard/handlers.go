package dashboard

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/trucore-ai/x402fuel/internal/keystore"
	"github.com/trucore-ai/x402fuel/internal/policy"
)

func Register(r *http.ServeMux, ks *keystore.KeyStore, pol *policy.Engine) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(indexHTML, truCoreHeader+fuelChrome)))
	})

	r.HandleFunc("/wallets", func(w http.ResponseWriter, r *http.Request) {
		addrs, _ := ks.List()
		var rows strings.Builder
		for _, addr := range addrs {
			rows.WriteString(fmt.Sprintf(`<tr><td class="font-mono text-sm">%s</td></tr>`, addr))
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(walletsHTML, truCoreHeader+fuelChrome, rows.String())))
	})

	r.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		p := pol.GetPolicy()
		paused := "No"
		if p.Paused {
			paused = "Yes \u23f8"
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(settingsHTML, truCoreHeader+fuelChrome, p.MaxPerTxn, p.DailyCap, p.DailySpent, paused)))
	})
}

// truCoreHeader is the shared TruCore site header rendered above the app nav
// on every dashboard page. Keep in sync with the header on www.trucore.xyz.
const truCoreHeader = `
<style>
  .tc-header { position: sticky; top: 0; z-index: 50; background: linear-gradient(to bottom, rgba(5,10,20,0.25) 0%, rgba(5,10,20,0.25) 55%, rgba(5,10,20,0) 100%); }
  .tc-header-content { transition: opacity 320ms ease; }
  .tc-header.is-scrolled .tc-header-content { opacity: 0.75; }
  .tc-header-inner { max-width: 1180px; margin: 0 auto; display: flex; flex-wrap: wrap; align-items: center; gap: 6px 24px; padding: 10px 20px; min-height: 56px; }
  .tc-logo { display: inline-flex; align-items: center; gap: 10px; text-decoration: none; border-radius: 6px; }
  .tc-logo img { width: 36px; height: 36px; border-radius: 6px; object-fit: contain; }
  .tc-logo span { font-size: 20px; font-weight: 700; letter-spacing: -0.02em; color: #fff; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; }
  .tc-nav { display: flex; flex-wrap: wrap; align-items: center; gap: 2px; margin-left: auto; }
  .tc-nav a { font-size: 14px; font-weight: 500; color: #cbd5e1; text-decoration: none; padding: 6px 10px; border-radius: 6px; transition: background-color 200ms ease, color 200ms ease; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; }
  .tc-nav a:hover { background: rgba(255,255,255,0.05); color: #f8fafc; }
  .tc-nav a[aria-current="true"] { color: #f8fafc; }
  @media (min-width: 640px) {
    .tc-header-inner { min-height: 68px; }
    .tc-logo img { width: 40px; height: 40px; }
  }
</style>
<header class="tc-header">
  <div class="tc-header-content">
    <div class="tc-header-inner">
      <a class="tc-logo" href="https://www.trucore.xyz" aria-label="TruCore home">
        <img src="https://www.trucore.xyz/media/trucore-mark-128.png" alt="TruCore logo" width="40" height="40" />
        <span>TruCore</span>
      </a>
      <nav class="tc-nav" aria-label="Primary">
        <a href="https://www.trucore.xyz/atf">ATF</a>
        <a href="https://meshdns.trucore.xyz">MeshDNS</a>
        <a href="https://x402fuel.trucore.xyz" aria-current="true">x402Fuel</a>
        <a href="https://www.trucore.xyz/builders">Builders</a>
      </nav>
    </div>
  </div>
</header>
<script>
  (function () {
    var h = document.querySelector('.tc-header');
    if (!h) return;
    function u() { h.classList.toggle('is-scrolled', window.scrollY > 12); }
    window.addEventListener('scroll', u, { passive: true });
    u();
  })();
</script>`

// fuelChrome restyles the dashboard with TruCore brand tokens while keeping
// x402Fuel's own identity: sky-blue brand shell, amber fuel accents.
const fuelChrome = `
<style>
  :root {
    --tc-bg: #050a14;
    --tc-bg2: #0b1220;
    --tc-panel: rgba(11, 18, 32, 0.92);
    --tc-border: rgba(255, 255, 255, 0.08);
    --tc-text: #e5e7eb;
    --tc-muted: #9aa3b2;
    --fuel-sky: #38bdf8;
    --fuel-amber: #fbbf24;
    --fuel-orange: #f08a1f;
    --sans: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    --mono: ui-monospace, "SFMono-Regular", Menlo, Monaco, Consolas, "Liberation Mono", monospace;
  }
  body {
    font-family: var(--sans);
    color: var(--tc-text);
    background:
      radial-gradient(circle at 15% 0%, rgba(56, 189, 248, 0.12), transparent 32%),
      radial-gradient(circle at 85% 8%, rgba(251, 191, 36, 0.07), transparent 30%),
      linear-gradient(180deg, var(--tc-bg), var(--tc-bg2)) !important;
  }
  nav.bg-gray-800 { background: rgba(11, 18, 32, 0.6) !important; border-color: var(--tc-border) !important; }
  nav.bg-gray-800 a { color: #cbd5e1; transition: color 200ms ease; }
  nav.bg-gray-800 a:hover { color: var(--fuel-sky) !important; }
  .bg-gray-800 { background-color: var(--tc-panel) !important; border: 1px solid var(--tc-border); border-radius: 14px; box-shadow: 0 12px 32px rgba(0, 0, 0, 0.35); }
  .text-gray-400 { color: var(--tc-muted) !important; }
  .text-gray-500 { color: var(--tc-muted) !important; }
  input.bg-gray-700 { background-color: rgba(255, 255, 255, 0.04) !important; border: 1px solid var(--tc-border); border-radius: 10px; color: var(--tc-text); font-family: var(--mono); }
  button.bg-yellow-500 { background: linear-gradient(180deg, var(--fuel-amber), var(--fuel-orange)) !important; color: #050a14 !important; border-radius: 10px; }
  button.bg-yellow-500:hover { filter: brightness(1.1); }
  h1, h2 { letter-spacing: -0.03em; }
  .fuel-hero { display: flex; flex-wrap: wrap; align-items: center; justify-content: space-between; gap: 24px; margin-bottom: 28px; }
  .fuel-eyebrow { color: var(--fuel-sky); text-transform: uppercase; letter-spacing: 0.18em; font-size: 12px; margin-bottom: 10px; }
  .fuel-h1 { font-size: clamp(26px, 3.4vw, 40px); font-weight: 700; letter-spacing: -0.04em; line-height: 1.05; margin: 0 0 10px; max-width: 22ch; }
  .fuel-lede { color: var(--tc-muted); max-width: 52ch; margin: 0; font-size: 15px; }
  .fuel-mark { width: 170px; height: 170px; flex-shrink: 0; filter: drop-shadow(0 0 28px rgba(56, 189, 248, 0.25)); }
  .fuel-charge { animation: fuelCharge 2.8s ease-in-out infinite; }
  .fuel-charge-2 { animation-delay: 0.9s; }
  .fuel-charge-3 { animation-delay: 1.8s; }
  @keyframes fuelCharge { 0%, 100% { opacity: 0.25; } 50% { opacity: 1; } }
  .fuel-mini { color: var(--fuel-sky); margin-right: 8px; }
  .section-h { font-size: 1.25rem; font-weight: 700; margin-bottom: 1rem; }
  @media (prefers-reduced-motion: reduce) { .fuel-charge { animation: none; } }
  @media (max-width: 640px) { .fuel-mark { width: 120px; height: 120px; } }
</style>`

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>x402Fuel</title>
<script src="https://unpkg.com/htmx.org@2.0.4"></script>
<link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-gray-900 text-gray-100 min-h-screen">
%s
<nav class="bg-gray-800 border-b border-gray-700 px-6 py-3 flex gap-6">
  <span class="font-bold text-yellow-400">⚡ x402Fuel</span>
  <a href="/" class="hover:text-yellow-300" hx-get="/" hx-target="#main" hx-push-url="true">Dashboard</a>
  <a href="/wallets" class="hover:text-yellow-300" hx-get="/wallets" hx-target="#main" hx-push-url="true">Wallets</a>
  <a href="/settings" class="hover:text-yellow-300" hx-get="/settings" hx-target="#main" hx-push-url="true">Settings</a>
</nav>
<main id="main" class="p-6 max-w-4xl mx-auto">
  <section class="fuel-hero">
    <div>
      <div class="fuel-eyebrow"><span aria-hidden="true">◇</span> x402Fuel · USDC wallets for agents</div>
      <h1 class="fuel-h1">Fuel for agents that pay their own way.</h1>
      <p class="fuel-lede">Non-custodial USDC wallets, spend policies, and HTTP 402 settlement on Base. One daemon, self-hosted, MIT-licensed.</p>
    </div>
    <svg class="fuel-mark" viewBox="0 0 200 200" role="img" aria-label="x402Fuel fuel cell mark">
      <defs>
        <linearGradient id="fuel-hex" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0" stop-color="#38bdf8" />
          <stop offset="1" stop-color="#2582cb" />
        </linearGradient>
        <linearGradient id="fuel-bolt" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0" stop-color="#fbbf24" />
          <stop offset="1" stop-color="#f08a1f" />
        </linearGradient>
      </defs>
      <polygon points="100,14 174,57 174,143 100,186 26,143 26,57" fill="rgba(56,189,248,0.05)" stroke="url(#fuel-hex)" stroke-width="2" />
      <polygon points="100,34 157,67 157,133 100,166 43,133 43,67" fill="none" stroke="rgba(255,255,255,0.08)" stroke-width="1" />
      <path d="M112 44 L68 110 L98 110 L88 156 L134 86 L103 86 Z" fill="url(#fuel-bolt)" />
      <g stroke="#38bdf8" stroke-width="2" stroke-linecap="round">
        <line class="fuel-charge" x1="100" y1="2" x2="100" y2="10" />
        <line class="fuel-charge fuel-charge-2" x1="184" y1="52" x2="177" y2="56" />
        <line class="fuel-charge fuel-charge-3" x1="16" y1="52" x2="23" y2="56" />
      </g>
    </svg>
  </section>
  <h2 class="section-h">Dashboard</h2>
  <div class="grid grid-cols-2 gap-4">
    <div class="bg-gray-800 rounded-lg p-4"><div class="text-gray-400 text-sm">Wallet Status</div><div class="text-2xl font-bold text-green-400">Active</div></div>
    <div class="bg-gray-800 rounded-lg p-4"><div class="text-gray-400 text-sm">Spending Today</div><div class="text-2xl font-bold">$0.00</div></div>
    <div class="bg-gray-800 rounded-lg p-4"><div class="text-gray-400 text-sm">Balance</div><div class="text-2xl font-bold">—</div></div>
    <div class="bg-gray-800 rounded-lg p-4"><div class="text-gray-400 text-sm">Network</div><div class="text-2xl font-bold">Base</div></div>
  </div>
  <div class="mt-6 bg-gray-800 rounded-lg p-4">
    <h2 class="text-lg font-semibold mb-2">Recent Activity</h2>
    <p class="text-gray-500">No transactions yet. Fund your wallet and point your agent at the proxy to get started.</p>
  </div>
</main>
</body></html>`

const walletsHTML = `<!DOCTYPE html><html lang="en">
<head><meta charset="UTF-8"><title>x402Fuel — Wallets</title>
<script src="https://unpkg.com/htmx.org@2.0.4"></script>
<link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-gray-900 text-gray-100 min-h-screen">
%s
<nav class="bg-gray-800 border-b border-gray-700 px-6 py-3 flex gap-6">
  <span class="font-bold text-yellow-400">⚡ x402Fuel</span>
  <a href="/" class="hover:text-yellow-300" hx-get="/" hx-target="#main" hx-push-url="true">Dashboard</a>
  <a href="/wallets" class="hover:text-yellow-300" hx-get="/wallets" hx-target="#main" hx-push-url="true">Wallets</a>
  <a href="/settings" class="hover:text-yellow-300" hx-get="/settings" hx-target="#main" hx-push-url="true">Settings</a>
</nav>
<main id="main" class="p-6 max-w-4xl mx-auto">
<h1 class="text-2xl font-bold mb-4"><span class="fuel-mini" aria-hidden="true">◇</span>Wallets</h1>
<div class="bg-gray-800 rounded-lg p-4">
<table class="w-full"><thead><tr class="text-left text-gray-400"><th>Address</th></tr></thead><tbody>%s</tbody></table>
</div>
<form class="mt-4 bg-gray-800 rounded-lg p-4" hx-post="/api/wallets" hx-swap="none">
  <input name="label" placeholder="Wallet label" class="bg-gray-700 rounded px-3 py-2 mr-2 text-white">
  <input name="passphrase" type="password" placeholder="Passphrase" class="bg-gray-700 rounded px-3 py-2 mr-2 text-white">
  <button type="submit" class="bg-yellow-500 text-black font-bold px-4 py-2 rounded hover:bg-yellow-400">Create Wallet</button>
</form>
</main></body></html>`

const settingsHTML = `<!DOCTYPE html><html lang="en">
<head><meta charset="UTF-8"><title>x402Fuel — Settings</title>
<script src="https://unpkg.com/htmx.org@2.0.4"></script>
<link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-gray-900 text-gray-100 min-h-screen">
%s
<nav class="bg-gray-800 border-b border-gray-700 px-6 py-3 flex gap-6">
  <span class="font-bold text-yellow-400">⚡ x402Fuel</span>
  <a href="/" class="hover:text-yellow-300" hx-get="/" hx-target="#main" hx-push-url="true">Dashboard</a>
  <a href="/wallets" class="hover:text-yellow-300" hx-get="/wallets" hx-target="#main" hx-push-url="true">Wallets</a>
  <a href="/settings" class="hover:text-yellow-300" hx-get="/settings" hx-target="#main" hx-push-url="true">Settings</a>
</nav>
<main id="main" class="p-6 max-w-4xl mx-auto">
<h1 class="text-2xl font-bold mb-4"><span class="fuel-mini" aria-hidden="true">◇</span>Settings</h1>
<div class="bg-gray-800 rounded-lg p-4 space-y-3">
  <div><span class="text-gray-400">Max Per Transaction:</span> <span class="font-mono">%s USDC</span></div>
  <div><span class="text-gray-400">Daily Cap:</span> <span class="font-mono">%s USDC</span></div>
  <div><span class="text-gray-400">Spent Today:</span> <span class="font-mono">%s USDC</span></div>
  <div><span class="text-gray-400">Paused:</span> <span class="font-bold">%s</span></div>
</div>
</main></body></html>`