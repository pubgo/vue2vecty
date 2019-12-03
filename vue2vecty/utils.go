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

var trim = strings.TrimSpace

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
var ternaryBrace = xerror.PanicErr(regexp.Compile(`(.*){{{(.+)}}}(.*)`)).(*regexp.Regexp)
var twoBrace = xerror.PanicErr(regexp.Compile(`(.*){{(.+)}}(.*)`)).(*regexp.Regexp)
var forReg = xerror.PanicErr(regexp.Compile(`,|\s+in\s+`)).(*regexp.Regexp)
var forExp = func(file *jen.File, child *jen.Statement, value string) *jen.Statement {
	xerror.PanicT(value == "", "params is zero")

	var params []string
	for _, p := range forReg.Split(value, -1) {
		if p = trim(p); p == "" {
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

	return jen.Func().Params().Params(jen.Id("e").Qual(vectyPackage, "List")).BlockFunc(func(g *jen.Group) {
		g.For(s.Op(":=").Range().Add(exp(file, "t."+params[len(params)-1]))).BlockFunc(func(f *jen.Group) {
			f.Id("e=").Append(jen.Id("e"), child)
		})
		g.Return()
	}).Call()
}
var onExp = func(v string) *jen.Statement {
	v = strings.ReplaceAll(v, "$event", "value")
	if strings.Contains(v, "=") {
		return jen.Func().Params(jen.Id("value").String()).Block(jen.Id("t." + v))
	} else {
		return jen.Id("t." + v)
	}
}

var IfExp string
var OptExp string

func styleExp(file *jen.File, g *jen.Group, e string) {
	e = trim(e)

	switch {
	case e == "":
		return
	case len(e) > 2 && e[:1] == "{" && e[len(e)-1:] == "}":
		for _, c := range strings.Split(trim(e[1:len(e)-1]), ",") {
			if c = trim(c); c == "" {
				continue
			}

			if _c := strings.Split(c, ":"); len(_c) == 2 && trim(_c[0]) != "" && trim(_c[1]) != "" {
				g.Qual(vectyPackage, "Style").Call(jen.Lit(trim(_c[0])), exp(file, _c[1]))
			}
		}
		return
	default:
		return
	}
}

func classExp(file *jen.File, g *jen.Group, e string) {
	e = trim(e)

	switch {
	case e == "":
		return
	case len(e) > 2 && e[:1] == "[" && e[len(e)-1:] == "]": //[]
		for _, c := range strings.Split(trim(e[1:len(e)-1]), ",") {
			if c = trim(strings.Trim(strings.Trim(trim(c), "{"), "}")); c == "" {
				continue
			}

			if _c := strings.Split(c, ":"); len(_c) == 2 && trim(_c[0]) != "" && trim(_c[1]) != "" {
				g.Qual(vectyPackage, "MarkupIf").Call(exp(file, _c[1]), jen.Qual(vectyPackage, "Class").Call(jen.Lit(trim(_c[0]))))
			} else {
				g.Qual(vectyPackage, "Class").Call(exp(file, c))
			}
		}
		return
	case len(e) > 2 && e[:1] == "{" && e[len(e)-1:] == "}": // {}
		for _, c := range strings.Split(trim(e[1:len(e)-1]), ",") {
			if c = trim(c); c == "" {
				continue
			}

			if _c := strings.Split(c, ":"); len(_c) == 2 && trim(_c[0]) != "" && trim(_c[1]) != "" {
				g.Qual(vectyPackage, "ClassMap").Call(jen.Lit(trim(_c[0])), exp(file, _c[1]))
			}
		}
		return
	default:
		return
	}
}

func exp(file *jen.File, e string) *jen.Statement {
	e = trim(e)

	switch {
	case e == "":
		return nil
	case strings.Contains(e, "$"):
		return exp(file, strings.ReplaceAll(e, "$", "t."))
	case len(e) > 2 && e[0:1] == "'" && e[len(e)-1:] == "'":
		return exp(file, fmt.Sprintf(`"%s"`, e[1:len(e)-1]))
	case len(e) > 2 && e[0:1] == `"` && e[len(e)-1:] == `"`:
		return jen.Id(e)
	case strings.Contains(e, "{{{") && strings.Contains(e, "}}}") && ternaryBrace.MatchString(e): //{{{}}}
		_d := ternaryBrace.FindStringSubmatch(e)
		return jen.Qual(vectyPackage, "Markup").Call(jen.Qual(vectyPackage, "UnsafeHTML").CallFunc(func(g *jen.Group) {
			_exp := exp(file, trim(_d[2]))
			if _d[1] == "" && _d[3] == "" {
				g.Add(_exp)
			}

			if _d[1] == "" && _d[3] != "" {
				g.Add(_exp).Op("+").Lit(_d[3])
			}

			if _d[1] != "" && _d[3] == "" {
				g.Lit(_d[1]).Op("+").Add(_exp)
			}
		}))
	case twoBrace.MatchString(e): //{{}}
		_d := twoBrace.FindStringSubmatch(e)
		return jen.Qual(vectyPackage, "Text").CallFunc(func(g *jen.Group) {
			_exp := exp(file, trim(_d[2]))
			if _d[1] == "" && _d[3] == "" {
				g.Add(_exp)
			}

			if _d[1] == "" && _d[3] != "" {
				g.Add(_exp).Op("+").Lit(_d[3])
			}

			if _d[1] != "" && _d[3] == "" {
				g.Lit(_d[1]).Op("+").Add(_exp)
			}
		})
	case len(e) > 2 && e[:1] == "[" && e[len(e)-1:] == "]": //[]
		return nil
	case len(e) > 2 && e[:1] == "{" && e[len(e)-1:] == "}": // {}
		return nil
	case len(strings.Split(e, "?:")) == 2:
		_s := strings.Split(e, "?:")
		s0 := trim(_s[0])
		s1 := trim(_s[1])
		return jen.Func().Params().String().BlockFunc(func(g *jen.Group) {
			g.If(jen.Qual("github.com/pubgo/g/pkg", "IsNone").Call(exp(file, s0))).BlockFunc(func(g *jen.Group) {
				g.Return(exp(file, s1))
			}).Else().BlockFunc(func(g *jen.Group) {
				g.Return(exp(file, s0))
			})
		}).Call()
	case ternaryReg.MatchString(e):
		_s := ternaryReg.FindStringSubmatch(e)
		s1 := trim(_s[1])
		s2 := trim(_s[2])
		s3 := trim(_s[3])
		return jen.Func().Params().String().BlockFunc(func(g *jen.Group) {
			g.If(exp(file, s1)).BlockFunc(func(g *jen.Group) {
				g.Return(exp(file, s2))
			}).Else().BlockFunc(func(g *jen.Group) {
				g.Return(exp(file, s3))
			})
		}).Call()
	default:
		return jen.Id(e)
	}
}

func componentAttr(file *jen.File, d jen.Dict, g *jen.Group, attr xml.Attr) {
	ns := trim(attr.Name.Space)
	key := trim(attr.Name.Local)
	key = strings.TrimLeft(key, ns)
	value := trim(attr.Value)
	//value = strings.ReplaceAll(value, "$", "t.")
	value = strings.ReplaceAll(value, `'`, `"`)

	switch {
	case ns == "" && (key == "v-for" || key == "xmlns" || key == "v-if"):
		return
	case ns == "" && len(key) > 5 && key[:5] == "data-" && value != "":
		if key = trim(key[5:]); key == "" {
			return
		}

		if ternaryBrace.MatchString(key) || twoBrace.MatchString(key) {
			if _exp := exp(file, key); _exp != nil {
				file.ImportName(vectyPackage, "vecty")
				if d != nil {
					d[jen.Lit(strings.Title(key))] = _exp
				}

				if g != nil {
					g.Qual(vectyPackage, "Data").Call(jen.Lit(key), _exp)
				}
			}
			return
		}

		if d != nil {
			d[jen.Lit(strings.Title(key))] = jen.Lit(key)
		}

		if g != nil {
			g.Qual(vectyPackage, "Data").Call(jen.Lit(key), jen.Lit(key))
		}
		return

	case ns == "v-on" || key[0] == '@':

		key = strings.ReplaceAll(strings.Title(key[1:]), "-", "")
		if d != nil {
			d[jen.Lit("On"+key)] = onExp(value)
		}

		if g != nil {
			g.Qual(vectyEventPackage, key).Call(onExp(value))
		}
		return
	case ns == "v-bind" || key[0] == ':':
		key = trim(key[1:])

		if _key, ok := stringProps[key]; ok {
			g.Qual(vectyPropPackage, _key).Call(exp(file, value))
			return
		}

		switch key {
		case "style":
			styleExp(file, g, value)
			return
		case "class":
			classExp(file, g, value)
			return
		default:
			if len(key) > 2 && key[0] == '[' && key[len(key)-1] == ']' {
				key = trim(key[1 : len(key)-1])
				if d != nil {
					d[jen.Id("t."+key)] = exp(file, value)
				}

				if g != nil {
					g.Qual(vectyPackage, "Data").Call(jen.Id("t."+key), exp(file, value))
				}
				return
			}

			if d != nil {
				d[jen.Lit(strings.Title(key))] = exp(file, value)
			}

			if g != nil {
				g.Qual(vectyPackage, "Data").Call(jen.Lit(key), exp(file, value))
			}
		}
		return
	case key == "v-model":
		if value == "" {
			return
		}

		value = "t." + value
		if d != nil {
			d[jen.Lit("Value")] = jen.Id(value)
			d[jen.Lit("OnInput")] = jen.Func().Params(jen.Id("_value").String()).BlockFunc(func(g *jen.Group) {
				g.Id(value).Op("=").Id("_value")
			})
		}

		if g != nil {
			file.ImportName(vectyPropPackage, "prop")
			file.ImportName(vectyEventPackage, "event")
			file.ImportAlias(dom, "dom")

			g.Qual(vectyPropPackage, "Value").Call(jen.Id(value))
			g.Qual(vectyEventPackage, "Input").CallFunc(func(g *jen.Group) {
				g.Func().Params(jen.Id("e").Op("*").Qual(vectyPackage, "Event")).BlockFunc(func(g *jen.Group) {
					g.Id(value).Op("=").Qual(dom, "WrapEvent(e.Target).Target().TextContent()")
					g.Qual(dom, "WrapEvent(e.Target).PreventDefault()")
				})
			})
		}
		return
	case key == "v-html":
		if value != "" {
			g.Add(exp(file, "{{{"+value+"}}}"))
		}
		return
	case key == "v-text":
		if value != "" {
			g.Add(exp(file, "{{"+value+"}}"))
		}
		return
	case key == "v-focus":
		if g == nil {
			return
		}

		g.Qual(vectyPackage, "Autofocus").Call()
		return
	case key == "v-pre":
		return
	case key == "v-once":
		return
	case key == "style":
		if g == nil {
			return
		}

		css, err := parser.ParseDeclarations(value)
		xerror.PanicM(err, "css parsing error")

		for _, dec := range css {
			if dec.Important {
				dec.Value += "!important"
			}
			g.Qual(vectyPackage, "Style").Call(
				jen.Lit(dec.Property),
				exp(file, dec.Value),
			)
		}
		return
	case key == "class":
		if g == nil {
			return
		}

		g.Qual(vectyPackage, "Class").CallFunc(func(g *jen.Group) {
			if twoBrace.MatchString(value) {
				_d := twoBrace.FindStringSubmatch(value)
				g.Add(exp(file, trim(_d[2])))

				for _, c := range strings.Split(_d[1]+" "+_d[3], " ") {
					if trim(c) == "" {
						continue
					}

					g.Lit(trim(c))
				}
			} else {
				for _, c := range strings.Split(value, " ") {
					if trim(c) == "" {
						continue
					}

					g.Lit(trim(c))
				}
			}
		})
	default:
		fmt.Println(ns, key, value, "ok")
	}

}

func componentElement(file *jen.File, appPackage string, token xml.StartElement) *jen.Statement {
	ts := strings.Split(trim(token.Name.Local), ":")[1:]
	name := trim(ts[len(ts)-1])

	appPackage += "/components"
	if len(ts) > 1 {
		appPackage += "/" + strings.Join(ts[:len(ts)-1], "/")
		file.ImportAlias(appPackage, ts[len(ts)-2])
	} else {
		file.ImportAlias(appPackage, "components")
	}

	return jen.Qual(appPackage, strings.ReplaceAll(strings.Title(name), "-", "")).CallFunc(func(g *jen.Group) {
		if len(token.Attr) > 0 {
			g.Map(jen.String()).Interface().Values(jen.DictFunc(func(d jen.Dict) {
				for _, v := range token.Attr {
					componentAttr(file, d, nil, v)
				}
			}))

			for _, v := range token.Attr {
				componentAttr(file, nil, g, v)
			}
		}
	})
}
