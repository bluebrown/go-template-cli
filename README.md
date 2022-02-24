# tpl

Render yaml and json with go templates from the command line.

The templates are executed with the [text/template](https://pkg.go.dev/text/template) package. This means they come with the additional risks and benefits the text templates provide. Additionally, [Sprig functions](http://masterminds.github.io/sprig/) are available.

## Usage

The input data is read from stdin via pipe or redirection.

```bash
# Redirection
tpl < path/to/input.json
# Pipe
curl localhost | tpl
```

The template is either read from the first positional argument or from a path specified via `--template` or `-t` flag.

```bash
# Positional argument
echo '{"place": "bar"}' | tpl 'lets go to the {{.place}}!'
# File
echo '{"place": "bar"}' | tpl --template path/to/template
```

### Flags

```console
Usage of tpl:
  -t string
  -template string
        alternative way to specify template
  -n
  -no-newline
        do not print a new line at the end
  -h
  -help
        show this message
```

## Installation

## Binary

Download the binary from the [release page](https://github.com/bluebrown/tpl/releases). For example

```bash
curl -fsSLO https://github.com/bluebrown/tpl/releases/download/v1.0.0/tpl-amd64-static.tar.gz
tar -xzf tpl-amd64-static.tar.gz
mv tpl-1.0.0-amd64-static /usr/local/bin/tpl
```

## Go

If you have go installed, you can use the `go install` command to install the binary.

```bash
go install github.com/bluebrown/tpl
```

### Docker

The binary is also available as a docker image.

```shell
curl -s https://jsonplaceholder.typicode.com/users/1 | docker run -i bluebrown/tpl '{{ .name }}'
```

## Examples

## Convert YAML to JSON

```bash
echo 'foo: [bar, baz]' | tpl '{{ toPrettyJson . }}'
```

<details>
<summary>Output</summary>

```json
{
  "foo": [
    "bar",
    "baz"
  ]
}
```

</details>

## Render HTML from JSON

```bash
curl -s https://jsonplaceholder.typicode.com/users | tpl '<table>
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
        <li><strong>{{$key}}:</strong> &nbsp; {{$val}}</li>
        {{- end -}}
        {{ end }}
      </ul>
    </td>
  </tr>
  {{- end -}}
</table>'
```

<details>
<summary>Output</summary>

<table>
  <caption>My Address Nook</caption>
  <tr>
    <th>Name</th>
    <th>Email</th>
    <th>Phone</th>
    <th>Address</th>
  </tr>
  <tr>
    <th>Leanne Graham</th>
    <td>Sincere@april.biz</td>
    <td>1-770-736-8031 x56442</td>
    <td>
      <ul>
        <li><strong>city:</strong> &nbsp; Gwenborough</li>
        <li><strong>street:</strong> &nbsp; Kulas Light</li>
        <li><strong>suite:</strong> &nbsp; Apt. 556</li>
        <li><strong>zipcode:</strong> &nbsp; 92998-3874</li>
      </ul>
    </td>
  </tr>
  <tr>
    <th>Ervin Howell</th>
    <td>Shanna@melissa.tv</td>
    <td>010-692-6593 x09125</td>
    <td>
      <ul>
        <li><strong>city:</strong> &nbsp; Wisokyburgh</li>
        <li><strong>street:</strong> &nbsp; Victor Plains</li>
        <li><strong>suite:</strong> &nbsp; Suite 879</li>
        <li><strong>zipcode:</strong> &nbsp; 90566-7771</li>
      </ul>
    </td>
  </tr>
  <tr>
    <th>Clementine Bauch</th>
    <td>Nathan@yesenia.net</td>
    <td>1-463-123-4447</td>
    <td>
      <ul>
        <li><strong>city:</strong> &nbsp; McKenziehaven</li>
        <li><strong>street:</strong> &nbsp; Douglas Extension</li>
        <li><strong>suite:</strong> &nbsp; Suite 847</li>
        <li><strong>zipcode:</strong> &nbsp; 59590-4157</li>
      </ul>
    </td>
  </tr>
  <tr>
    <th>Patricia Lebsack</th>
    <td>Julianne.OConner@kory.org</td>
    <td>493-170-9623 x156</td>
    <td>
      <ul>
        <li><strong>city:</strong> &nbsp; South Elvis</li>
        <li><strong>street:</strong> &nbsp; Hoeger Mall</li>
        <li><strong>suite:</strong> &nbsp; Apt. 692</li>
        <li><strong>zipcode:</strong> &nbsp; 53919-4257</li>
      </ul>
    </td>
  </tr>
  <tr>
    <th>Chelsey Dietrich</th>
    <td>Lucio_Hettinger@annie.ca</td>
    <td>(254)954-1289</td>
    <td>
      <ul>
        <li><strong>city:</strong> &nbsp; Roscoeview</li>
        <li><strong>street:</strong> &nbsp; Skiles Walks</li>
        <li><strong>suite:</strong> &nbsp; Suite 351</li>
        <li><strong>zipcode:</strong> &nbsp; 33263</li>
      </ul>
    </td>
  </tr>
  <tr>
    <th>Mrs. Dennis Schulist</th>
    <td>Karley_Dach@jasper.info</td>
    <td>1-477-935-8478 x6430</td>
    <td>
      <ul>
        <li><strong>city:</strong> &nbsp; South Christy</li>
        <li><strong>street:</strong> &nbsp; Norberto Crossing</li>
        <li><strong>suite:</strong> &nbsp; Apt. 950</li>
        <li><strong>zipcode:</strong> &nbsp; 23505-1337</li>
      </ul>
    </td>
  </tr>
  <tr>
    <th>Kurtis Weissnat</th>
    <td>Telly.Hoeger@billy.biz</td>
    <td>210.067.6132</td>
    <td>
      <ul>
        <li><strong>city:</strong> &nbsp; Howemouth</li>
        <li><strong>street:</strong> &nbsp; Rex Trail</li>
        <li><strong>suite:</strong> &nbsp; Suite 280</li>
        <li><strong>zipcode:</strong> &nbsp; 58804-1099</li>
      </ul>
    </td>
  </tr>
  <tr>
    <th>Nicholas Runolfsdottir V</th>
    <td>Sherwood@rosamond.me</td>
    <td>586.493.6943 x140</td>
    <td>
      <ul>
        <li><strong>city:</strong> &nbsp; Aliyaview</li>
        <li><strong>street:</strong> &nbsp; Ellsworth Summit</li>
        <li><strong>suite:</strong> &nbsp; Suite 729</li>
        <li><strong>zipcode:</strong> &nbsp; 45169</li>
      </ul>
    </td>
  </tr>
  <tr>
    <th>Glenna Reichert</th>
    <td>Chaim_McDermott@dana.io</td>
    <td>(775)976-6794 x41206</td>
    <td>
      <ul>
        <li><strong>city:</strong> &nbsp; Bartholomebury</li>
        <li><strong>street:</strong> &nbsp; Dayna Park</li>
        <li><strong>suite:</strong> &nbsp; Suite 449</li>
        <li><strong>zipcode:</strong> &nbsp; 76495-3109</li>
      </ul>
    </td>
  </tr>
  <tr>
    <th>Clementina DuBuque</th>
    <td>Rey.Padberg@karina.biz</td>
    <td>024-648-3804</td>
    <td>
      <ul>
        <li><strong>city:</strong> &nbsp; Lebsackbury</li>
        <li><strong>street:</strong> &nbsp; Kattie Turnpike</li>
        <li><strong>suite:</strong> &nbsp; Suite 198</li>
        <li><strong>zipcode:</strong> &nbsp; 31428-2261</li>
      </ul>
    </td>
  </tr></table>

</details>
