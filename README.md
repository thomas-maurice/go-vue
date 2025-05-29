# Go Vue JS test app with users

This is a test VueJS app, it does nothing besides logging users in and showing them their profiles.

You can login using a local user (in the DB) or through OIDC

The only local admin right now is `admin:admin`

For OIDC you connect your provider and see the magic happen.

## How to build

:warning: You need [swag](https://github.com/swaggo/swag) and docker installed

```
$ make
```
Should take care of that. The binary is self suficient and will serve the UI

### Build steps explained

1. We do a swaggo that builds the swagger
2. We do a swagger to javascript shenanigan that builds the js api client the UI will use
3. We bundle that up for the go program to serve
4. We build the binary

## Sample config file
```yaml
debug: true
storage:
  driver: sqlite3
  url: db.sqlite3
http:
  listen: :8080
security:
  # it's "admin" in normal speak
  adminPassword: $2a$12$iUfsLM1ZPqjAFuUFhA1.aeBMbIkFCHb.2iJs9u/IzQCp1CqES39LW
  oidc:
    authentik:
      display_name: Authentik OIDC
      issuer: your.issuer
      clientId: your.client.id
      clientSecret: [redacted]
      scopes:
        - openid
        - profile
        - email
        - groups
  # generate a new one using go run main.go genkey
  signingKey: |
    -----BEGIN ECDSA PRIVATE KEY-----
    MIHcAgEBBEIAYlrVDOmSjvhlP6rdn5zFl+PuqSQPw6XaD4ysKujIHownoj0voXpl
    zqkQ1Tgw6OvFtZ3FLBEGyHorL3y7g1H2HdKgBwYFK4EEACOhgYkDgYYABADhuAVE
    VQxC6T/x+5jqhWBpafRjMOjcdi70dmPzcVW+XA68VBvOQQ3+PDojv5npCruzBZJ/
    dIBB1Slyb5Bqi/5OngAJ5m+D+ab19gBr8MGZza6mn1QUrPWur2Fbj+2AItyDysJj
    EkQjYJKyjt/x3nEFmTnnHz6s9XPAQ6KOoNuvSQg6Wg==
    -----END ECDSA PRIVATE KEY-----
```
