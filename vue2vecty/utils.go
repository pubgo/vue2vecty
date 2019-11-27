package vue2vecty

import (
	"encoding/xml"
	"github.com/aymerick/douceur/parser"
	"github.com/dave/jennifer/jen"
	"github.com/pubgo/g/xerror"
	"regexp"
	"strings"
)

func createStruct(packageName, componentName string) *jen.File {
	file := jen.NewFile(packageName)
	file.PackageComment("This file was created with https://github.com/pubgo/vue2vecty")
	file.ImportName("github.com/gopherjs/vecty", "vecty")

	componentName = strings.ReplaceAll(strings.Title(componentName), "-", "")
	file.Type().Id(componentName).Struct(
		jen.Qual("github.com/gopherjs/vecty", "Core"),
	)
	return file
}

func exp(e string) *jen.Statement {
	e = strings.TrimSpace(e)
	_, _ = regexp.MatchString(`.+\?.+:.+`, "")
	/*
		[activeClass, errorClass]
		["activeClass", "errorClass"]
		[isActive ? activeClass : "", errorClass]
		[{ active: isActive }, errorClass]
		{ color: activeColor, fontSize: fontSize + 'px' }
		{ display: ['-webkit-box', '-ms-flexbox', 'flex'] }
		{{ item.message }}
		counter += 1
	*/

	m, err := regexp.MatchString(`.+\?.+:.+`, "")

	switch {
	case len(e) > 4 && e[:2] == "{{" && e[len(e)-2:] == "}}":
	case len(e) > 2 && e[:1] == "[" && e[len(e)-1:] == "]":
	case len(e) > 2 && e[:1] == "{" && e[len(e)-1:] == "}":
	case !strings.Contains(e, " "):
		return jen.Id(e)

	case strings.Contains(e, "?:") && len(strings.Split(e, "?:")) == 2: //{{ name?:"hello" }}
		return jen.If(jen.Id("name").Op("==").Lit("")).BlockFunc(func(g *jen.Group) {
			g.Return()
		}).Else().BlockFunc(func(g *jen.Group) {
			g.Return()
		})
	case strings.Contains(e, "?:") && len(strings.Split(e, "?:")) == 2: //{{ name==""? "hello" : "world" }}
		return jen.If(jen.Id("name").Op("==").Lit("")).BlockFunc(func(g *jen.Group) {
			g.Return()
		}).Else().BlockFunc(func(g *jen.Group) {
			g.Return()
		})
	case m && err == nil:
	default:
		return jen.Id(e)
	}

	return nil
}

