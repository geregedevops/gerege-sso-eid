package dan

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

// ComputeHMAC generates an HMAC-SHA256 signature over sorted key=value pairs.
// The "signature" and "image" keys are excluded from the computation.
func ComputeHMAC(data map[string]string, hmacKey string) string {
	params := url.Values{}
	for k, v := range data {
		if k != "signature" && k != "image" {
			params.Set(k, v)
		}
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(params.Get(k)))
	}

	mac := hmac.New(sha256.New, []byte(hmacKey))
	mac.Write([]byte(buf.String()))
	return hex.EncodeToString(mac.Sum(nil))
}
