---
title: Grant Types
---

For a list of supported or unsupported `Grant Types` please have a look at the table below.

| Grant Type                                            | Supported           |
|:------------------------------------------------------|:--------------------|
| Authorization Code                                    | yes                 |
| Authorization Code with PKCE                          | yes                 |
| Client Credentials                                    | no                  |
| Device Authorization                                  | under consideration |
| Implicit                                              | yes                 |
| JSON Web Token (JWT) Profile                          | yes                 |
| Refresh Token                                         | yes                 |
| Resource Owner Password Credentials                   | no                  |
| Security Assertion Markup Language (SAML) 2.0 Profile | no                  |
| Token Exchange                                        | no    |

## Authorization Code

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.3.1](https://tools.ietf.org/html/rfc6749#section-1.3.1)

## Proof Key for Code Exchange

**Link to spec.** [Proof Key for Code Exchange by OAuth Public Clients](https://tools.ietf.org/html/rfc7636)

## Implicit

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.3.2](https://tools.ietf.org/html/rfc6749#section-1.3.2)

## Client Credentials

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.3.4](https://tools.ietf.org/html/rfc6749#section-1.3.4)

## Refresh Token

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.5](https://tools.ietf.org/html/rfc6749#section-1.5)

## JSON Web Token (JWT) Profile

**Link to spec.** [JSON Web Token (JWT) Profile for OAuth 2.0 Client Authentication and Authorization Grants](https://tools.ietf.org/html/rfc7523)

### Using JWTs as Authorization Grants

Our service user work with the JWT profile to authenticate them against ZITADEL.

1. Create or use an existing service user
2. Create a new key and download it
3. Generate a JWT with the structure below and sign it with the downloaded key
4. Send the JWT Base64 encoded to ZITADEL's token endpoint
5. Use the received access token

---

Key JSON

| Key    | Example                                                             | Description                                                        |
|:-------|:--------------------------------------------------------------------|:-------------------------------------------------------------------|
| type   | `"serviceaccount"`                                                  | The type of account, right now only serviceaccount is valid        |
| keyId  | `"81693565968772648"`                                               | This is unique ID of the key                                       |
| key    | `"-----BEGIN RSA PRIVATE KEY-----...-----END RSA PRIVATE KEY-----"` | The private key generated by ZITADEL, this can not be regenerated! |
| userId | `78366401571647008`                                                 | The service users ID, this is the same as the subject from tokens  |

```JSON
{
	"type": "serviceaccount",
	"keyId": "81693565968772648",
	"key": "-----BEGIN RSA PRIVATE KEY-----...-----END RSA PRIVATE KEY-----",
	"userId": "78366401571647008"
}
```

---

JWT

| Claim | Example                   | Description                                                                                                   |
|:------|:--------------------------|:--------------------------------------------------------------------------------------------------------------|
| aud   | `"https://{your_domain}"` | String or Array of intended audiences MUST include ZITADEL's issuing domain                                   |
| exp   | `1605183582`              | Unix timestamp of the expiry                                                                                  |
| iat   | `1605179982`              | Unix timestamp of the creation singing time of the JWT, MUST NOT be older than 1h                             |
| iss   | `"77479219772321307"`     | String which represents the requesting party (owner of the key), normally the `userId` from the json key file |
| sub   | `"77479219772321307"`     | The subject ID of the service user, normally the `userId` from the json key file                              |

```JSON
{
	"iss": "77479219772321307",
	"sub": "77479219772321307",
	"aud": "https://{your_domain}",
	"exp": 1605183582,
	"iat": 1605179982
}
```

> To identify your key, it is necessary that you provide a JWT with a `kid` header claim representing your keyId from the Key JSON:
>
> ```json
> {
> 	"alg": "RS256",
> 	"kid": "81693565968772648"
> }
> ```

---

See [JWT Profile Grant on Token Endpoint](endpoints#token_endpoint) for usage.

### Using JWTs for Client Authentication

See how to build a [JWT for client authentication](authn-methods#jwt-with-private-key) from the downloaded key.

Find out how to use it on the [token endpoint](endpoints#token_endpoint) or the [introspection endpoint](endpoints#introspection_endpoint).

## Token Exchange

**Link to spec.** [OAuth 2.0 Token Exchange](https://tools.ietf.org/html/rfc8693)

## Device Authorization

**Link to spec.** [OAuth 2.0 Device Authorization Grant](https://tools.ietf.org/html/rfc8628)

## Not Supported Grant Types

### Resource Owner Password Credentials

> Due to growing security concerns we do not support this grant type. With OAuth 2.1 it looks like this grant will be removed.

**Link to spec.** [OThe OAuth 2.0 Authorization Framework Section 1.3.3](https://tools.ietf.org/html/rfc6749#section-1.3.3)

### Security Assertion Markup Language (SAML) 2.0 Profile

**Link to spec.** [Security Assertion Markup Language (SAML) 2.0 Profile for OAuth 2.0 Client Authentication and Authorization Grants](https://tools.ietf.org/html/rfc7522)
