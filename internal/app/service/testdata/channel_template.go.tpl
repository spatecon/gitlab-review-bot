{{- /*gotype: github.com/spatecon/gitlab-review-bot/internal/app/ds.ChannelNotification*/ -}}
:colleagues-kollegi: Статистика по ревью :colleagues-kollegi:
В среднем: {{.AverageCount}} MR на разработчике
Всего: {{.TotalCount}} уникальных MR
{{- if .FirstEditedMR }}
Самый старый MR менялся: {{.FirstEditedMR.UpdatedAt | since}}
{{- end -}}