func componentAttr(custom bool, file *jen.File, d jen.Dict, g jen.Group, attr xml.Attr) (*jen.Statement, *jen.Statement) {
	ns := strings.TrimSpace(attr.Name.Space)
	key := strings.TrimSpace(attr.Name.Local)
	key = strings.TrimLeft(key, ns)
	value := strings.TrimSpace(attr.Value)

	switch {
	case ns == "":
		switch key {
		case "v-for", "xmlns":
			return nil, nil
		case "v-if":
			file.ImportName("github.com/gopherjs/vecty", "vecty")
			return nil, nil
		}
	case ns == "v-on" || key[0] == '@':
		key = strings.ReplaceAll(strings.Title(strings.TrimLeft(key, ns)[1:]), "-", "")
		if custom {
			d[jen.Lit("On"+key)] = jen.Id(value)
			return nil, nil
		} else {
			g.Qual("github.com/gopherjs/vecty/event", key).Call(jen.Id(value))
			return nil, nil
		}
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
		value = "t." + value
		if custom {
			d[jen.Id("Value")] = jen.Id(value)
			d[jen.Id("OnInput")] = jen.Func().Params(jen.Id("v").String()).BlockFunc(func(g *jen.Group) {
				g.Id(value).Op("=").Id("v")
			})
		} else {
			file.ImportName("github.com/gopherjs/vecty/prop", "prop")
			file.ImportName("github.com/gopherjs/vecty/event", "event")

			g.Qual("github.com/gopherjs/vecty/prop", "Value").Call(jen.Id(value))
			g.Qual("github.com/gopherjs/vecty/event", "Input").CallFunc(func(g *jen.Group) {
				g.Func().Params(jen.Id("e").Qual("github.com/gopherjs/vecty", "Event")).BlockFunc(func(g *jen.Group) {
					// t.Value = dom.WrapEvent(event.Target).Target().TextContent()
					// dom.WrapEvent(event.Target).PreventDefault()
					g.Id(value).Op("=").Id("e")
				})
			})
		}

	case key == "v-html":
		return nil, nil
	case key == "style":
		css, err := parser.ParseDeclarations(value)
		xerror.PanicM(err, "css parsing error")

		for _, dec := range css {
			if dec.Important {
				dec.Value += "!important"
			}
			file.ImportName("github.com/gopherjs/vecty", "vecty")
			jen.Qual("github.com/gopherjs/vecty", "Style").Call(
				jen.Lit(dec.Property),
				jen.Lit(dec.Value),
			)
		}
	case key == "class":

	default:
		key = strings.ToUpper(string(key[0])) + key[1:]
		d[jen.Id(key)] = jen.Lit(value)
	}

	return nil, nil
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
	case key == "class":
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

	case strings.HasPrefix(key, "data-"):
		attribute := strings.TrimPrefix(key, "data-")
		g.Qual("github.com/gopherjs/vecty", "Data").Call(
			jen.Lit(attribute),
			jen.Lit(value),
		)

	case boolProps[key] != "":
		value := value == "true"
		g.Qual("github.com/gopherjs/vecty/prop", boolProps[key]).Call(
			jen.Lit(value),
		)
	case stringProps[key] != "":
		if strings.HasPrefix(value, "{vecty-field:") {
			field := strings.TrimLeft(value, "{vecty-field:")
			field = field[:len(field)-1]
			g.Qual("github.com/gopherjs/vecty/prop", stringProps[key]).Call(
				jen.Id("p").Dot(field),
			)
		} else {
			g.Qual("github.com/gopherjs/vecty/prop", stringProps[key]).Call(
				jen.Lit(value),
			)
		}
	case strings.HasPrefix(ns, "vecty"):
		field := strings.TrimLeft(key, "on")
		field = strings.ToLower(field)
		g.Qual("github.com/gopherjs/vecty/event", strings.Title(field)).Call(
			jen.Id("p").Dot(value),
		)
	case strings.HasPrefix(ns, "components"):
		component := strings.TrimLeft(key, "components.")
		jen.Op("&").Id(component).Values()
	case key == "xmlns":
		g.Qual("github.com/gopherjs/vecty", "Namespace").Call(
			jen.Lit(value),
		)
	case key == "type" && typeProps[value] != "":
		g.Qual("github.com/gopherjs/vecty/prop", "Type").Call(
			jen.Qual("github.com/gopherjs/vecty/prop", typeProps[value]),
		)

	case key == "v-for":

	default:
		g.Qual("github.com/gopherjs/vecty", "Attribute").Call(
			jen.Lit(key),
			jen.Lit(value),
		)
	}
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

	name = strings.ReplaceAll(strings.Title(name), "-", "")
	return jen.Qual(vectyPackage, name).CallFunc(func(g *jen.Group) {
		g.Map(jen.String()).Interface().Values(jen.DictFunc(func(d jen.Dict) {
			for _, v := range token.Attr {
				if v.Name.Local == "class" {
					continue
				}

				d[jen.Id(v.Name.Local)] = jen.Lit(v.Value)
			}
		}))

		for _, v := range token.Attr {
			if v.Name.Local == "class" {
				g.Qual("github.com/gopherjs/vecty", "Style").CallFunc(func(g *jen.Group) {
					jen.Lit(v.Value)
					jen.Lit(v.Value)
				})
			}
		}
	})
}
