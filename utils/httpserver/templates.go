package httpserver

import "html/template"

var templates = template.Must(template.New("httpserver").Parse(`
{{define "errors/generic.html"}}<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>{{.title}}</title>
</head>
<body>
<h1>{{.status}} {{.title}}</h1>
<p>{{.message}}</p>
</body>
</html>{{end}}

{{define "errors/validation.html"}}<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>{{.Message}}</title>
</head>
<body>
<h1>{{.Status}} {{.Message}}</h1>
{{if .Errors}}
<ul>
{{range $field, $err := .Errors}}
<li><strong>{{$field}}</strong>: {{$err}}</li>
{{end}}
</ul>
{{end}}
</body>
</html>{{end}}

{{define "form/success.html"}}<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Success</title>
</head>
<body>
<h1>Form submitted</h1>
</body>
</html>{{end}}

{{define "form/redirect.html"}}<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Redirecting</title>
</head>
<body>
<p>Redirecting…</p>
</body>
</html>{{end}}
`))
