package vue2vecty

import (
	"bytes"
	"fmt"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/vue2vecty/xml"
	"io"
	"io/ioutil"
	"regexp"
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
	xerror.PanicM(s.transform(), "transform error")
	return s
}

type Transpiler struct {
	reader        io.Reader
	appPackage    string
	componentName string
	packageName   string
	html, code    string
}

func (t *Transpiler) read() {
	bb, err := ioutil.ReadAll(t.reader)
	xerror.PanicM(err, "reading component template error")
	t.html = string(bb)
}

func (t *Transpiler) Code() string {
	return t.code
}

func (t *Transpiler) transform() (err error) {
	defer xerror.RespErr(&err)

	file := jen.NewFile(t.packageName)
	file.PackageComment("This file was created with https://github.com/pubgo/vue2vecty")
	file.PackageComment("using https://jsgo.io/pubgo/vue2vecty")
	file.ImportNames(map[string]string{
		vectyPackage:      "vecty",
		vectyElemPackage:  "elem",
		vectyPropPackage:  "prop",
		vectyEventPackage: "event",
		vectyStylePackage: "style",
	})

	decoder := xml.NewDecoder(bytes.NewBufferString(t.html))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	decoder.Entity = xml.HTMLEntity

	var _transform func(*xml.Decoder) ([]jen.Code, error)
	_transform = func(decoder *xml.Decoder) (code []jen.Code, err error) {
		defer xerror.RespErr(&err)

		token, err := decoder.Token()
		if err == io.EOF || token == nil {
			return nil, xerror.Errs.Done
		}
		xerror.Panic(err)

		var ce *jen.Statement
		var _appPackage = vectyElemPackage

		switch token := token.(type) {
		case xml.StartElement:
			tag := trim(token.Name.Local)
			tagName, ok := elemNames[tag]
			if !ok {
				ts := strings.Split(tag, ":")
				name := trim(ts[len(ts)-1])

				_appPackage = t.appPackage + "/components"
				if len(ts) > 1 {
					_appPackage += "/" + strings.Join(ts[:len(ts)-1], "/")
					file.ImportName(_appPackage, ts[len(ts)-2])
				} else {
					file.ImportName(_appPackage, "components")
					if t.packageName == "components" {
						_appPackage = ""
					}
				}

				name = strings.ReplaceAll(strings.Title(name), "-", "")
				ce = jen.Qual(_appPackage, name).CallFunc(func(g *jen.Group) {
					if len(token.Attr) > 0 {
						g.Map(jen.String()).Interface().Values(jen.DictFunc(func(d jen.Dict) {
							for _, v := range token.Attr {
								componentAttr(file, d, nil, v)
							}
						}))
					}

					for {
						c, err := _transform(decoder)
						if err != nil {
							if err == xerror.Errs.Done {
								break
							}
							xerror.Panic(err)
						}
						if c != nil {
							g.Add(c...)
						}
					}
				})
			} else {
				ce = jen.Qual(_appPackage, tagName).CallFunc(func(g *jen.Group) {
					if len(token.Attr) > 0 {
						g.Qual(vectyPackage, "Markup").CallFunc(func(g *jen.Group) {
							for _, v := range token.Attr {
								componentAttr(file, nil, g, v)
							}
						})
					}

					for {
						c, err := _transform(decoder)
						if err != nil {
							if err == xerror.Errs.Done {
								break
							}
							xerror.Panic(err)
						}
						if c != nil {
							g.Add(c...)
						}
					}
				})
			}

			for _, attr := range token.Attr {
				ns := trim(attr.Name.Space)
				key := trim(attr.Name.Local)
				key = strings.TrimLeft(key, ns)
				value := trim(attr.Value)

				if ns == "" && (key == "v-for" || key == "v-range") && value != "" {
					ce = forExp(file, ce, value)
				}

				if ns == "" && (key == "v-if" || key == "v-show") && value != "" {
					ce = jen.Qual(vectyPackage, "If").Call(exp(file, value), ce)
				}
			}

			return []jen.Code{ce}, nil
		case xml.CharData:
			e := trim(string(token))
			if e == "" {
				return nil, nil
			}

			var _code []jen.Code
			for _, _e := range strings.Split(e, "\n") {
				if _e = trim(_e); _e == "" {
					continue
				}

				// strings.Contains(_e, "{{") && strings.Contains(_e, "}}")
				_check := xerror.PanicErr(regexp.Compile(`(.*){{(.*)}}(.*)`)).(*regexp.Regexp)
				if _check.MatchString(_e) {
					if ternaryBrace.MatchString(_e) || twoBrace.MatchString(_e) {
						if _exp := exp(file, _e); _exp != nil {
							_code = append(_code, _exp)
						}
					}
					continue
				}
				_code = append(_code, jen.Qual(vectyPackage, "Text").Call(jen.Lit(_e)))
			}

			if len(_code) == 0 {
				return nil, nil
			}
			return _code, nil
		case xml.EndElement:
			return nil, xerror.Errs.Done
		case xml.Comment:
			return nil, nil
		default:
			fmt.Printf("%T %#v \n", token, token)
		}
		return nil, nil
	}
	var elements []jen.Code
	for {
		c, err := _transform(decoder)
		if err != nil {
			if err == io.EOF || err == xerror.Errs.Done {
				break
			}
			xerror.Panic(err)
		}

		if c != nil {
			elements = append(elements, c...)
		}
	}

	if t.packageName == "routes" || t.packageName == "pages" || t.packageName == "views" {
		file.Func().Params(jen.Id("t").Id("*"+t.componentName)).Id("Render").Params().Qual(vectyPackage, "ComponentOrHTML").Block(
			jen.Qual(vectyPackage, "SetTitle").Call(
				jen.Id("p").Dot("GetTitle").Call(),
			),
			jen.Return(
				jen.Qual(vectyElemPackage, "Body").Call(elements...),
			),
		)
	} else {
		file.Func().Params(jen.Id("t").Op("*").Id(t.componentName)).Id("Render").Params().Qual(vectyPackage, "ComponentOrHTML").Block(
			jen.Return(elements...),
		)
	}

	buf := &bytes.Buffer{}
	defer buf.Reset()
	xerror.Panic(file.Render(buf))
	t.code = buf.String()
	return
}
