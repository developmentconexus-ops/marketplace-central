package shopee

import "testing"

func TestSignShopeeRequestAuthorize(t *testing.T) {
	t.Parallel()

	got, err := signShopeeRequest("10090", "test-partner-key", "/api/v2/shop/auth_partner", 1594897040, "")
	if err != nil {
		t.Fatalf("signShopeeRequest() error = %v", err)
	}
	if got, want := got, "5c3357a5d7ede41e954ce12fe4a885c77dd0f4ff59666ddf16fbab373fdfe636"; got != want {
		t.Fatalf("signShopeeRequest() = %q, want %q", got, want)
	}
}

func TestSignShopeeRequestTokenRequests(t *testing.T) {
	t.Parallel()

	got, err := signShopeeRequest("10090", "test-partner-key", "/api/v2/auth/token/get", 1594897040, "9001")
	if err != nil {
		t.Fatalf("signShopeeRequest() error = %v", err)
	}
	if got, want := got, "fdca6bda23a875e5a2b836c825009085487ff88af73c681b6019d5a5b70d3b34"; got != want {
		t.Fatalf("signShopeeRequest() = %q, want %q", got, want)
	}

	got, err = signShopeeRequest("10090", "test-partner-key", "/api/v2/auth/access_token/get", 1594897040, "9001")
	if err != nil {
		t.Fatalf("signShopeeRequest() error = %v", err)
	}
	if got, want := got, "095229669dda19f50e82da0c4b2067215c24f11b582282705d507fc117a79838"; got != want {
		t.Fatalf("signShopeeRequest() = %q, want %q", got, want)
	}
}
