<h1>Docs for Event EndPoints 🗓️📍</h1>

<p>Note: requests to this endpoint require a <a href="https://github.com/golang-jwt/jwt">JWT</a> token because it is a protected resource</p>

<h3>GET EndPoint ⬇️</h3>
<p>This endpoint exists to allow clients requests all the events associated with them.</p>
<p>Make a GET request to http://localhost:3000/v1/events</p>
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
	"message":      "events were gotten",
	"events":       <events>,
}

```

<p>Expected response body for failed request:</p>

```
{
    "status" : "failed",
    "message": <error_description>,
}
```

<h3>POST EndPoint ⬆️</h3>
<p>This endpoint exists to allow clients with the "teacher" role create events.</p>
<p>Note: event id must be unique, this would eventually be changed to an auto-increment or similar in the future</p>
<p>Make a POST request to http://localhost:3000/v1/events</p>
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
    "id": <id>
	"course_id": <course_id>
	"title": <title>
	"description": <description>
	"start_time": <time stamp in "2025-07-23T14:00:00Z" format>
	"end_time": <time stamp in "2025-07-23T14:00:00Z" format>
	"location": <location>
	"event_type": x ∈ {'assignment_due' | 'exam'| 'lecture'| 'holiday'}
	"created_by": <teacher_id>
	"created_at": <time stamp in "2025-07-23T14:00:00Z" format>
}

```

<p>Expected response body for a successful request:</p>

```
{
    "status": "successful",
	"message": "event was saved",
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
<p>This endpoint exists to allow clients with the "teacher" role edit events.</p>
<p>Make a PATCH request to http://localhost:3000/v1/events</p>
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
    "event_id": <event_id>,
    "patch_field": <column_name>,
    "new_content": <data>
}

```

<p>Expected response body for a successful request:</p>

```
{
    "status":       "successful",
	"message":      "event was updated",
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
<p>This endpoint exists to allow clients with the "teacher" role delete events.</p>
<p>Make a DELETE request to http://localhost:3000/v1/events</p>
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
    "event_id": <event_id>,
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
