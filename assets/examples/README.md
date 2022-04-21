# tpl Examples

## JSON to HTML

```bash
curl -s https://jsonplaceholder.typicode.com/users | tpl -f assets/examples/users.html.tpl
```

## JSON to Table

```bash
curl -s https://jsonplaceholder.typicode.com/todos | tpl '{{ table . }}'
```

## Convert YAML to JSON

```bash
echo 'foo: [bar, baz]' | tpl '{{ toPrettyJson . }}' -d yaml
```

## Create a Certificate

```bash
tpl -f assets/examples/cert.yaml.tpl
```
