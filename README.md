# go-ems-json-proxy

A tiny reverse proxy written in Go that:

- Accepts two positional arguments: EMS server IP/URL and proxy port.
- Starts an HTTP proxy server on the given port.
- Logs the request line and all headers to the console.
- Forwards the request to the EMS server, adjusting headers as needed for a proper reverse proxy (Host, X-Forwarded-*, hop-by-hop stripping handled by Go’s httputil.ReverseProxy).
- **Extra feature**: for POST requests with Content-Type: application/json, if the body is a single JSON string with escaped quotes (e.g. "{\"k\":\"v\"}"), it will unescape it and forward the sanitized JSON (e.g. {"k":"v"}).

Tip: You can pass http:// or https:// in the EMS argument. If you pass just an IP/host, it defaults to http://.




## Build & Run (Go)


```bash
go build -o ems-proxy ./cmd/ems-proxy
./ems-proxy --ems http://10.0.0.12:8081 --port 8080

```