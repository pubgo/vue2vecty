package tests_test

import (
	"bytes"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/vue2vecty/vue2vecty"
	"github.com/pubgo/vue2vecty/xml"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	defer xerror.Assert()

	a := `<b:a:todo-item
                    v-for="groceryList"
                    v-bind:todo="item"
                    v-bind:key="item.id"
                    v-on:key="item.id"
                    v-key="item.id" 
					@key="item.id"
					.key="item.id"
					:key="item.id"
					v-bind:[key]="item.id"
					v-focus=true
            >
sss
{{}}

<li></li>
<input/>

</b:a:todo-item>`

	decoder := xml.NewDecoder(bytes.NewBufferString(a))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	decoder.Entity = xml.HTMLEntity

	for {
		token := xerror.PanicErr(decoder.Token()).(xml.Token)

		switch token := token.(type) {
		case xml.StartElement:
			fmt.Println(token.Name.Local)
			fmt.Println(token.Name.Space)
			for _, v := range token.Attr {
				fmt.Println(v.Name.Space, v.Name.Local, "data", strings.TrimLeft(v.Name.Local, v.Name.Space), v.Value)
			}
		case xml.EndElement:
			continue
		case xml.CharData:
			fmt.Println(string(token.Copy()))
		}
	}
}

func TestName1(t *testing.T) {
	_g := func(s *jen.Statement, value string) *jen.Statement {
		vs := strings.Split(value, "in")
		_var := vs[0]
		_data := vs[1]
		return jen.Qual("github.com/gopherjs/vecty", "Map").CallFunc(func(g *jen.Group) {
			g.Id(_data)
			g.Func().Params(jen.Id("i").Int()).Qual("github.com/gopherjs/vecty", "Tag").Block(
				jen.Id(_var).Op(":=").Id(_data).Index(jen.Id("i")),
				jen.Return(s),
			)
		})
	}

	//fmt.Println(_g(jen.Id("p"),"i in data").Render(os.Stdout))
	//fmt.Println(_g(jen.Empty(),"i in data").Render(os.Stdout))
	//fmt.Println(_g(jen.Func().Params().Qual("github.com/gopherjs/vecty", "tag").Call(), "i in data").Render(os.Stdout))
	fmt.Println(_g(jen.Func().Params().Qual("github.com/gopherjs/vecty", "tag").Call(), "d in data").Render(os.Stdout))
}

func TestName2(t *testing.T) {
	_g := func(child *jen.Statement, params ...string) *jen.Statement {
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

		return jen.Return().Func().Params().Params(jen.Id("e").Qual("github.com/gopherjs/vecty", "List")).BlockFunc(func(g *jen.Group) {
			g.For(s.Op(":=").Id("range").Qual("t", params[len(params)-1])).BlockFunc(func(f *jen.Group) {
				f.Id("e").Op("=").Id("append").Call(jen.Id("e"), child)
				f.If(jen.Id(`len("")>0`)).Block()
			})
			g.Return()
		}).Call()
	}

	fmt.Println(_g(jen.Qual("elem", "Heading2").Call(), "data").Render(os.Stdout))
	fmt.Println(_g(jen.Qual("elem", "Heading2").Call(), "k", "data").Render(os.Stdout))
	fmt.Println(_g(jen.Qual("elem", "Heading2").Call(), "k", "v", "data").Render(os.Stdout))

	//func() (e vecty.List) {
	//	for k := range t.Data {
	//		e = append(e, elem.Heading2(vecty.Text("Class attributes"), vecty.Markup(
	//			vecty.Style("border", k),
	//			vecty.Style("color", "red!important"),
	//		), ))
	//	}
	//	return
	//}()
}

func TestName3(t *testing.T) {
	a := `
<ul>

<li
                    v-for="$groceryList"
                    v-model="$todo"
					@key="$click"
					@click="$onClick"
					:key="item.id"
					class="a b c"
					data-hello="hello"
					v-focus=true
            >
hello
{{$moke}}

<li> {{{$world}}} </li>
<input/>

</li>


<p v-if="$hello>0"></p>

</ul>
`

	v := vue2vecty.NewTranspiler(bytes.NewBufferString(a), "github.com/pubgo/vue2vecty", "Test", "views")
	fmt.Println(v.Code())
}

func TestName4(t *testing.T) {
	fmt.Println(strings.ReplaceAll(strings.Title("hello-hello"), "-", ""))
}

func TestA5(t *testing.T) {
	ternary := regexp.MustCompile(`(.+)\?(.+):(.+)`)
	//ternary := regexp.MustCompile(`(.+)\?:(.+)`)
	fmt.Println(ternary.Split("a>b ? dd : 22", -1))
	fmt.Println(ternary.MatchString("a>b?dd:22"))
	fmt.Println(ternary.FindStringSubmatch("a>b?dd:22"))
	fmt.Println(ternary.FindStringSubmatch("a>b?:22"))

	//var ternaryBrace = xerror.PanicErr(regexp.Compile(`.*{{{(.+)}}}.*`)).(*regexp.Regexp)
	var twoBrace = xerror.PanicErr(regexp.Compile(`.*{{(.+)}}.*`)).(*regexp.Regexp)
	fmt.Println(twoBrace.FindStringSubmatch(`sss{{{<li>sss</li>}}}sss`))
	fmt.Println(jen.Lit("11").Render(os.Stdout))
	fmt.Println(jen.Id("11").Render(os.Stdout))
	fmt.Println(jen.Lit(11).Render(os.Stdout))
	fmt.Println(jen.Lit(`'nn+1/2'`).Render(os.Stdout))
	fmt.Println(jen.Id(`"nn"+1/2`).Render(os.Stdout))

	var forReg = xerror.PanicErr(regexp.Compile(`,|\s+in\s+`)).(*regexp.Regexp)
	fmt.Println(forReg.Split(`m`, -1))
	fmt.Println(forReg.Split(`a in m`, -1))
	fmt.Println(forReg.Split(`a,b in m`, -1))
	fmt.Println(forReg.Split(`a,b,i in m`, -1))
	fmt.Println(jen.Id(strings.ToTitle("id-hello")).Render(os.Stdout))
}
