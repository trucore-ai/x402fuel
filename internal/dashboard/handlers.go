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
		w.Write([]byte(indexHTML))
	})

	r.HandleFunc("/wallets", func(w http.ResponseWriter, r *http.Request) {
		addrs, _ := ks.List()
		var rows strings.Builder
		for _, addr := range addrs {
			rows.WriteString(fmt.Sprintf(`<tr><td class="font-mono text-sm">%s</td></tr>`, addr))
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(walletsHTML, rows.String())))
	})

	r.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		p := pol.GetPolicy()
		paused := "No"
		if p.Paused {
			paused = "Yes \u23f8"
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(settingsHTML, p.MaxPerTxn, p.DailyCap, p.DailySpent, paused)))
	})
}

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>x402Fuel</title>
<script src="https://unpkg.com/htmx.org@2.0.4"></script>
<link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-gray-900 text-gray-100 min-h-screen">
<nav class="bg-gray-800 border-b border-gray-700 px-6 py-3 flex gap-6">
  <span class="font-bold text-yellow-400">⚡ x402Fuel</span>
  <a href="/" class="hover:text-yellow-300" hx-get="/" hx-target="#main" hx-push-url="true">Dashboard</a>
  <a href="/wallets" class="hover:text-yellow-300" hx-get="/wallets" hx-target="#main" hx-push-url="true">Wallets</a>
  <a href="/settings" class="hover:text-yellow-300" hx-get="/settings" hx-target="#main" hx-push-url="true">Settings</a>
</nav>
<main id="main" class="p-6 max-w-4xl mx-auto">
  <h1 class="text-2xl font-bold mb-4">Dashboard</h1>
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
<nav class="bg-gray-800 border-b border-gray-700 px-6 py-3 flex gap-6">
  <span class="font-bold text-yellow-400">⚡ x402Fuel</span>
  <a href="/" class="hover:text-yellow-300" hx-get="/" hx-target="#main" hx-push-url="true">Dashboard</a>
  <a href="/wallets" class="hover:text-yellow-300" hx-get="/wallets" hx-target="#main" hx-push-url="true">Wallets</a>
  <a href="/settings" class="hover:text-yellow-300" hx-get="/settings" hx-target="#main" hx-push-url="true">Settings</a>
</nav>
<main id="main" class="p-6 max-w-4xl mx-auto">
<h1 class="text-2xl font-bold mb-4">Wallets</h1>
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
<nav class="bg-gray-800 border-b border-gray-700 px-6 py-3 flex gap-6">
  <span class="font-bold text-yellow-400">⚡ x402Fuel</span>
  <a href="/" class="hover:text-yellow-300" hx-get="/" hx-target="#main" hx-push-url="true">Dashboard</a>
  <a href="/wallets" class="hover:text-yellow-300" hx-get="/wallets" hx-target="#main" hx-push-url="true">Wallets</a>
  <a href="/settings" class="hover:text-yellow-300" hx-get="/settings" hx-target="#main" hx-push-url="true">Settings</a>
</nav>
<main id="main" class="p-6 max-w-4xl mx-auto">
<h1 class="text-2xl font-bold mb-4">Settings</h1>
<div class="bg-gray-800 rounded-lg p-4 space-y-3">
  <div><span class="text-gray-400">Max Per Transaction:</span> <span class="font-mono">%s USDC</span></div>
  <div><span class="text-gray-400">Daily Cap:</span> <span class="font-mono">%s USDC</span></div>
  <div><span class="text-gray-400">Spent Today:</span> <span class="font-mono">%s USDC</span></div>
  <div><span class="text-gray-400">Paused:</span> <span class="font-bold">%s</span></div>
</div>
</main></body></html>`