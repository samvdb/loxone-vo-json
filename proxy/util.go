package proxy

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
)

func SanitizeJSONBody(b []byte) []byte {
	orig := string(bytes.TrimSpace(b))
	s := orig

	// 1) If the whole payload is quoted JSON, keep unquoting until stable.
	for {
		unq, err := strconv.Unquote(s)
		if err != nil {
			break
		}
		s = unq
	}
	if json.Valid([]byte(s)) {
		if s != orig {
			slog.Info("sanitized JSON body (unquote loop)", "before", orig, "after", s)
		}
		return []byte(s)
	}

	// 2) Wrap-and-unquote strategy: treat the payload as a Go string literal
	// so that sequences like \" become ". This helps for bodies like:
	// {\"cmd\":\"dhw.onetime\",\"data\":\"on\"}
	if unq, err := strconv.Unquote(`"` + s + `"`); err == nil && json.Valid([]byte(unq)) {
		slog.Info("sanitized JSON body (wrap-and-unquote)", "before", orig, "after", unq)
		return []byte(unq)
	}

	// 3) Conservative fallback: replace \" with " and validate.
	// Only accept if it becomes valid JSON.
	candidate := strings.ReplaceAll(s, `\"`, `"`)
	if json.Valid([]byte(candidate)) {
		slog.Info("sanitized JSON body (replace-escaped-quotes)", "before", orig, "after", candidate)
		return []byte(candidate)
	}

	// 4) Give up: return original.
	slog.Info("left JSON body unchanged (no valid sanitation)", "body", orig)
	return b
}

func ParseTarget(raw string) (*url.URL, error) {
	// Allow plain host/IP (default to http), or full URL with scheme.
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return url.Parse(raw)
	}
	return url.Parse("http://" + raw)
}
