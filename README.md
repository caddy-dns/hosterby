hoster.by module for Caddy
===========================

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records with [hoster.by](https://hoster.by) through [libdns/hosterby](https://github.com/libdns/hosterby).

## Caddy module name

```
dns.providers.hosterby
```

## Build

```sh
xcaddy build --with github.com/caddy-dns/hosterby
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
	"module": "acme",
	"challenges": {
		"dns": {
			"provider": {
				"name": "hosterby",
				"access_key": "{env.HOSTERBY_ACCESS_KEY}",
				"secret_key": "{env.HOSTERBY_SECRET_KEY}"
			},
			"propagation_delay": "2m"
		}
	}
}
```

or with the Caddyfile:

```
# globally
{
	cert_issuer acme {
		dns hosterby {
			access_key {env.HOSTERBY_ACCESS_KEY}
			secret_key {env.HOSTERBY_SECRET_KEY}
		}
		propagation_delay 2m
	}
}
```

```
# one site
tls {
	dns hosterby {
		access_key {env.HOSTERBY_ACCESS_KEY}
		secret_key {env.HOSTERBY_SECRET_KEY}
	}
	propagation_delay 2m
}
```

`access_key` and `secret_key` are the keys of a hoster.by API service account. For the ACME DNS challenge, the keys need permissions to list the DNS hosting services and to list, add, change and delete TXT records. See [libdns/hosterby](https://github.com/libdns/hosterby) for how to get the keys and for caveats.

At the time of writing, the hoster.by API fails to return the records of internationalized zones such as `.бел`, so certificates for them can't be obtained yet.

## Propagation delay

At the time of writing, a new record appears on u2.hoster.by about a minute after it appears on u1.hoster.by (59 to 60 seconds in four measurements). Caddy considers the record propagated once one nameserver returns it, so without a delay Let's Encrypt queries u2 too early and the challenge fails with NXDOMAIN. A `propagation_delay` of 2 minutes worked in all tests.
