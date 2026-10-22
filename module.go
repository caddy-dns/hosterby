// Package hosterby is a Caddy DNS provider module for hoster.by, built on
// github.com/libdns/hosterby.
package hosterby

import (
	"errors"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/certmagic"
	libdnshosterby "github.com/libdns/hosterby"
)

// Provider lets Caddy read and manipulate DNS records hosted by hoster.by.
type Provider struct{ *libdnshosterby.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.hosterby",
		New: func() caddy.Module { return &Provider{new(libdnshosterby.Provider)} },
	}
}

// Provision replaces placeholders such as {env.HOSTERBY_ACCESS_KEY} in the
// keys. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.AccessKey = repl.ReplaceAll(p.Provider.AccessKey, "")
	p.Provider.SecretKey = repl.ReplaceAll(p.Provider.SecretKey, "")
	if p.Provider.AccessKey == "" || p.Provider.SecretKey == "" {
		return errors.New("access_key and secret_key must not be empty")
	}
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	hosterby {
//	    access_key <access_key>
//	    secret_key <secret_key>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			var field *string
			switch d.Val() {
			case "access_key":
				field = &p.Provider.AccessKey
			case "secret_key":
				field = &p.Provider.SecretKey
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
			if *field != "" {
				return d.Errf("%s already set", d.Val())
			}
			if !d.AllArgs(field) {
				return d.ArgErr()
			}
		}
	}
	if p.Provider.AccessKey == "" {
		return d.Err("missing access_key")
	}
	if p.Provider.SecretKey == "" {
		return d.Err("missing secret_key")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
	_ certmagic.DNSProvider = (*Provider)(nil)
)
