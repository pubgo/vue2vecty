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

var ternaryReg = xerror.PanicErr(regexp.Compile(`(.+)\?(.+):(.+)`)).(*regexp.Regexp)
var ternaryBraceReg = xerror.PanicErr(regexp.Compile(`(.*){{{(.+)}}}(.*)`)).(*regexp.Regexp)
var twoBraceReg = xerror.PanicErr(regexp.Compile(`(.*){{(.+)}}(.*)`)).(*regexp.Regexp)
var forReg = xerror.PanicErr(regexp.Compile(`,|\s+in\s+`)).(*regexp.Regexp)

func braceExp(file *jen.File, e string) *jen.Statement {
	e = trim(e)

	var _braceExp = func(_d []string) *jen.Statement {
		x1, x2, x3 := _d[1], _d[2], _d[3]
		_brace := jen.Empty()
		if _x := braceExp(file, x1); _x != nil {
			_brace = _brace.Add(_x).Op("+")
		}

		if _x := exp(file, x2); _x != nil {
			_brace = _brace.Add(_x)
		}

		if _x := braceExp(file, x3); _x != nil {
			_brace = _brace.Op("+").Add(_x)
		}

		return _brace
	}

	switch {
	case e == "":
		return nil
	case strings.Contains(e, "{{{") && strings.Contains(e, "}}}") && ternaryBraceReg.MatchString(e):
		return _braceExp(ternaryBraceReg.FindStringSubmatch(e))
	case strings.Contains(e, "{{") && strings.Contains(e, "}}") && twoBraceReg.MatchString(e):
		return _braceExp(twoBraceReg.FindStringSubmatch(e))
	default:
		return jen.Lit(e)
	}
}

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
	var __s *jen.Statement
	if len(params) == 1 {
		__s = jen.Id("__key,__value")
		s = jen.Id("key,value")
	}

	if len(params) == 2 {
		__s = jen.Id("__" + params[0])
		s = jen.Id(params[0])
	}

	if len(params) == 3 {
		__s = jen.Id("__" + params[0] + "," + "__" + params[1])
		s = jen.Id(params[0] + "," + params[1])
	}

	xerror.PanicT(s == nil, "statements error")

	return jen.Func().Params().Params(jen.Id("e").Qual(vectyPackage, "List")).BlockFunc(func(g *jen.Group) {
		g.For(jen.Add(__s).Op(":=").Range().Add(exp(file, "t."+params[len(params)-1]))).BlockFunc(func(f *jen.Group) {
			f.Add(s).Op(":=").Add(__s)
			f.Id("e=").Append(jen.Id("e"), child)
		})
		g.Return()
	}).Call()
}
var onExp = func(v string) *jen.Statement {
	v = strings.ReplaceAll(v, "$event", "value")
	if strings.Contains(v, "=") {
		return jen.Func().Params(jen.Id("value").String()).BlockFunc(func(g *jen.Group) {
			g.Id("t." + v)
			g.Qual(vectyPackage, "Rerender").Call(jen.Id("t"))
		})
	} else {
		return jen.Id("t." + v)
	}
}
var styleExp = func(file *jen.File, g *jen.Group, e string) {
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
var classExp = func(file *jen.File, g *jen.Group, e string) {
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

//var IfExp string
//var OptExp string

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
	case strings.Contains(e, "{{{") && strings.Contains(e, "}}}") && ternaryBraceReg.MatchString(e):
		return jen.Qual(vectyPackage, "Markup").Call(jen.Qual(vectyPackage, "UnsafeHTML").Call(braceExp(file, e)))
	case strings.Contains(e, "{{") && strings.Contains(e, "}}") && twoBraceReg.MatchString(e):
		return jen.Qual(vectyPackage, "Text").Call(braceExp(file, e))
	case len(strings.SplitN(e, "||", 2)) == 2:
		_s := strings.SplitN(e, "||", 2)
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

		if ternaryBraceReg.MatchString(key) || twoBraceReg.MatchString(key) {
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
					g.Qual(vectyPackage, "Property").Call(jen.Id("t."+key), exp(file, value))
				}
				return
			}

			if d != nil {
				d[jen.Lit(strings.Title(key))] = exp(file, value)
			}

			if g != nil {
				g.Qual(vectyPackage, "Property").Call(jen.Lit(strings.ToLower(key)), exp(file, value))
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
				g.Qual(vectyPackage, "Rerender").Call(jen.Id("t"))
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
			if twoBraceReg.MatchString(value) {
				_d := twoBraceReg.FindStringSubmatch(value)
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
		if value == "" {
			return
		}

		key = strings.ToLower(key)
		//switch strings.ToLower(key) {
		//case "id":
		//	key = "ID"
		//case "for", "htmlfor", "html-for":
		//	key = "For"
		//case "focus", "autofocus":
		//	key = "Autofocus"
		//case "disabled":
		//	key = "Alt"
		//}

		if d != nil {
			d[jen.Lit(strings.Title(key))] = jen.Lit(value)
		}

		if g != nil {
			g.Qual(vectyPackage, "Property").Call(jen.Lit(key), jen.Lit(value))
		}
	}

}

func CreateComponent(packageName, componentName string) *jen.File {
	file := jen.NewFile(packageName)
	file.PackageComment("This file was created with https://github.com/pubgo/vue2vecty")
	file.ImportName(vectyPackage, "vecty")

	componentName = strings.ReplaceAll(strings.Title(componentName), "-", "")
	_componentName := "_" + componentName

	file.Func().Id(componentName).Params(jen.Id("data").Qual(js, "M"), jen.Id("slots ...vecty.ComponentOrHTML")).Id("vecty.ComponentOrHTML").BlockFunc(func(g *jen.Group) {
		g.Id("t").Op(":=").Op("&").Id(_componentName).Values(jen.Dict{jen.Id("Slot"): jen.Id("slots")})
		g.IfFunc(func(g *jen.Group) {
			g.Id("data").Op("!=").Nil().BlockFunc(func(g *jen.Group) {
				g.IfFunc(func(g *jen.Group) {
					file.ImportName(mapstructure, "mapstructure")
					g.Id("err:=").Qual(mapstructure, "Decode").Call(jen.Id("data"), jen.Id("t")).Id("; err != nil").BlockFunc(func(g *jen.Group) {
						file.ImportName("log", "log")
						g.Qual("log", "Fatalf").Call(jen.Lit("%#v"), jen.Id("err"))
					})
				})
			})
		})
		g.Return(jen.Id("t"))
	})

	file.Type().Id(_componentName).Struct(
		jen.Qual(vectyPackage, "Core"),
		jen.Id("Slot").Qual(vectyPackage, "List"),
	)

	file.Func().Params(jen.Id("t").Op("*").Id(_componentName)).Id("Render").Params().Qual(vectyPackage, "ComponentOrHTML").Block(
		jen.Qual(vectyPackage, "SetTitle").Call(
			jen.Id("t").Dot("GetTitle").Call(),
		),
		jen.Return(jen.Id("t._Render()")),
	)

	return file
}
