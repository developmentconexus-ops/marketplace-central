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

exec ngrok http --url="$callback_host" frontend:5174
