# tpl Examples

## JSON to HTML

```bash
curl -s https://jsonplaceholder.typicode.com/users | tpl -t examples/users.html.tpl
```

## JSON to Table

```bash
curl -s https://jsonplaceholder.typicode.com/todos | tpl '{{ table . }}'
```

## Convert YAML to JSON

```bash
echo 'foo: [bar, baz]' | tpl '{{ toPrettyJson . }}'
```

## Create a Certificate

```bash
tpl -t examples/cert.yaml.tpl
```
