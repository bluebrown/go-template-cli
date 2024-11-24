package textfunc

import (
	"text/template"
)

func Register(baseMap map[string]any) map[string]any {
	baseMap["table"] = Table
	baseMap["fromJson"] = FromJson
	baseMap["toYaml"] = ToYaml
	baseMap["mustToYaml"] = MustToYaml
	baseMap["fromYaml"] = FromYaml
	baseMap["toToml"] = ToToml
	baseMap["mustToToml"] = MustToToml
	baseMap["fromToml"] = FromToml
	baseMap["iter"] = Iter
	return baseMap
}

func Map() map[string]any {
	funcMap := make(map[string]any)
	return Register(funcMap)

}

func MapClosure(baseMap map[string]any, t *template.Template) map[string]any {
	funcMap := Register(baseMap)
	funcMap["include"] = MakeInclude(t)
	funcMap["tpl"] = MakeTpl(t)
	return funcMap
}
