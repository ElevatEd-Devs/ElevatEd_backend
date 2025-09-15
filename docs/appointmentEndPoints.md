<h1>Docs for Appointment EndPoints ⏰</h1>

<p>Note: requests to this endpoint require a <a href="https://github.com/golang-jwt/jwt">JWT</a> token because it is a protected resource</p>

<h3>GET EndPoint ⬇️</h3>
<p>This endpoint exists to allow clients requests all the appointments associated with them.</p>
<p>Make a GET request to http://localhost:3000/v1/appointments</p>
<p>Ensure that the Content-Type is "application/json"</p>
<p>Also add an Authentication field to the header</p>
<p>Minimum header content:</p>

```
Content-Type: application/json
Authorization: Bearer <JWT>
```

<p>Expected response body for a successful request:</p>

```
{
    "status":       "successful",
	"message":      "appointments were gotten",
	"appointments": <appointments>,
}

```

<p>Expected response body for failed request:</p>

```
{
    "status" : "failed",
    "message": <error_description>,
    "appointments" : null
}
```

<h3>POST EndPoint ⬆️</h3>
<p>This endpoint exists to allow clients create appointments.</p>
<p>Note: appointment id must be unique, this would eventually be changed to an auto-increment or similar in the future</p>
<p>Make a POST request to http://localhost:3000/v1/appointments</p>
<p>Ensure that the Content-Type is "application/json"</p>
<p>Also add an Authentication field to the header</p>
<p>Minimum header content:</p>

```
Content-Type: application/json
Authorization: Bearer <JWT>
```

<p>The body should contain at minimum:</p>

```
{
    "id": <appointment_id>,
    "office_hours_id": <office_hours_id>,
    "title": <title>,
    "student_id": <student_id>,
    "teacher_id": <teacher_id>,
    "start_time": <time stamp in "2025-07-23T14:00:00Z" format>,
    "end_time": <time stamp in "2025-07-23T14:00:00Z" format>
}

```

<p>Expected response body for a successful request:</p>

```
{
    "status":       "successful",
	"message":      "appointment was saved",
}

```

<p>Expected response body for failed request:</p>

```
{
    "status" : "failed",
    "message": <error_description>,
}
```

<h3>PATCH EndPoint 🪡</h3>
<p>This endpoint exists to allow clients edit appointments.</p>
<p>Make a PATCH request to http://localhost:3000/v1/appointments</p>
<p>Ensure that the Content-Type is "application/json"</p>
<p>Also add an Authentication field to the header</p>
<p>Minimum header content:</p>

```
Content-Type: application/json
Authorization: Bearer <JWT>
```

<p>The body should contain at minimum:</p>

```
{
    "appointment_id": <appointment_id>,
    "patch_field": <column_name>,
    "new_content": <data>
}

```

<p>Expected response body for a successful request:</p>

```
{
    "status":       "successful",
	"message":      "appointment was updated",
}

```

<p>Expected response body for failed request:</p>

```
{
    "status" : "failed",
    "message": <error_description>,
}
```

<h3>DELETE EndPoint 🗑️</h3>
<p>This endpoint exists to allow clients delete appointments.</p>
<p>Make a DELETE request to http://localhost:3000/v1/appointments</p>
<p>Ensure that the Content-Type is "application/json"</p>
<p>Also add an Authentication field to the header</p>
<p>Minimum header content:</p>

```
Content-Type: application/json
Authorization: Bearer <JWT>
```

<p>The body should contain at minimum:</p>

```
{
    "appointment_id": <appointment_id>,
}

```

<p>Expected response body for a successful request:</p>

```
{
    "status":       "successful",
	"message":      "appointment was deleted",
}

```

<p>Expected response body for failed request:</p>

```
{
    "status" : "failed",
    "message": <error_description>,
}
```
