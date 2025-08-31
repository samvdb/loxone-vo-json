package proxy

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"moul.io/http2curl"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	rp := httputil.NewSingleHostReverseProxy(target)

	rp.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		// Handle JSON un-escaping for POST bodies with application/json and no content-encoding.
		if strings.EqualFold(req.Method, http.MethodPost) {
			if ct := req.Header.Get("Content-Type"); ct != "" {
				if mt, _, err := mime.ParseMediaType(ct); err == nil && mt == "application/json" {
					if req.Header.Get("Content-Encoding") == "" { // skip gzipped/etc.
						command, _ := http2curl.GetCurlCommand(req)
						fmt.Println(command)
						b, err := io.ReadAll(req.Body)
						if err == nil {
							_ = req.Body.Close()
							clean := SanitizeJSONBody(b)
							if !bytes.Equal(b, clean) {
								slog.Info("sanitized JSON body by unquoting escaped payload")
							}
							fmt.Print(string(clean))
							req.Body = io.NopCloser(bytes.NewReader(clean))
							req.ContentLength = int64(len(clean))
							req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))
						}
					}
				}
			}
		}
	}

	rp.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
		ForceAttemptHTTP2:     true,
	}

	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		slog.Error("proxy error", "error", err)
		http.Error(w, "proxy error", http.StatusBadGateway)
	}

	return rp
}
