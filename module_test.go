package hosterby

import (
	"encoding/json"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	libdnshosterby "github.com/libdns/hosterby"
)

func TestUnmarshalCaddyfile(t *testing.T) {
	for _, tc := range []struct {
		input   string
		wantErr bool
	}{
		{"hosterby {\n access_key ak\n secret_key sk\n}", false},
		{"hosterby {\n secret_key sk\n access_key ak\n}", false},
		{"hosterby", true},
		{"hosterby ak sk", true},
		{"hosterby {\n access_key ak\n}", true},
		{"hosterby {\n secret_key sk\n}", true},
		{"hosterby {\n access_key\n secret_key sk\n}", true},
		{"hosterby {\n access_key ak ak2\n secret_key sk\n}", true},
		{"hosterby {\n access_key ak\n access_key ak\n secret_key sk\n}", true},
		{"hosterby {\n access_key ak\n secret_key sk\n base_url http://localhost\n}", true},
	} {
		p := Provider{new(libdnshosterby.Provider)}
		err := p.UnmarshalCaddyfile(caddyfile.NewTestDispenser(tc.input))
		if (err != nil) != tc.wantErr {
			t.Errorf("%q: got error %v, want error %v", tc.input, err, tc.wantErr)
			continue
		}
		if err == nil && (p.AccessKey != "ak" || p.SecretKey != "sk") {
			t.Errorf("%q: got access key %q, secret key %q", tc.input, p.AccessKey, p.SecretKey)
		}
	}
}

func TestProvision(t *testing.T) {
	t.Setenv("HOSTERBY_TEST_ACCESS_KEY", "ak")
	t.Setenv("HOSTERBY_TEST_SECRET_KEY", "sk")

	p := Provider{new(libdnshosterby.Provider)}
	config := `{"access_key": "{env.HOSTERBY_TEST_ACCESS_KEY}", "secret_key": "{env.HOSTERBY_TEST_SECRET_KEY}"}`
	if err := json.Unmarshal([]byte(config), &p); err != nil {
		t.Fatal(err)
	}
	if err := p.Provision(caddy.Context{}); err != nil {
		t.Fatal(err)
	}
	if p.AccessKey != "ak" || p.SecretKey != "sk" {
		t.Errorf("got access key %q, secret key %q", p.AccessKey, p.SecretKey)
	}

	p = Provider{&libdnshosterby.Provider{AccessKey: "ak", SecretKey: "{env.HOSTERBY_TEST_UNSET}"}}
	if err := p.Provision(caddy.Context{}); err == nil {
		t.Error("an empty key after replacing placeholders must be an error")
	}
}
