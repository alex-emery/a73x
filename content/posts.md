---
title: posts
nav: /posts
---
{{range .Collections.posts}}
- [{{.Meta.title}}]({{.Path}})
{{end}}