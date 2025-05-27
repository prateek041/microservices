---
title: "JWT Auth"
description: "JWT Authentication mechanisms"
date: 2025-05-27
---

Write more in-depth of various ways of creating and signing tokens. I should
refer my previous articles here.

## Claims

In JSON Web Tokens (JWTs), claims are pieces of information about the user and
the token itself that are encoded within the token's payload. They are key-value
pairs represented as a JSON object.

There are three types of claims:

- **Registered Claims**: These are a set of predefined claim names that are
  recommended for interoperability. They are not mandatory but are often used.
  Examples include:

  - `iss` (issuer): Identifies the principal that issued the JWT.
  - `sub` (subject): Identifies the principal that is the subject of the JWT.
  - `aud` (audience): Identifies the recipients for which the JWT is intended.
  - `exp` (expiration time): Identifies the time after which the JWT MUST NOT be
    accepted for processing.
  - `nbf` (not before): Identifies the time before which the JWT MUST NOT be
    accepted for processing.
  - `iat` (issued at): Identifies the time at which the JWT was issued.
  - `jti` (JWT ID): A unique identifier for the JWT.

- **Public Claims**: These are claim names that are defined in the JWT
  specification but are not registered. To avoid collisions, it's recommended to
  define them in a namespace or use a URL format. For example, you might define
  a public claim like <https://example.com/roles>.

- **Private Claims**: These are custom claim names that are specific to your
  application. You can include any relevant information here, such as user IDs,
  roles, permissions, or any other application-specific data.

Example:

Let's say a user with ID `user123` logs into our system. The User Management
Service might generate a JWT with the following claims in its payload (represented
as a JSON object):

```json
{
  "iss": "user-management-service",
  "sub": "user123",
  "aud": "api-gateway",
  "exp": 1678886400,
  "iat": 1678882800,
  "user_id": "user123",
  "roles": ["user", "premium"]
}
```

Breakdown of the Example Claims:

- iss: The issuer of this token is user-management-service.
- sub: The subject of this token is the user with ID user123.
- aud: This token is intended for the api-gateway (the audience).
- exp: This token will expire at the Unix timestamp 1678886400 (e.g., March 15,
  2023, at 8:00 AM UTC). After this time, the API Gateway should reject this token.
- iat: This token was issued at the Unix timestamp 1678882800 (e.g., March 15,
  2023, at 7:00 AM UTC).
- user_id (Private Claim): Our application's specific user ID.
- roles (Private Claim): An array indicating the roles of this user (user and premium).

When the API Gateway receives this JWT, after verifying its signature, it can read
these claims from the payload. It can then use this information for:

- **Identifying the user**: The sub or `user_id` claim tells the gateway which user
  the request is for.
- **Authorization (basic)**: The roles claim could be used by the gateway to
  determine if the user has the necessary permissions to access a certain
  resource. For more complex authorization, something like an Open Policy Agent,
  the gateway might forward the token or user information to the OPA service.
- **Audience verification**: The gateway can check if the `aud` claim matches
  itself to ensure the token was intended for it.
- **Expiration check**: The gateway must ensure the current time is before the
  exp time.

In our Login handler, we are currently setting the sub, issuedAt, and expiresAt
registered claims, and you can see where you would add private claims like `user_id`
or roles to the claims map before creating the token.
