<h1>Docs for Auth EndPoints 🔐</h1>

<p>Note: requests to protected resources may require a <a href="https://github.com/golang-jwt/jwt">JWT</a> token</p>

<h3>Sign Up EndPoint 🔑</h3>
<p>This endpoint exists for the sole reason of creating accounts for users.</p>
<p>It creates a session by exchanging a JWT and refresh token with the client</p>
<p>It would also attempt to log the user in.</p>
<p>Make a POST request to http://localhost:3000/users</p>
<p>Ensure that the Content-Type is "application/json"</p>
<p>The expected body format:</p>

```
{
    "role": "student",
    "email": "example@exmail.ext",
    "password": "password"
}
```

<p>Valid roles would include one of {"student", "teacher", "admin"}</p>

<p>Expected response body for a successful request:</p>

```
{
    "status":  "success",
	"message": "user created",
	"jwt":     jwtToken,
	"token":   refreshToken,
}

```

<p>Expected response body for failed request:</p>

```
{
    "status" : "failed"
    "message": <error_description>
}
```

<h3>Login EndPoint 🔓</h3>
<p>This endpoint exists for the sole reason of logging users in.</p>
<p>It creates a session by exchanging a JWT and refresh token with the client</p>
<p>Make a POST request to http://localhost:3000/login</p>
<p>Ensure that the Content-Type is "application/json"</p>
<p>The expected body format:</p>

```
{
    "role": "student",
    "email": "example@exmail.ext",
    "password": "password"
}
```

<p>Valid roles would include one of {"student", "teacher", "admin"}</p>

<p>Expected response body for a successful request:</p>

```
{
    "status":  "success",
	"message": "session created",
	"jwt":     jwtToken,
	"token":   refreshToken,
}
```

<p>Expected response body for failed request:</p>

```
{
    "status" : "failed"
    "message": <error_description>
}
```

<h3>Refresh Token Endpoint 🔄</h3>
<p>This endpoint exists to generate a new JWT for the client</p>
<p>Proper implementation client-side should involve setting a timer based on the expiration date of the JWT</p>
<p>Once the JWT expires send a request to the server for re-generation</p>
<p>If the refresh token has expired the request will fail to generate a new JWT</p>
<p>Make a POST request to http://localhost:3000/jwtToken</p>
<p>Ensure that the Content-Type is "application/json"</p>
<p>The expected body format:</p>

```
{
    "email": "example@exmail.ext",
    "token": <refresh_token>
}
```

<p>Expected response body for a successful request:</p>

```
{
    "status":  "success",
	"message": "session continued",
	"jwt":     jwtToken,
	"token":   refreshToken,
}
```

<p>Expected response body for failed request:</p>

```
{
    "status" : "failed"
    "message": <error_description>
}
```
