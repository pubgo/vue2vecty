package vue2vecty

import (
	"fmt"
	"github.com/aymerick/douceur/parser"
	"github.com/dave/jennifer/jen"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/vue2vecty/xml"
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

var ternaryReg = xerror.PanicErr(regexp.Compile(`(.+)\?(.+):(.+)`)).(*regexp.Regexp)
var ternaryBrace = xerror.PanicErr(regexp.Compile(`.*{{{(.+)}}}.*`)).(*regexp.Regexp)
var twoBrace = xerror.PanicErr(regexp.Compile(`.*{{(.+)}}.*`)).(*regexp.Regexp)
var forReg = xerror.PanicErr(regexp.Compile(`,|\s+in\s+`)).(*regexp.Regexp)
var _for = func(child *jen.Statement, value string) *jen.Statement {
	xerror.PanicT(value == "", "params is zero")

	var params []string
	for _, p := range forReg.Split(value, -1) {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		params = append(params, p)
	}

	var s *jen.Statement
	if len(params) == 1 {
		s = jen.Id("key,value")
	}

	if len(params) == 2 {
		s = jen.Id(params[0])
	}

	if len(params) == 3 {
		s = jen.Id(params[0] + "," + params[1])
	}

	xerror.PanicT(s == nil, "statements error")

	return jen.Func().Params().Params(jen.Id("e").Qual("github.com/gopherjs/vecty", "List")).BlockFunc(func(g *jen.Group) {
		g.For(s.Op(":=").Range().Add(exp(params[len(params)-1]))).BlockFunc(func(f *jen.Group) {
			f.Id("e=").Append(jen.Id("e"), child)
		})
		g.Return()
	}).Call()
}

func exp(e string) *jen.Statement {
	e = strings.TrimSpace(e)

	switch {
	case e == "":
		return nil
	case strings.Contains(e, "$"):
		return exp(strings.ReplaceAll(e, "$", "t."))
	case !strings.Contains(e, " ") && strings.Contains(e, "="):
		return jen.Id(strings.ReplaceAll(strings.Title(e), "-", ""))
	case len(e) > 2 && e[0:1] == "'" && e[len(e)-1:] == "'":
		return jen.Id(fmt.Sprintf(`"%s"`, e[1:len(e)-1]))
	case strings.Contains(e, "{{{") && strings.Contains(e, "}}}") && ternaryBrace.MatchString(e):
		_d := ternaryBrace.FindStringSubmatch(e)
		if len(_d) == 1 {
			return nil
		}
		return jen.Qual("github.com/gopherjs/vecty", "UnsafeHTML").Call(exp(_d[1]))
	case twoBrace.MatchString(e):
		_d := twoBrace.FindStringSubmatch(e)
		if len(_d) == 1 {
			return nil
		}
		return jen.Qual("github.com/gopherjs/vecty", "Text").Call(exp(_d[1]))
	case len(e) > 2 && e[:1] == "[" && e[len(e)-1:] == "]":
		return nil
	case len(e) > 2 && e[:1] == "{" && e[len(e)-1:] == "}":
		return nil
	case len(strings.Split(e, "?:")) == 2:
		_s := strings.Split(e, "?:")
		s0 := strings.TrimSpace(_s[0])
		s1 := strings.TrimSpace(_s[1])
		return jen.Func().Interface().BlockFunc(func(g *jen.Group) {
			g.If(jen.Id(s0).Op("==").Lit("")).BlockFunc(func(g *jen.Group) {
				g.Return(jen.Id(s1))
			}).Else().BlockFunc(func(g *jen.Group) {
				g.Return(jen.Id(s0))
			})
		}).Call()
	case ternaryReg.MatchString(e):
		_s := ternaryReg.FindStringSubmatch(e)
		s1 := strings.TrimSpace(_s[1])
		s2 := strings.TrimSpace(_s[2])
		s3 := strings.TrimSpace(_s[3])
		return jen.Func().Params().Interface().BlockFunc(func(g *jen.Group) {
			g.If(jen.Id(s1)).BlockFunc(func(g *jen.Group) {
				g.Return(jen.Id(s2))
			}).Else().BlockFunc(func(g *jen.Group) {
				g.Return(jen.Id(s3))
			})
		}).Call()
	default:
		return jen.Id(e)
	}
}

func componentAttr(file *jen.File, d jen.Dict, g *jen.Group, attr xml.Attr) {
	ns := strings.TrimSpace(attr.Name.Space)
	key := strings.TrimSpace(attr.Name.Local)
	key = strings.TrimLeft(key, ns)
	value := strings.ReplaceAll(strings.TrimSpace(attr.Value), "$", "t.")

	fmt.Println(ns, key, value)

	switch {
	case ns == "" && (key == "v-for" || key == "xmlns" || key == "v-if"):
		return
	case ns == "" && len(key) > 5 && key[:5] == "data-" && value != "":
		key = key[5:]
		if d != nil {
			d[jen.Lit(strings.Title(key))] = exp(value)
		}

		if g != nil {
			g.Qual("github.com/gopherjs/vecty", "Data").Call(jen.Lit(key), exp(value))
		}
		return
	case ns == "v-on" || key[0] == '@':
		key = strings.ReplaceAll(strings.Title(key[1:]), "-", "")
		//fmt.Println(g != nil, d != nil, ns, key, len(key), key[0] == '@',value, "\n\n")
		if d != nil {
			d[jen.Lit("On"+key)] = jen.Func().Block(exp(value)).Call()
		}

		if g != nil {
			g.Qual("github.com/gopherjs/vecty/event", key).Call(exp(value))
		}
		return
	case ns == "v-bind" || key[0] == ':':
		key = key[1:]
		switch key {
		case "autofocus", "alt", "checked", "htmlFor", "href", "id", "placeholder", "src", "type", "value", "name", "disabled":
			if key == "id" {
				key = "ID"
			} else if key == "htmlFor" {
				key = "For"
			}

			g.Qual(vectyPropPackage, strings.Title(key)).Call(exp(value))
			return
		case "style":
			//style(file, g, value)
			return
		case "v-html":
			return
		case "class":
			return
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
			return
		}
		return
	case key == "v-model":
		_exp := exp(value)
		if _exp == nil {
			return
		}

		if d != nil {
			d[jen.Id("Value")] = _exp
			d[jen.Id("OnInput")] = jen.Func().Params(jen.Id("v").String()).BlockFunc(func(g *jen.Group) {
				g.Id(value).Op("=").Id("v")
			})
		}

		if g != nil {
			file.ImportName("github.com/gopherjs/vecty/prop", "prop")
			file.ImportName("github.com/gopherjs/vecty/event", "event")
			file.ImportName("honnef.co/go/js/dom/v2", "dom")

			g.Qual("github.com/gopherjs/vecty/prop", "Value").Call(_exp)
			g.Qual("github.com/gopherjs/vecty/event", "Input").CallFunc(func(g *jen.Group) {
				g.Func().Params(jen.Id("e").Op("*").Qual("github.com/gopherjs/vecty", "Event")).BlockFunc(func(g *jen.Group) {
					// t.Value = dom.WrapEvent(event.Target).Target().TextContent()
					// dom.WrapEvent(event.Target).PreventDefault()
					g.Id(value).Op("=").Id("dom.WrapEvent(e.Target).Target().TextContent()")
					g.Id("dom.WrapEvent(e.Target).PreventDefault()")
				})
			})
		}
	case key == "v-html":
		g.Qual("github.com/gopherjs/vecty", "UnsafeHTML").Call(exp(value))
		return
	case key == "v-text":
		return
	case key == "v-focus":
		return
	case key == "v-pre":
		return
	case key == "v-once":
		return
	case key == "class":
		return
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
		return
	default:
		fmt.Println(ns, key, value, "ok")
		return
	}

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

func componentElement(file *jen.File, appPackage string, token xml.StartElement) *jen.Statement {
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
				componentAttr(file, d, nil, v)
			}
		}))

		for _, v := range token.Attr {
			componentAttr(file, nil, g, v)
		}
	})
}
