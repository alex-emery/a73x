FROM nginx:1.27.0-alpine

RUN apk update && apk add ca-certificates iptables ip6tables && rm -rf /var/cache/apk/*

COPY --from=docker.io/tailscale/tailscale:stable /usr/local/bin/tailscaled /app/tailscaled
COPY --from=docker.io/tailscale/tailscale:stable /usr/local/bin/tailscale /app/tailscale
RUN mkdir -p /var/run/tailscale /var/cache/tailscale /var/lib/tailscale

WORKDIR /app/

COPY nginx.conf /etc/nginx/conf.d/nginx.conf
COPY /start.sh /app/start.sh

RUN chmod +x /app/start.sh

CMD ["/app/start.sh"]
