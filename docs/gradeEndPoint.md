<h1>Docs for Grade EndPoint ✅</h1>

<p>Note: requests to this endpoint require a <a href="https://github.com/golang-jwt/jwt">JWT</a> because it is a protected resource.</p>

<h3>Grade EndPoint ⬆️</h3>
<p>This endpoint exists for the sole reason of grading assessments.</p>
<p>Make a POST request to http://localhost:3000/v1/grade</p>
<p>Ensure that the Content-Type is "application/json"</p>
<p>The expected body format:</p>

```
{
    "assessment_id": <id>,
    "choices": [
        {
            "question_id": <id>,
            "choice_text": <answer>
        },
        {
            "question_id": <id>,
            "choice_text": <answer>
        },
        ...
    ]
}

```

<p>Expected response body for a successful request:</p>

```
{
    "status":  "success",
	"message": "user created",
	"message": "assessment has been graded",
	"grade_metadata": gradeMeta,
}

```

<p>Expected response body for failed request:</p>

```
{
    "status" : "failed"
    "message": <error_description>
}
```
