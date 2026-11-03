#!/usr/bin/env bash
# Rasteriza a marca do Sintaxe (web/static/favicon.svg) nos PNGs que o PWA usa.
# Chrome headless porque é o mesmo motor que desenha o SVG no browser: o PNG
# instalado sai idêntico ao que aparece na aba. O SVG é a fonte única — nada de
# redesenhar a marca aqui.
#
# Uso: scripts/build-icons.sh [--chrome <caminho>]
set -euo pipefail

cd "$(dirname "$0")/.."

CHROME="${CHROME:-/Applications/Google Chrome.app/Contents/MacOS/Google Chrome}"
if [ "${1:-}" = "--chrome" ] && [ -n "${2:-}" ]; then
  CHROME="$2"
fi
if [ ! -x "$CHROME" ]; then
  echo "chrome não encontrado em $CHROME (passe --chrome <caminho> ou defina CHROME)" >&2
  exit 1
fi

SVG="$(pwd)/web/static/favicon.svg"
OUT="web/static/icons"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

cat >"$TMP/page.html" <<HTML
<!doctype html>
<html><head><meta charset="utf-8"><style>
  html, body { margin: 0; padding: 0; background: #161513; }
  img { display: block; }
</style></head><body><img id="mark" src="file://$SVG" width="512" height="512"></body></html>
HTML

render() { # size file
  local size="$1" file="$2"
  "$CHROME" --headless=new --disable-gpu --hide-scrollbars --force-device-scale-factor=1 \
    --window-size="$size,$size" --screenshot="$TMP/$file" "file://$TMP/page.html#s=$size" >/dev/null 2>&1
  # o tamanho vem no atributo: reescreve a página para cada medida
  cat >"$TMP/page.html" <<HTML
<!doctype html>
<html><head><meta charset="utf-8"><style>
  html, body { margin: 0; padding: 0; background: #161513; }
  img { display: block; }
</style></head><body><img src="file://$SVG" width="$size" height="$size"></body></html>
HTML
  mv "$TMP/$file" "$OUT/$file"
  printf '  %s (%sx%s)\n' "$OUT/$file" "$size" "$size"
}

echo "→ rasterizando web/static/favicon.svg"
render 512 icon-512.png
render 512 icon-512-maskable.png
render 192 icon-192.png
render 180 icon.png
render 32 favicon-32.png
echo "→ pronto. Mudou o desenho? Suba o ?v= nas referências (manifest, layout, sw.js) para furar o cache."
