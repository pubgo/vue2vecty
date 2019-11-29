package tests_test

import (
	"bytes"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/pubgo/g/logs"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/vue2vecty/vue2vecty"
	"github.com/pubgo/vue2vecty/xml"
	"io"
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

<c:a:b:d
                    v-for="$groceryList"
                    v-model="$todo"
					@key="$click"
					@click="$onClick"
					:key="key.id"
					class="a b c"
					data-hello="hello"
					v-focus=true
            >
</c:a:b:d>

<li></li>
<input/>

<p data-click-sss=name>0?world:"hello"></p>

</b:a:todo-item>`

	decoder := xml.NewDecoder(bytes.NewBufferString(a))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	decoder.Entity = xml.HTMLEntity

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}

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
//	a := `
//<ul>
//
//<li
//                    v-for="$groceryList"
//                    v-model="$todo"
//					@key="$click"
//					@click="$onClick"
//					:key="key.id"
//					class="a b c"
//					data-hello="hello"
//					v-focus=true
//            >
//hello
//{{$moke}}
//
//<li> {{{$world}}} </li>
//<input/>
//
//</li>
//
//
//<p v-if="$hello>0"></p>
//<p @click="a"></p>
//<p @click="a1"></p>
//<p data-click-sss="name?:'hello'"></p>
//<p data-click-sss="name>0?world:'hello'"></p>
//
//<p>hello1111 {{name>0?world:hello}}</p>
//<p>hello1111 {{"hello1111"+hello}}</p>
//
//<li> <li> <p>测试 {{hello}}</p> </li></li>
//<li> <li> <p>测试</p> </li></li>
//
//<c:a:b:Hello
//                    v-for="$groceryList"
//                    v-model="$todo"
//					@key="$click"
//					@click="$onClick"
//					:key="key.id"
//					class="a b c"
//					data-hello="hello"
//					v-focus=true
//            >
//		<li> <li> <p>测试</p> </li></li>
//		<li> <li> <p>测试 {{hhhh}}</p> </li></li>
//</c:a:b:Hello>
//
//</ul>
//`
	b:=`<div>
    <div>
        <nav class="navbar navbar-expand-md navbar-dark bg-dark fixed-top">
            <a class="navbar-brand" href="#">Navbar</a>
            <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarsExampleDefault"
                    aria-controls="navbarsExampleDefault" aria-expanded="false" aria-label="Toggle navigation">
                <span class="navbar-toggler-icon"></span>
            </button>

            <div class="collapse navbar-collapse" id="navbarsExampleDefault">
                <ul class="navbar-nav mr-auto">
                    <li class="nav-item active">
                        <a class="nav-link" href="#">Home <span class="sr-only">(current)</span></a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="#">Link</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link disabled" href="#">Disabled</a>
                    </li>
                    <li class="nav-item dropdown">
                        <a class="nav-link dropdown-toggle" href="https://example.com" id="dropdown01"
                           data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">Dropdown</a>
                        <div class="dropdown-menu" aria-labelledby="dropdown01">
                            <a class="dropdown-item" href="#">Action</a>
                            <a class="dropdown-item" href="#">Another action</a>
                            <a class="dropdown-item" href="#">Something else here</a>
                        </div>
                    </li>
                </ul>
                <form class="form-inline my-2 my-lg-0">
                    <input class="form-control mr-sm-2" type="text" placeholder="Search" aria-label="Search"/>
                    <button class="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>
                </form>
            </div>
        </nav>
        <div id="app-4">
            <ol>
                <li v-for="todo in todos">
                    {{ todo.text }}
                </li>
            </ol>
        </div>
        <div style="float: right;">
            <label>
            <textarea style="font-family: monospace;" cols="70" rows="14"
                      @input="texthandler">{vecty-field:Input}</textarea>
            </label>
        </div>

        <div id="app">
            {{ message }}
        </div>

        <div id="app-2">
            <span v-bind:title="message">鼠标悬停几秒钟查看此处动态绑定的提示信息！</span>
        </div>

        <div id="app-3">
            <p v-if="seen">现在你看到我了</p>
        </div>

        <div id="app-5">
            <p>{{ message }}</p>
            <button v-on:click="reverseMessage">反转消息</button>
            <button @click="reverseMessage">反转消息</button>
            <div id="app-5">
                <p>{{ message }}</p>
                <button v-on:click="reverseMessage">反转消息</button>
            </div>

            <div id="app-6">
                <p>{{ message }}</p>
                <input v-model="message">
            </div>

            <div id="app-7">
                <ol>
                    <!--
                      现在我们为每个 todo-item 提供 todo 对象
                      todo 对象是变量，即其内容可以是动态的。
                      我们也需要为每个组件提供一个“key”，稍后再
                      作详细解释。
                    -->
                    <c:b:todo-item
                            v-for="item in groceryList"
                            v-bind:todo="item"
                            v-bind:key="item.id"
                    >
                        <div id="app">
                            <c:app-nav></c:app-nav>
                            <c:app-view>
                                <c:app-sidebar></c:app-sidebar>
                                <c:app-content></c:app-content>
                            </c:app-view>
                        </div>

                    </c:b:todo-item>
                </ol>
            </div>

        </div>

        <div id="app-6">
            <p>{{ message }}</p>
            <input v-model="message">
        </div>
        <ol>
            <!-- 创建一个 todo-item 组件的实例 -->
            <c:todo-item></c:todo-item>
        </ol>
        <div id="app-7">
            <ol>
                <!--
                  现在我们为每个 todo-item 提供 todo 对象
                  todo 对象是变量，即其内容可以是动态的。
                  我们也需要为每个组件提供一个“key”，稍后再
                  作详细解释。
                -->
                <c:todo-item
                        v-for="item in groceryList"
                        v-bind:todo="item"
                        v-bind:key="item.id"
                ></c:todo-item>
            </ol>
        </div>
        <div id="app">
            <c:app-nav></c:app-nav>
            <c:app-view>
                <c:app-sidebar></c:app-sidebar>
                <c:app-content></c:app-content>
            </c:app-view>
        </div>
    </div>
</div>
`

	v := vue2vecty.NewTranspiler(bytes.NewBufferString(b), "github.com/pubgo/vue2vecty", "Test", "views")
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
	var twoBrace = xerror.PanicErr(regexp.Compile(`(.*){{(.+)}}(.*)`)).(*regexp.Regexp)
	logs.Debug(twoBrace.FindStringSubmatch(`sss{{<li>sss</li>}}sss`))
	logs.Debug(twoBrace.FindStringSubmatch(`sss{{}}sss`))
	logs.Debug(twoBrace.MatchString(`sss{{}}sss`))
	logs.Debug(twoBrace.FindStringSubmatch(`sss{{<li>sss</li>}}`))
	logs.Debug(twoBrace.FindStringSubmatch(`{{<li>sss</li>}}sss`))
	logs.Debug(twoBrace.FindStringSubmatch(`{{<li>sss</li>}}`))

	//fmt.Println(jen.Lit("11").Render(os.Stdout))
	//fmt.Println(jen.Id("11").Render(os.Stdout))
	//fmt.Println(jen.Lit(11).Render(os.Stdout))
	//fmt.Println(jen.Lit(`'nn+1/2'`).Render(os.Stdout))
	//fmt.Println(jen.Id(`"nn"+1/2`).Render(os.Stdout))
	//
	//var forReg = xerror.PanicErr(regexp.Compile(`,|\s+in\s+`)).(*regexp.Regexp)
	//fmt.Println(forReg.Split(`m`, -1))
	//fmt.Println(forReg.Split(`a in m`, -1))
	//fmt.Println(forReg.Split(`a,b in m`, -1))
	//fmt.Println(forReg.Split(`a,b,i in m`, -1))
	//fmt.Println(jen.Id(strings.ToTitle("id-hello")).Render(os.Stdout))
	//fmt.Println(jen.Id(`"hello"`).Render(os.Stdout))
	//
	//var ternaryReg = xerror.PanicErr(regexp.Compile(`(.+)\?(.+):(.+)`)).(*regexp.Regexp)
	//fmt.Println(ternaryReg.MatchString(`name>0?world:hello`))
	//fmt.Println(ternaryReg.FindStringSubmatch(`name>0?world:hello`))
	//
	//css, err := parser.ParseDeclarations("width: 0%;")
	//xerror.PanicM(err, "css parsing error")
	//
	//for _, dec := range css {
	//	if dec.Important {
	//		dec.Value += "!important"
	//	}
	//	jen.Lit(dec.Property).Render(os.Stdout)
	//	fmt.Println(dec.Property,dec.Value)
	//}
}
