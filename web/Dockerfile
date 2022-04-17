FROM caddy:latest

COPY dist /usr/share/caddy/
COPY Caddyfile /etc/caddy/Caddyfile

ENV APIADDRESS = "http://127.0.0.1:8080"

CMD cat /etc/caddy/Caddyfile && echo $APIADDRESS > /usr/share/caddy/api_address && caddy run --config /etc/caddy/Caddyfile --adapter caddyfile