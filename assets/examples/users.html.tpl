<table>
  <caption>My Address Nook</caption>
  <tr>
    <th>Name</th>
    <th>Email</th>
    <th>Phone</th>
    <th>Address</th>
  </tr>
  {{- range . }}
  <tr>
    <th>{{ .name }}</th>
    <td>{{ .email }}</td>
    <td>{{ .phone }}</td>
    <td>
      <ul>
        {{- range $key, $val := .address }} {{ if ne $key "geo" }}
        <li><strong>{{ $key }}:</strong> &nbsp; {{ $val }}</li>
        {{- end -}}
        {{ end }}
      </ul>
    </td>
  </tr>
  {{- end -}}
</table>
