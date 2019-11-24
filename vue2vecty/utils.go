package vue2vecty

import (
	"encoding/xml"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/gopherjs/vecty"
	"reflect"
	"strings"
)

func createStruct(appPackage, componentName, packageName string) *jen.File {
	file := jen.NewFile(packageName)
	file.PackageComment("This file was created with https://github.com/pubgo/factor")
	file.PackageComment("using https://jsgo.io/dave/git_repo")
	file.ImportName(vectyPackage, "vecty")
	file.ImportName(vectyElemPackage, "elem")
	file.ImportName(vectyPropPackage, "prop")
	file.ImportName(vectyEventPackage, "event")
	file.ImportName(vectyStylePackage, "style")

	file.Type().Id(componentName).Struct(
		jen.Qual("github.com/gopherjs/vecty", "Core"),
	)
	return file
}

func componentElement(appPackage, componentName string, token *xml.StartElement) *jen.Statement {
	vectyPackage := appPackage + "/components"
	// vectyFunction = component
	// vectyParamater = tag
	var component = strings.TrimLeft(token.Name.Local, "components.")
	// I'm not sure what qual was intended to mean (it's always true now) but it looks like perhaps you're
	// trying to avoid using Qual if the package path == local path? If so, no need! Qual handles this
	// gracefully... See: https://github.com/dave/jennifer#qual
	return jen.Op("&").Qual(vectyPackage, component).Values(jen.DictFunc(func(d jen.Dict) {
		for _, v := range token.Attr {
			d[jen.Id(v.Name.Local)] = jen.Lit(v.Value)
		}
	}))

	//baseDecl := jen.Op("&").Id(component).Values(jen.DictFunc(func(d jen.Dict) {
	//	for _, v := range token.Attr {
	//		d[jen.Id(v.Name.Local)] = jen.Lit(v.Value)
	//	}
	//}))
	//return baseDecl
}

func Components(c ...vecty.ComponentOrHTML) (l vecty.List) {
	return c
}

func Str(c ...string) []string {
	return c
}

func Css(m ...string) vecty.ClassMap {
	return ClassMap(m...)
}

func ClassMap(m ...string) vecty.ClassMap {
	_css := vecty.ClassMap{}
	for _, i := range m {
		if i == "" {
			continue
		}

		for _, _s := range strings.Split(i, " ") {
			_css[_s] = true
		}

	}
	return _css
}

type EleHandle func(markup ...vecty.MarkupOrChild) *vecty.HTML
type EventHandle func(*vecty.Event)

func EleIf(b bool, a, a1 EleHandle) EleHandle {
	if b {
		return a
	}

	return a1
}

type StrM map[string]string

type M map[string]interface{}

func (t M) Map(fn func(k string, v interface{}) vecty.ComponentOrHTML) vecty.List {
	var _dm vecty.List
	for k, v := range t {
		_dm = append(_dm, fn(k, v))
	}
	return _dm
}

func init() {
	fmt.Println([]StrM{{"s": "22"}})
}

type MapComponent map[string]vecty.ComponentOrHTML

type ComponentFn func(style ...string) func(c ...vecty.ComponentOrHTML) vecty.ComponentOrHTML

func MapElem(cpt []vecty.ComponentOrHTML, fn func(c vecty.ComponentOrHTML) vecty.ComponentOrHTML) vecty.List {
	var _cs vecty.List
	for _, _i := range cpt {
		_cs = append(_cs, fn(_i))
	}
	return _cs
}

func MapE(e vecty.ComponentOrHTML, data interface{}, fn interface{}) vecty.List {
	var _cs vecty.List
	for _, _i := range cpt {
		_cs = append(_cs, fn(_i))
	}
	return _cs
}

func Map(data interface{}, fn interface{}) (_l vecty.List) {
	_d := reflect.ValueOf(data)
	for i := 0; i < _d.Len(); i++ {
		_dt := _d.Index(i)
		if !_dt.IsValid() {
			continue
		}
		_v := reflect.ValueOf(fn).Call([]reflect.Value{_d.Index(i)})
		_l = append(_l, _v[0].Interface().(vecty.ComponentOrHTML))
	}
	return
}

/*
data=[1,2,3,4,5]
Map(data,func(i int)vecty.ComponentOrHTML{
	ddd:=data[i]
	return elem.Button()
})
*/
