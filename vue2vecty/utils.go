package vue2vecty

import (
	"encoding/xml"
	"github.com/dave/jennifer/jen"
	"strings"
)

func createStruct(packageName, componentName string) *jen.File {
	file := jen.NewFile(packageName)
	file.PackageComment("This file was created with https://github.com/pubgo/vue2vecty")
	file.ImportName("github.com/gopherjs/vecty", "vecty")

	var cs string
	for _, s := range strings.Split(componentName, "-") {
		cs += strings.ToUpper(string(s[0])) + s[1:]
	}

	file.Type().Id(cs).Struct(
		jen.Qual("github.com/gopherjs/vecty", "Core"),
	)

	return file
}

func componentElement(file *jen.File, appPackage string, token *xml.StartElement) *jen.Statement {
	appPackage += "/components"
	ts := strings.Split(token.Name.Local, ":")
	_l := len(ts)
	name := ""
	if _l == 1 {
		name = ts[0]
		file.ImportName(appPackage, "components")
	} else if _l >= 2 {
		name = ts[_l-1]
		file.ImportName(appPackage+"/"+strings.Join(ts[:_l-1], "/"), ts[_l-2])
	}

	var cn string
	for _, s := range strings.Split(name, "-") {
		cn += strings.ToUpper(string(s[0])) + s[1:]
	}

	return jen.Op("&").Qual(vectyPackage, cn).Values(jen.DictFunc(func(d jen.Dict) {
		for _, v := range token.Attr {
			d[jen.Id(v.Name.Local)] = jen.Lit(v.Value)
		}
	}))
}
