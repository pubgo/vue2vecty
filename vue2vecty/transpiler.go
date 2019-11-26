package vue2vecty

import (
	"errors"
	"fmt"
	"github.com/pubgo/g/xerror"
	"io/ioutil"
	"regexp"

	"bytes"
	"io"

	"strings"

	"encoding/xml"

	"github.com/aymerick/douceur/parser"
	"github.com/dave/jennifer/jen"
)

var callRegexp = regexp.MustCompile(`{vecty-call:([a-zA-Z0-9_\-]+)}`)
var fieldRegexp = regexp.MustCompile(`{vecty-field:([a-zA-Z0-9_\-]+})`)
var EOT = errors.New("end of tag")

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

	call := jen.Options{
		Close:     ")",
		Multi:     true,
		Open:      "(",
		Separator: ",",
	}

	decoder := xml.NewDecoder(bytes.NewBufferString(s.html))

	var _transcode func(*xml.Decoder) ([]jen.Code, error)
	_transcode = func(decoder *xml.Decoder) (code []jen.Code, err error) {
		defer xerror.RespErr(&err)

		token := xerror.PanicErr(decoder.Token()).(xml.Token)
		switch token := token.(type) {
		case xml.StartElement:
			tag := token.Name.Local
			ns := token.Name.Space

			vectyFunction, ok := elemNames[tag]
			_vectyPackage := vectyPackage
			vectyParamater := ""
			var ce *jen.Statement
			if !ok {
				if strings.HasPrefix(ns, "c") {
					// not sure if we need this?
					if strings.HasPrefix(ns, "c:") {
						ce = componentElement(file, s.appPackage, &token)
						vectyFunction = ""
					} else {
						vectyFunction = "Tag"
						_vectyPackage = "github.com/gopherjs/vecty"
						vectyParamater = tag
					}
				}

				// 处理for
				_for := func(child *jen.Statement, params ...string) *jen.Statement {
					xerror.PanicT(len(params) == 0, "params is zero")

					var s *jen.Statement

					if len(params) == 1 {
						s = jen.Id("key").Id(",").Id("value")
					}

					if len(params) == 2 {
						s = jen.Id(params[0])
					}

					if len(params) == 3 {
						s = jen.Id(params[0]).Id(",").Id(params[1])
					}

					xerror.PanicT(s == nil, "statements error")

					return jen.Func().Params().Params(jen.Id("e").Qual("github.com/gopherjs/vecty", "List")).BlockFunc(func(g *jen.Group) {
						g.For(s.Op(":=").Id("range").Qual("t", params[len(params)-1])).BlockFunc(func(f *jen.Group) {
							f.Id("e").Op("=").Id("append").Call(jen.Id("e"), child)
						})
						g.Return()
					}).Call()
				}

				var outer error
				q := jen.Qual(_vectyPackage, vectyFunction).CustomFunc(call, func(g *jen.Group) {
					if vectyParamater != "" {
						g.Lit(vectyParamater)
					}

					if ce == nil && len(token.Attr) > 0 {
						g.Qual("github.com/gopherjs/vecty", "Markup").CustomFunc(call, func(g *jen.Group) {
							for _, v := range token.Attr {
								space := v.Name.Space
								local := v.Name.Local
								switch {
								case local == "v-if":
									file.ImportName("github.com/gopherjs/vecty", "If")
									continue
								case local == "v-for":
									continue
								case local == "v-model":
								case space == "v-bind" || space == "v":
								case local == "v-on":
								case v.Name.Local == "style":
									css, err := parser.ParseDeclarations(v.Value)
									if err != nil {
										outer = err
										return
									}
									for _, dec := range css {
										if dec.Important {
											dec.Value += "!important"
										}
										g.Qual("github.com/gopherjs/vecty", "Style").Call(
											jen.Lit(dec.Property),
											jen.Lit(dec.Value),
										)
									}
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
						})
					}

					for {

						// 这是 子元素
						c, err := _transcode(decoder)
						if err != nil {
							if err == EOT {
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

				if ce != nil {
					return []jen.Code{ce}, nil
				}

				for _, v := range token.Attr {
					// 处理 v-for
					if v.Name.Space == "" && v.Name.Local == "v-for" {
						vs := strings.Split(v.Value, "in")
						return []jen.Code{_for(q, strings.TrimSpace(vs[0]), strings.TrimSpace(vs[1]))}, nil
					}

					// 处理 v-if
					//if v.Name.Space == "" && v.Name.Local == "v-if" {
					//	return []jen.Code{_g(q, v.Value)}, nil
					//}

					// 处理 {}
					// 处理 {{  }}
				}

				return []jen.Code{q}, nil
			}
		case xml.CharData:
			str := string(token)
			hasCall := callRegexp.MatchString(str)
			hasField := fieldRegexp.MatchString(str)
			hasSpecial := hasCall || hasField

			if hasSpecial {
				if hasCall {
					var statements []jen.Code
					crResult := callRegexp.FindAllStringIndex(str, -1)
					index := 0
					for matchNumber, match := range crResult {
						var before, between, after string
						before = str[index:match[0]]
						fnCall := str[match[0]:match[1]]
						fnCall = strings.TrimLeft(fnCall, "{vecty-call:")
						fnCall = strings.Replace(fnCall, "}", "", -1)
						if matchNumber < len(crResult)-1 {
							// there's another match
							between = str[match[1]:crResult[matchNumber+1][0]]
						}
						after = str[match[1]:]

						if before != "" && !strings.Contains(before, "vecty-call") {
							statements = append(statements, jen.Qual("github.com/gopherjs/vecty", "Text").Call(
								jen.Lit(before),
							))
						}
						statements = append(statements, jen.Qual("github.com/gopherjs/vecty", "Text").Call(
							jen.Id("p").Dot(fnCall).Call(),
						))
						if between != "" && !strings.Contains(between, "vecty-call") {
							statements = append(statements, jen.Qual("github.com/gopherjs/vecty", "Text").Call(
								jen.Lit(between),
							))
						}
						if after != "" && !strings.Contains(after, "vecty-call") {
							statements = append(statements, jen.Qual("github.com/gopherjs/vecty", "Text").Call(
								jen.Lit(after),
							))
						}
					}
					return statements, nil

				}
				if hasField {
					var statements []jen.Code
					crResult := fieldRegexp.FindAllStringIndex(str, -1)
					index := 0
					for matchNumber, match := range crResult {
						var before, between, after string
						before = str[index:match[0]]
						field := str[match[0]:match[1]]
						field = strings.TrimLeft(field, "{vecty-field:")
						field = strings.Replace(field, "}", "", -1)
						if matchNumber < len(crResult)-1 {
							// there's another match
							between = str[match[1]:crResult[matchNumber+1][0]]
						}
						after = str[match[1]:]
						if before != "" && !strings.Contains(before, "vecty-field") {
							statements = append(statements, jen.Qual("github.com/gopherjs/vecty", "Text").Call(
								jen.Lit(before),
							))
						}
						statements = append(statements, jen.Qual("github.com/gopherjs/vecty", "Text").Call(
							jen.Id("p").Dot(field),
						))
						if between != "" && !strings.Contains(between, "vecty-field") {
							statements = append(statements, jen.Qual("github.com/gopherjs/vecty", "Text").Call(
								jen.Lit(between),
							))
						}
						if after != "" && !strings.Contains(after, "vecty-field") {
							statements = append(statements, jen.Qual("github.com/gopherjs/vecty", "Text").Call(
								jen.Lit(after),
							))
						}
					}
					return statements, nil

				}
			}
			s := strings.TrimSpace(string(token))
			if s == "" {
				return nil, nil
			}
			return []jen.Code{jen.Qual("github.com/gopherjs/vecty", "Text").Call(jen.Lit(s))}, nil
		case xml.EndElement:
			return nil, EOT
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
			if err == io.EOF || err == EOT {
				break
			}
			s.code = fmt.Sprintf("%s", err)
			//return
		}
		if c != nil {
			elements = append(elements, c...)
		}
	}

	if s.packageName == "routes" || s.packageName == "pages" {
		file.Func().Params(jen.Id("p").Op("*").Id(s.componentName)).Id("Render").Params().Qual("github.com/gopherjs/vecty", "ComponentOrHTML").Block(
			jen.Qual("github.com/gopherjs/vecty", "SetTitle").Call(
				jen.Id("p").Dot("GetTitle").Call(),
			),
			jen.Return(
				// TODO: wrap in if - only body for a "route"
				jen.Qual("github.com/gopherjs/vecty/elem", "Body").Custom(call, elements...),
			),
		)
	} else {
		file.Func().Params(jen.Id("p").Op("*").Id(s.componentName)).Id("Render").Params().Qual("github.com/gopherjs/vecty", "ComponentOrHTML").Block(
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
