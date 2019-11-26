package vue2vecty

import (
	"encoding/xml"
	"github.com/aymerick/douceur/parser"
	"github.com/dave/jennifer/jen"
	"github.com/pubgo/g/xerror"
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

func componentAttr(file *jen.File, d jen.Dict, attr xml.Attr) {
	ns := attr.Name.Space
	key := attr.Name.Local
	value := attr.Value

	switch {
	case key == "v-if" && ns == "":
		file.ImportName("github.com/gopherjs/vecty", "If")
		return
	case (key == "v-for" || key == "style" || key == "class" || key == "v-on" || key == "xmlns") && ns == "":
		return
	case key == "v-model":
		return
	case ns == "v-bind" || key[0] == ':':
		key = string(strings.TrimLeft(key, ns)[1:])
		key = strings.ToUpper(string(key[0])) + key[1:]
		d[jen.Id(key)] = jen.Id("t." + value)
	default:
		key = strings.ToUpper(string(key[0])) + key[1:]
		d[jen.Id(key)] = jen.Lit(value)
	}

	return
}

func style(file *jen.File, g *jen.Group, value string) {
	css, err := parser.ParseDeclarations(value)
	xerror.PanicM(err, "css parsing error")

	for _, dec := range css {
		if dec.Important {
			dec.Value += "!important"
		}
		file.ImportName("github.com/gopherjs/vecty", "vecty")
		g.Qual("github.com/gopherjs/vecty", "Style").Call(
			jen.Lit(dec.Property),
			jen.Lit(dec.Value),
		)
	}
}

func tagAttr(file *jen.File, g *jen.Group, attr xml.Attr) {
	ns := attr.Name.Space
	key := attr.Name.Local
	value := attr.Value

	switch {
	case ns == "" && key == "v-if":
		file.ImportName("github.com/gopherjs/vecty", "vecty")
		return
	case ns == "" && (key == "v-for" || key == "xmlns"):
		return
	case ns == "v-bind" || key[0] == ':':
		key = string(strings.TrimLeft(key, ns)[1:])
		switch key {
		case "style":
			style(file, g, value)
		case "v-html":
		case "class":
			file.ImportName("github.com/gopherjs/vecty", "vecty")
			g.Qual("github.com/gopherjs/vecty", "Class").CallFunc(func(g *jen.Group) {
				classes := strings.Split(value, " ")
				for _, class := range classes {
					if strings.HasPrefix(class, "{vecty-field:") {
						field := strings.TrimLeft(class, "{vecty-field:")
						field = field[:len(field)-1]
						g.Add(jen.Id("p").Dot(field))
					} else {
						g.Lit(class)
					}
				}
			})
		}
	case key == "v-model":
	case key == "v-on":
	case key == "style":
		style(file, g, value)
	case v.Name.Local == "class":
		g.Qual("github.com/gopherjs/vecty", "Class").CallFunc(func(g *jen.Group) {
			classes := strings.Split(v.Value, " ")
			for _, class := range classes {
				if strings.HasPrefix(class, "{vecty-field:") {
					field := strings.TrimLeft(class, "{vecty-field:")
					field = field[:len(field)-1]
					g.Add(jen.Id("p").Dot(field))
				} else {
					g.Lit(class)
				}
			}
		})

	case strings.HasPrefix(v.Name.Local, "data-"):
		attribute := strings.TrimPrefix(v.Name.Local, "data-")
		g.Qual("github.com/gopherjs/vecty", "Data").Call(
			jen.Lit(attribute),
			jen.Lit(v.Value),
		)

	case boolProps[v.Name.Local] != "":
		value := v.Value == "true"
		g.Qual("github.com/gopherjs/vecty/prop", boolProps[v.Name.Local]).Call(
			jen.Lit(value),
		)
	case stringProps[v.Name.Local] != "":
		if strings.HasPrefix(v.Value, "{vecty-field:") {
			field := strings.TrimLeft(v.Value, "{vecty-field:")
			field = field[:len(field)-1]
			g.Qual("github.com/gopherjs/vecty/prop", stringProps[v.Name.Local]).Call(
				jen.Id("p").Dot(field),
			)
		} else {
			g.Qual("github.com/gopherjs/vecty/prop", stringProps[v.Name.Local]).Call(
				jen.Lit(v.Value),
			)
		}
	case strings.HasPrefix(v.Name.Space, "vecty"):
		field := strings.TrimLeft(v.Name.Local, "on")
		field = strings.ToLower(field)
		g.Qual("github.com/gopherjs/vecty/event", strings.Title(field)).Call(
			jen.Id("p").Dot(v.Value),
		)
	case strings.HasPrefix(v.Name.Space, "components"):
		component := strings.TrimLeft(v.Name.Local, "components.")
		jen.Op("&").Id(component).Values()
	case v.Name.Local == "xmlns":
		g.Qual("github.com/gopherjs/vecty", "Namespace").Call(
			jen.Lit(v.Value),
		)
	case v.Name.Local == "type" && typeProps[v.Value] != "":
		g.Qual("github.com/gopherjs/vecty/prop", "Type").Call(
			jen.Qual("github.com/gopherjs/vecty/prop", typeProps[v.Value]),
		)

	case v.Name.Local == "v-for":

	default:
		g.Qual("github.com/gopherjs/vecty", "Attribute").Call(
			jen.Lit(v.Name.Local),
			jen.Lit(v.Value),
		)
	}
}

func componentElement(file *jen.File, appPackage string, token *xml.StartElement) *jen.Statement {
	* /
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
