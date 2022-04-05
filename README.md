# Go Template CLI (tpl)

Render yaml and json with go templates from the command line.

The templates are executed with the [text/template](https://pkg.go.dev/text/template) package. This means they come with the additional risks and benefits the text templates provide. Additionally, [Sprig functions](http://masterminds.github.io/sprig/) are available.

## Synopsis

```console
Options:
  -f, --file stringArray      template file path. Can be specified multiple times
  -g, --glob stringArray      template file glob. Can be specified multiple times
  -n, --name string           if specified, execute the template with the given name
      --options stringArray   options to pass to the template engine
  -d, --decoder string        decoder to use for input data. Supported values: json, yaml, toml, xml (default "json")
      --no-newline            do not print newline at the end of the output
  -h, --help                  show the help text
  -v, --version               show the version
      --debug                 enable debug mode
```

## Input Data

The input data is read from stdin via pipe or redirection. It is actually not required to provide any input data. If no input data is provided, the template is executed with nil data.

```bash
# Redirection
tpl '{{ . }}' < path/to/input.json
# Pipe
curl localhost | tpl '{{ . }}'

# nil data
tpl '{{ . }}'
```

## Templates

The root templates name is `_tpl.root` and positional arguments are parsed into this root template. That means while its possible to specify multiple arguments, they will overwrite each other unless they use the `define` keyword to define a named template that can be referenced later when executing the template. If a named templates is specified multiple times, the last one will override the previous ones.

Templates from the flags --file and --glob are parsed in the order they are specified. So the override rules of the text/template package apply. If a file with the same name is specified multiple times, the last one wins. Even if they are in different directories.

The behavior of the cli tries to stay consistent with the actual behavior of the go template engine.

By default the root template is executed if at least one positional argument has been provided. Otherwise the first parsed file name is used to to determine which named template to execute since the root templates body is empty. It is always possible to choose another template to execute by using the --name flag.

```bash
tpl '{{ . }}' --file foo.tpl --glob templates/*.tpl
```

The ability to parse multiple templates makes sense when defining helper snippets and other named templates to reference using the builtin `template` keyword or the custom `include` function which can be used in pipelines.

## Decoders

By default input data is decoded as json and passed to the template to execute. It is possible to use an alternative decoder. The supported decoders are:

- json
- yaml
- toml
- xml

While json could technically be decoded using the yaml decoder, this is not done by default for performance reasons.

## Functions

Next to the builtin functionality [Sprig functions](http://masterminds.github.io/sprig/) and [treasre-map functions](https://github.com/bluebrown/treasure-map) are available.

## Installation

### Binary

Download the binary from the [release page](https://github.com/bluebrown/tpl/releases). For example

```bash
curl -fsSLO https://github.com/bluebrown/go-template-cli/releases/download/v0.1.0/tpl-amd64-static.tar.gz
tar -xzf tpl-amd64-static.tar.gz
mv tpl-0.1.0-amd64-static /usr/local/bin/tpl
```

### Go

If you have go installed, you can use the `go install` command to install the binary.

```bash
go install github.com/bluebrown/tpl
```

### Docker

The binary is also available as a [docker image](https://hub.docker.com/repository/docker/bluebrown/tpl).

```shell
curl -s curl -s https://jsonplaceholder.typicode.com/todos | docker run -i bluebrown/tpl '{{ table . }}'
```

### From source

Clone the repo and use the makefile to build the binary. The make install command will move the binary to /usr/local/bin.

```bash
git clone https://github.com/bluebrown/tpl
cd tpl && make install
```

## Example

Review the [examples](https://github.com/bluebrown/go-template-cli/tree/main/assets/examples) directory, for more examples.

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
        <li><strong>{{ $key }}:</strong> &nbsp; {{ $val }}</li>
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
