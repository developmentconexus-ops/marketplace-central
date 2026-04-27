package shopee

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func signShopeeRequest(partnerID, partnerKey, path string, timestamp int64, shopID string) (string, error) {
	partnerID = strings.TrimSpace(partnerID)
	partnerKey = strings.TrimSpace(partnerKey)
	path = strings.TrimSpace(path)
	shopID = strings.TrimSpace(shopID)

	if partnerID == "" || partnerKey == "" || path == "" || timestamp <= 0 {
		return "", fmt.Errorf("invalid shopee signer input")
	}

	base := partnerID + path + fmt.Sprintf("%d", timestamp)
	if shopID != "" {
		base += shopID
	}

	mac := hmac.New(sha256.New, []byte(partnerKey))
	_, _ = mac.Write([]byte(base))
	return hex.EncodeToString(mac.Sum(nil)), nil
}
