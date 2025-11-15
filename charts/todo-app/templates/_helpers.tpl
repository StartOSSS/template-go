{{- define "todo-app.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{ .Values.fullnameOverride }}
{{- else -}}
{{- printf "%s-%s" .Release.Name "api" | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "todo-app.selectorLabels" -}}
app.kubernetes.io/name: todo-api
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
