{{ define "vendir" }}
apiVersion: vendir.k14s.io/v1alpha1
kind: Config

directories:
{{ range . }}
  - path: {{ .Path }}
    {{ range .Content }}
    contents:
      - path: {{ .Path }}
        git:
          url: {{ .Git.Url }}
          ref: {{ .Git.Ref }}
        newRootPath: {{ .NewRootPath }}
        includePaths:
          - {{ .IncludePaths }}
    {{ end }}
{{ end }}
{{ end }}