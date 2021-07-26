{{ define "vendir" }}
apiVersion: vendir.k14s.io/v1alpha1
kind: Config

directories:
{{ range . }}
  - path: {{ .Path }}
    contents:
    {{ range .Content }}
      - path: {{ .Path }}
      {{ if .Git }}
        git:
          url: {{ .Git.Url }}
          ref: {{ .Git.Ref }}
      {{ end }}
      {{ if .Dir }}
        directory:
          path: {{ .Dir.Path }}
      {{ end }}
        newRootPath: {{ .NewRootPath }}
        includePaths:
          - {{ .IncludePaths }}
    {{ end }}
{{ end }}
{{ end }}