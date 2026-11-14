package panel

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/cuijianzhuang/polaris/internal/store"
)

func TestPortalDevicesArePrivateAndRemovable(t *testing.T) {
	s, base, _ := cdnRegionFixture(t)
	now := time.Now().UTC()
	for _, d := range []store.UserDevice{
		{ID: "dev-alice", UserID: "ra", HWID: "alice-private-hardware-id", Model: "Alice phone", OS: "iOS", IP: "198.51.100.10", FirstSeen: now, LastSeen: now},
		{ID: "dev-bob", UserID: "rb", HWID: "bob-private-hardware-id", Model: "Bob phone", OS: "Android", IP: "198.51.100.11", FirstSeen: now, LastSeen: now},
	} {
		if err := s.st.DB.Create(&d).Error; err != nil {
			t.Fatal(err)
		}
	}
	unauth, _ := http.Get(base + portalPath + "/devices")
	if unauth.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated status = %d", unauth.StatusCode)
	}
	unauth.Body.Close()

	alice := regionPortalClient(t, s, base, "ra")
	if code, body := alice.req("GET", portalPath+"/devices", nil, false); code != http.StatusOK || !strings.Contains(body, "dev-alice") || strings.Contains(body, "dev-bob") || strings.Contains(body, "alice-private-hardware-id") {
		t.Fatalf("portal devices = %d %s", code, body)
	}
	if code, _ := alice.req("DELETE", portalPath+"/devices/dev-alice", nil, false); code != http.StatusForbidden {
		t.Fatalf("delete without csrf = %d", code)
	}
	if code, _ := alice.req("DELETE", portalPath+"/devices/dev-bob", nil, true); code != http.StatusNotFound {
		t.Fatalf("delete other's device = %d", code)
	}
	if code, _ := alice.req("DELETE", portalPath+"/devices/dev-alice", nil, true); code != http.StatusNoContent {
		t.Fatalf("delete own device = %d", code)
	}
	var count int64
	s.st.DB.Model(&store.UserDevice{}).Where("id = ?", "dev-alice").Count(&count)
	if count != 0 {
		t.Fatal("own device was not deleted")
	}
	s.st.DB.Model(&store.UserDevice{}).Where("id = ?", "dev-bob").Count(&count)
	if count != 1 {
		t.Fatal("other user's device was deleted")
	}
}

func TestPortalLanguagePreference(t *testing.T) {
	s, base, _ := cdnRegionFixture(t)
	alice := regionPortalClient(t, s, base, "ra")
	if code, _ := alice.req(http.MethodPost, portalPath+"/language", map[string]string{"lang": "en"}, true); code != http.StatusNoContent {
		t.Fatalf("set language = %d", code)
	}
	if code, body := alice.req(http.MethodGet, portalPath, nil, false); code != http.StatusOK || !strings.Contains(body, `<html lang="en">`) || !strings.Contains(body, `>Dashboard<`) || !strings.Contains(body, `>Subscription<`) || !strings.Contains(body, `>Nodes<`) {
		t.Fatalf("English portal = %d %s", code, body)
	}
	if code, _ := alice.req(http.MethodPost, portalPath+"/language", map[string]string{"lang": "fr"}, true); code != http.StatusBadRequest {
		t.Fatalf("invalid language = %d", code)
	}
}

func TestHWIDHeadersPreferStandardFamily(t *testing.T) {
	s := newTestServer(t)
	u := store.User{ID: store.NewID(), Username: "d", UUID: store.NewID(), Password: "p", SubToken: store.NewID(), Status: "active"}
	if err := s.st.DB.Create(&u).Error; err != nil {
		t.Fatal(err)
	}
	r, _ := http.NewRequest(http.MethodGet, "/sub/x", nil)
	r.Header.Set("X-HWID", "phone-1")
	r.Header.Set("X-HWID-OS", "iOS")
	r.Header.Set("X-HWID-Model", "iPhone")
	if !s.checkDevice(&u, r) {
		t.Fatal("standard HWID headers refused")
	}
	var got store.UserDevice
	if err := s.st.DB.First(&got, "user_id = ? AND hwid = ?", u.ID, "phone-1").Error; err != nil {
		t.Fatal(err)
	}
	if got.OS != "iOS" || got.Model != "iPhone" {
		t.Fatalf("metadata = %+v", got)
	}
}
