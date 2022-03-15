{{- $serverTLS := "cool.com" -}}
{{- $ca := genCA (printf "%s-ca" $serverTLS) 365 -}}
{{- $cert := genSignedCert $serverTLS nil (list $serverTLS) 365 $ca -}}
---
data:
  ca.crt: {{ $ca.Cert | b64enc | quote }}
  tls.crt: {{ $cert.Key | b64enc | quote }}
  tls.key: {{ $cert.Cert | b64enc | quote }}
