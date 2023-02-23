// Package goinject
//
// @author: xwc1125
package goinject

import "go.uber.org/fx"

type (
	Option     = fx.Option // new可选项。必须要实现对app 初始化的方法apply(*App)，要实现打印接口fmt.Stringer方法
	In         = fx.In     // 让相同的对象按照tag能够赋值到一个结构体上面
	Out        = fx.Out    // 将当前结构体的字段按名字输出
	Annotation = fx.Annotation
	Annotated  = fx.Annotated
	Hook       = fx.Hook
	Lifecycle  = fx.Lifecycle
	App        = fx.App
)

var (
	New       = fx.New
	Options   = fx.Options
	NopLogger = fx.NopLogger
	Provide   = fx.Provide  // 将被依赖的对象的构造函数传进去，传进去的函数必须是个待返回值的函数指针.func NewGirl()*Girl-->fx.Provide(NewGirl)
	Invoke    = fx.Invoke   // 将函数依赖的对象作为参数传进函数然后调用函数。invoke:= func(girl* Girl)-->fx.Invoke(invoke)
	Supply    = fx.Supply   // 直接提供被依赖的对象。提供的不能是接口。girl:=Newgirl()-->fx.Supply(girl)
	Populate  = fx.Populate // 将通过容器内的值对外面的变量进行赋值。var girl *girl-->fx.Populate(&gay)
	Annotate  = fx.Annotate // 让相同的对象按照tag能够赋值到一个结构体上面，结构体必须内嵌fx.in。type Girl struct {fx.In}
	Extract   = fx.Extract
	Replace   = fx.Replace
	Decorate  = fx.Decorate
)
