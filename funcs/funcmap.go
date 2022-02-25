package funcs

func TxtFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"toYaml":     ToYaml,
		"mustToYaml": MustToYaml,
		"table":      Table,
	}
}
