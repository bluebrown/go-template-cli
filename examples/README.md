# tpl Examples

## Create a Certificate

```bash
tpl -t examples/cert.tpl
```

## Table

```bash
curl -fsSL https://jsonplaceholder.typicode.com/todos | tpl '{{ table . }}'
```
