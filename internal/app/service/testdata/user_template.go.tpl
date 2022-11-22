{{- /*gotype: github.com/spatecon/gitlab-review-bot/internal/app/ds.UserNotification*/ -}}
Привет, {{.User.BasicUser.Name}}!

{{ if len .ReviewerMR -}}
:crossed_fingers:  *Ты ревьювер {{len .ReviewerMR}} {{plural (len .ReviewerMR) "реквеста" "реквестов" "реквестов"}}:*
{{ range .ReviewerMR }}
*{{.Title}}*
{{.URL}} менялся *{{.UpdatedAt | since}}*
{{ end -}}
{{- end }}

{{ if .AuthoredMR -}}
:index_pointing_at_the_viewer: *Ты автор {{len .AuthoredMR}} {{plural (len .AuthoredMR) "реквеста" "реквестов" "реквестов"}} в ревью:*
{{ range .AuthoredMR }}
*{{.Title}}*
{{.URL}} менялся *{{.UpdatedAt | since}}*
{{ end -}}
{{- end -}}