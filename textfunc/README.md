# Treasure Map

A collection of go template functions. Some functions are inspired by
[helm](https://helm.sh/), i.e. include and tpl.

## Usage

Some function require to close over the template itself. Therefore, the
MakeFuncMap function takes the template as argument and returns the
fully constructed map. In the below example
[sprigs](http://masterminds.github.io/sprig/) TxtFuncMap is passed as
base map.

```go
import (
    "testing"
    "text/template"

    "github.com/Masterminds/sprig/v3"
    "github.com/bluebrown/go-template-cli/textfunc"
)

func main() {
    tpl := template.New("")
    tpl = tpl.Funcs(textfunc.MapClosure(sprig.TxtFuncMap(), tpl))
}
```

## Functions

Function    | Description                                         | Returns       | Example
------------|-----------------------------------------------------|---------------|------------------------------
include     | render a previously associated template             | string        | `{{ include "my-helper" . }}`
tpl         | render a template snippet                           | string        | `{{ tpl "{{ .foo.bar }}" . }}`
table       | convert to table, list of object or list of lists   | string        | `{{ table . }}`
fromJson    | convert a json string to map or slice (any)         | interface{}   | `{{ fromJson . }}`
toYaml      | convert to yaml                                     | string        | `{{ toYaml . }}`
mustToYaml  | convert to yaml, errors if encoding fails           | string, error | `{{ mustToYaml . }}`
fromYaml    | convert from yaml                                   | interface{}   | `{{ fromYaml . }}`
toToml      | convert to toml                                     | string        | `{{ toToml . }}`
mustToToml  | convert to toml, errors if encoding fails           | string, error | `{{ mustToToml . }}`
fromToml    | convert from toml                                   | interface{}   | `{{ fromToml . }}`
iter        | get an iterator to use in range                     | []int         | `{{ range iter 5 }}{{.}}{{end}}`
