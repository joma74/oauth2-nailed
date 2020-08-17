# Replay after Udemy "OAuth2.0 : Nailed the core framework"

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

- myrealm
- Endpoints: OpenID Endpoint

So that

<img src="./docs/KeycloakOAuthWellKnown.png" alt="Keycloaks Well-known Openid Configuration"
	title="Keycloaks Well-known Openid Configuration" width="700" height="auto" />

### User

- myuser/myuser
- Email Verified: Off

## Starting the OAuth Client

```
joma@edison:oauth2-nailed (master%=) $ cd src/client/
joma@edison:client (master%=) $ go run .
```

The OAuth Client page is then reachable via http://localhost:9110/

## References

- https://www.keycloak.org/docs-api/11.0/rest-api/index.html
- https://github.com/keycloak/keycloak-containers/blob/11.0.0/server/README.md
