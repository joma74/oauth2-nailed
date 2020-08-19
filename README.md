# Replay after Udemy "OAuth2.0 : Nailed the core framework"

See https://www.udemy.com/course/oauth-2-nailed-the-core-framework-with-hands-dirty

## Treated OAuth 2.0 Flows

### The Abstract OAuth 2.0 Flow

> The abstract OAuth 2.0 flow illustrated describes the interaction between the four roles.

<img src="./docs/OAuth2.0Standard_AbstractProtocolFlow.png" alt="The abstract OAuth 2.0 flow illustrated describes the interaction between the four roles."
	title="The abstract OAuth 2.0 flow illustrated describes the interaction between the four roles." width="700" height="auto" />

See https://tools.ietf.org/html/rfc6749#section-1.2

### The Authorization Code Grant

> The authorization code grant type is used to obtain both access
> tokens and refresh tokens and is optimized for confidential clients.
>
> Since this is a redirection-based flow, the client must be capable of
> interacting with the resource owner's user-agent (typically a web
> browser) and capable of receiving incoming requests (via redirection)
> from the authorization server.

<img src="./docs/OAuth2.0Standard_AuthorizationCodeGrant.png" alt="Authorization Code Flow"
	title="The authorization code grant type is used to obtain both access
   tokens and refresh tokens and is optimized for confidential clients." width="700" height="auto" />

See https://tools.ietf.org/html/rfc6749#section-4.1

## Setup Keycloak on Docker

```
docker run -p 9112:8080 -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin -e TZ=Europe/Vienna quay.io/keycloak/keycloak:11.0.0
docker stop amazing_kapitsa
docker rename amazing_kapitsa keycloak_1
docker start keycloak_1
docker logs -f  --tail 20  keycloak_1
```

The admin interface is then reachable via http://localhost:9112/auth/

## Administer a new Keycloak realm and users

### Realm

- Name: myrealm
- Client: oauth-nailed-app-1
  - Root URL: http://localhost:9110/
  - Valid Redirect URIs: http://localhost:9110/*
  - Admin URL: http://localhost:9110/
  - Web Origins: http://localhost:9110/
- Endpoints: OpenID Endpoint

So that

<img src="./docs/KeycloakOAuthWellKnown.png" alt="Keycloaks Well-known Openid Configuration"
	title="Keycloaks Well-known Openid Configuration" width="700" height="auto" />

### User

- Name/PWD: myuser/myuser
- Email Verified: Off

## Starting the OAuth Client

```
cd src/client/
go run .
```

The OAuth Client page is then reachable via http://localhost:9110/

It covers the flow of

- A (Authorization Request, 4.1.1)
- C (Authorization Response, 4.1.2)
- D (Access Token Request, 4.1.3)
- E (Access Token Response, 4.1.4)

* E (Accessing Protected Resources, 7.)
* F (Accessing Protected Resources, 7.)

_The numbers reference the related section in https://tools.ietf.org/html/rfc6749_

## Starting the OAuth Protected Resource

```
cd src/billingservice/
go run .
```

## References

- https://tools.ietf.org/html/rfc6749 (The OAuth 2.0 Authorization Framework Standard)
  - https://tools.ietf.org/html/rfc6749#section-4.1 (Authorization Code Grant)
  - https://tools.ietf.org/html/rfc6749#section-7 (Accessing Protected Resources)
- https://tools.ietf.org/html/rfc6750 (The OAuth 2.0 Bearer Token Usage)
- https://github.com/keycloak/keycloak-containers/blob/11.0.0/server/README.md
- https://www.keycloak.org/docs-api/11.0/rest-api/index.html
- https://www.keycloak.org/docs/latest/securing_apps

## Tooling

- http://json2struct.mervine.net (Derive a Go struct type from a Json instance)
