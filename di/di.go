// Package dig
//
// @author: xwc1125
package di

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var (
	ErrFactoryNotFound = errors.New("factory not found")
)

type factory = func() (interface{}, error)

// Container 容器
type Container struct {
	sync.Mutex
	singletons map[string]interface{}
	factories  map[string]factory
}

// NewContainer 容器实例化
func NewContainer() *Container {
	return &Container{
		singletons: make(map[string]interface{}),
		factories:  make(map[string]factory),
	}
}

// SetSingleton 注册单例对象
func (p *Container) SetSingleton(name string, singleton interface{}) {
	p.Lock()
	p.singletons[name] = singleton
	p.Unlock()
}

// GetSingleton 获取单例对象
func (p *Container) GetSingleton(name string) interface{} {
	return p.singletons[name]
}

// GetPrototype 获取实例对象
func (p *Container) GetPrototype(name string) (interface{}, error) {
	factory, ok := p.factories[name]
	if !ok {
		return nil, fmt.Errorf(name + " factory not found")
	}
	return factory()
}

// SetPrototype 设置实例对象工厂
func (p *Container) SetPrototype(name string, factory factory) {
	p.Lock()
	p.factories[name] = factory
	p.Unlock()
}

func (p *Container) Ensures(instances ...interface{}) error {
	for _, instance := range instances {
		if err := p.Ensure(instance); err != nil {
			return err
		}
	}
	return nil
}

// Ensure 注入依赖
// 该方法扫描实例的所有export字段，并读取di标签，如果有该标签则启动注入。
// 判断di标签的类型来确定注入singleton或者prototype对象
func (p *Container) Ensure(instance interface{}) error {
	objValueOf := reflect.ValueOf(instance)
	elemType := reflect.TypeOf(instance).Elem()
	ele := reflect.ValueOf(instance).Elem()
	for i := 0; i < elemType.NumField(); i++ { // 遍历字段
		fieldType := elemType.Field(i)
		tag := fieldType.Tag.Get("di") // 获取tag
		diName := p.injectName(tag)
		if diName == "" {
			continue
		}
		var (
			diInstance interface{}
			err        error
		)
		// 通过tag名称获取已设置的对象，并复制给对象
		if p.isSingleton(tag) {
			// 判断是否为单例
			diInstance = p.GetSingleton(diName)
		}
		if p.isPrototype(tag) {
			// 是否为工厂
			diInstance, err = p.GetPrototype(diName)
		}
		if err != nil {
			return err
		}
		if diInstance == nil {
			return errors.New(diName + " dependency not found")
		}
		field := ele.Field(i)
		if field.CanSet() {
			field.Set(reflect.ValueOf(diInstance))
		} else {
			setMethod := "Set" + FirstUpper(fieldType.Name)
			objMethod := objValueOf.MethodByName(setMethod)
			if !objMethod.IsValid() {
				return errors.New(setMethod + " method not found")
			}
			param := []reflect.Value{reflect.ValueOf(diInstance)} // 构造参数
			ret := objMethod.Call(param)                          // 调用参数
			if ret != nil {
				for _, value := range ret {
					valInterf := value.Interface()
					switch valInterf.(type) {
					case error:
						if err := valInterf.(error); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// injectName 获取需要注入的依赖名称
func (p *Container) injectName(tag string) string {
	tags := strings.Split(tag, ",")
	if len(tags) == 0 {
		return ""
	}
	return tags[0]
}

// isSingleton 检测是否单例依赖
func (p *Container) isSingleton(tag string) bool {
	tags := strings.Split(tag, ",")
	for _, name := range tags {
		if name == "prototype" {
			return false
		}
	}
	return true
}

// isPrototype 检测是否实例依赖
func (p *Container) isPrototype(tag string) bool {
	tags := strings.Split(tag, ",")
	for _, name := range tags {
		if name == "prototype" {
			return true
		}
	}
	return false
}

// String 打印容器内部实例
func (p *Container) String() string {
	lines := make([]string, 0, len(p.singletons)+len(p.factories)+2)
	lines = append(lines, "singletons:")
	for name, item := range p.singletons {
		line := fmt.Sprintf("  %s: %x %s", name, &item, reflect.TypeOf(item).String())
		lines = append(lines, line)
	}
	lines = append(lines, "factories:")
	for name, item := range p.factories {
		line := fmt.Sprintf("  %s: %x %s", name, &item, reflect.TypeOf(item).String())
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
