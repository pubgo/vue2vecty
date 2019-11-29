package vue2vecty

import (
	"bytes"
	"fmt"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/vue2vecty/xml"
	"io"
	"io/ioutil"
	"strings"

	"github.com/dave/jennifer/jen"
)

func NewTranspiler(r io.Reader, appPackage, componentName, packageName string) *Transpiler {
	defer xerror.Assert()

	s := &Transpiler{
		reader:        r,
		appPackage:    appPackage,
		packageName:   packageName,
		componentName: componentName,
	}
	s.read()
	xerror.PanicM(s.transcode(), "transcode error")
	return s
}

type Transpiler struct {
	reader        io.Reader
	appPackage    string
	componentName string
	packageName   string
	html, code    string
}

func (s *Transpiler) read() {
	bb, err := ioutil.ReadAll(s.reader)
	xerror.PanicM(err, "reading component template error")
	s.html = string(bb)
}

func (s *Transpiler) Code() string {
	return s.code
}

func (s *Transpiler) transcode() (err error) {
	defer xerror.RespErr(&err)

	file := jen.NewFile(s.packageName)
	file.PackageComment("This file was created with https://github.com/pubgo/factor")
	file.PackageComment("using https://jsgo.io/dave/html2vecty")
	file.ImportNames(map[string]string{
		"github.com/gopherjs/vecty":       "vecty",
		"github.com/gopherjs/vecty/elem":  "elem",
		"github.com/gopherjs/vecty/prop":  "prop",
		"github.com/gopherjs/vecty/event": "event",
		"github.com/gopherjs/vecty/style": "style",
	})

	decoder := xml.NewDecoder(bytes.NewBufferString(s.html))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	decoder.Entity = xml.HTMLEntity

	var _transcode func(*xml.Decoder) ([]jen.Code, error)
	_transcode = func(decoder *xml.Decoder) (code []jen.Code, err error) {
		defer xerror.Assert()

		token, err := decoder.Token()
		if err == io.EOF || token == nil {
			return nil, xerror.ErrDone
		}

		switch token := token.(type) {
		case xml.StartElement:
			ns := strings.TrimSpace(token.Name.Space)
			tag := strings.TrimSpace(token.Name.Local)
			tag = strings.TrimLeft(tag, ns)

			var ce *jen.Statement
			var outer error
			if strings.HasPrefix(ns, "c") {
				ts := strings.Split(trim(token.Name.Local), ":")[1:]
				name := trim(ts[len(ts)-1])

				_appPackage := s.appPackage + "/components"
				if len(ts) > 1 {
					_appPackage += "/" + strings.Join(ts[:len(ts)-1], "/")
					file.ImportAlias(_appPackage, ts[len(ts)-2])
				} else {
					file.ImportAlias(_appPackage, "components")
				}

				ce = jen.Qual(_appPackage, strings.ReplaceAll(strings.Title(name), "-", "")).CallFunc(func(g *jen.Group) {
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

					for {
						// 这是 子元素
						c, err := _transcode(decoder)
						if err != nil {
							if err == xerror.ErrDone {
								break
							}
							outer = err
							return
						}
						if c != nil {
							g.Add(c...)
						}
					}
				})
			} else {
				vectyFunction, ok := elemNames[tag]
				xerror.PanicT(!ok, "element(%s) not found", tag)
				ce = jen.Qual(vectyElemPackage, vectyFunction).CallFunc(func(g *jen.Group) {
					if len(token.Attr) > 0 {
						g.Qual(vectyPackage, "Markup").CallFunc(func(g *jen.Group) {
							for _, v := range token.Attr {
								componentAttr(file, nil, g, v)
							}
						})
					}

					for {
						// 这是 子元素
						c, err := _transcode(decoder)
						if err != nil {
							if err == xerror.ErrDone {
								break
							}
							outer = err
							return
						}
						if c != nil {
							g.Add(c...)
						}
					}
				})
				if outer != nil {
					return nil, outer
				}
			}

			for _, attr := range token.Attr {
				ns := strings.TrimSpace(attr.Name.Space)
				key := strings.TrimSpace(attr.Name.Local)
				key = strings.TrimLeft(key, ns)
				value := strings.TrimSpace(attr.Value)

				if ns == "" && (key == "v-if" || key == "v-show") && value != "" {
					ce = jen.Qual(vectyPackage, "If").Call(exp(file, value), ce)
					break
				}

				if ns == "" && key == "v-for" {
					ce = _for(file, ce, value)
					break
				}
			}

			return []jen.Code{ce}, nil
		case xml.CharData:
			e := trim(string(token))
			if e == "" {
				return nil, nil
			}

			if ternaryBrace.MatchString(e) || twoBrace.MatchString(e) {
				if _exp := exp(file, string(token)); _exp != nil {
					file.ImportName(vectyPackage, "vecty")
					return []jen.Code{_exp}, nil
				}
			} else {
				return []jen.Code{jen.Qual(vectyPackage, "Text").Call(jen.Lit(e))}, nil
			}
		case xml.EndElement:
			return nil, xerror.ErrDone
		case xml.Comment:
			return nil, nil
		default:
			fmt.Printf("%T %#v \n", token, token)
		}
		return nil, nil
	}
	var elements []jen.Code
	for {
		c, err := _transcode(decoder)
		if err != nil {
			if err == io.EOF || err == xerror.ErrDone {
				break
			}
			s.code = fmt.Sprintf("%s", err)
		}

		if c != nil {
			elements = append(elements, c...)
		}
	}

	if s.packageName == "routes" || s.packageName == "pages" {
		file.Func().Params(jen.Id("t").Id("*"+s.componentName)).Id("Render").Params().Qual("github.com/gopherjs/vecty", "ComponentOrHTML").Block(
			jen.Qual("github.com/gopherjs/vecty", "SetTitle").Call(
				jen.Id("p").Dot("GetTitle").Call(),
			),
			jen.Return(
				// TODO: wrap in if - only body for a "route"
				jen.Qual("github.com/gopherjs/vecty/elem", "Body").Call(elements...),
			),
		)
	} else {
		file.Func().Params(jen.Id("t").Op("*").Id(s.componentName)).Id("Render").Params().Qual("github.com/gopherjs/vecty", "ComponentOrHTML").Block(
			// TODO: wrap in if - only body for a "route"
			// TODO: Are you sure this is right? It looks like if len(elements) > 1 this will break?
			jen.Return(elements...),
		)
	}

	buf := &bytes.Buffer{}
	xerror.Panic(file.Render(buf))
	s.code = buf.String()
	return
}
