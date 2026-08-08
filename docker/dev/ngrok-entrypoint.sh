#!/bin/sh
set -eu

if [ -z "${NGROK_AUTHTOKEN:-}" ]; then
  echo "NGROK_AUTHTOKEN is required to start the reserved Mercado Livre callback tunnel." >&2
  exit 2
fi

redirect_uri="${MPC_OAUTH_REDIRECT_URI:-}"
if [ -z "$redirect_uri" ]; then
  echo "MPC_OAUTH_REDIRECT_URI is required so ngrok can use the registered callback host." >&2
  exit 2
fi

callback_host="${redirect_uri#http://}"
callback_host="${callback_host#https://}"
callback_host="${callback_host%%/*}"

if [ -z "$callback_host" ] || [ "$callback_host" = "$redirect_uri" ]; then
  echo "MPC_OAUTH_REDIRECT_URI must be an absolute http(s) URL." >&2
  exit 2
fi

# Target is oauth-edge, NOT frontend:5174. The Vite dev server proxies the whole
# route table to the backend, so pointing the tunnel at it published /orders —
# buyer name, CPF/CNPJ and address, with no identity check anywhere in the Go
# middleware chain — to the public internet for as long as the tunnel was up
# (issue #1). oauth-edge whitelists the OAuth callback and denies the rest.
exec ngrok http --url="$callback_host" oauth-edge:80
