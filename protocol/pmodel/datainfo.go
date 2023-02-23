// Package pmodel
//
// @author: xwc1125
package pmodel

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"

	"github.com/chain5j/chain5j-pkg/codec/json"
	"github.com/chain5j/chain5j-pkg/util/dateutil"
	"github.com/chain5j/chain5j-pkg/util/reflectutil"
)

type DataInfo map[string]interface{}

func NewDataInfo() *DataInfo {
	dataInfos := &DataInfo{}
	dataInfos.Put("t", dateutil.CurrentTime())
	return dataInfos
}

func (i DataInfo) String() string {
	bytes, _ := json.Marshal(i)
	return string(bytes)
}

func (i DataInfo) StringKV() string {
	var buffer bytes.Buffer
	for _, v := range i.Sort() {
		buffer.WriteString(v.Key)
		buffer.WriteString("=")
		buffer.WriteString(fmt.Sprintf("%s", v.Value))
		buffer.WriteString(",")
	}
	return buffer.String()[:buffer.Len()-1]
}

// 增加或者修改一个元素
func (i DataInfo) Put(k string, v interface{}) {
	i[k] = v
}

func (i DataInfo) Get(k string) (obj interface{}, objType string, isExist bool) {
	v, cb := i[k]
	var rv interface{} = nil
	var rt = ""
	var rs = false
	if cb {
		rv = v
		rs = true
		rt = reflect.TypeOf(v).String()
	}
	return rv, rt, rs
}

func (i DataInfo) GetObj(k string) (obj interface{}) {
	v, cb := i[k]
	var rv interface{} = nil
	if cb {
		rv = v
	}

	return rv
}

func (i DataInfo) GetValue(k string) reflect.Value {
	v, cb := i[k]
	var rv interface{} = nil
	var sv reflect.Value
	if cb {
		rv = v
		sv = reflectutil.GetValue(rv)
	}

	return sv
}

func (i DataInfo) GetValues(k string) []reflect.Value {
	v, cb := i[k]
	var rv interface{} = nil
	var sv []reflect.Value
	if cb {
		rv = v
		sv = reflectutil.GetValues(rv)
	}

	return sv
}

// 判断是否包括key，如果包含key返回value的类型
func (i DataInfo) ContainsKey(k string) (bool, string) {
	v, cb := i[k]
	var rs = false
	var rt = ""
	if cb {
		rs = true
		rt = reflect.TypeOf(v).String()
	}
	return rs, rt
}

func (i DataInfo) Exist(k string) bool {
	_, cb := i[k]
	var rs = false
	if cb {
		rs = true
	}
	return rs
}

// 移除一个元素
func (i DataInfo) Remove(k string) (interface{}, bool) {
	v, cb := i[k]
	var rs = false
	var rv interface{} = nil
	if cb {
		rv = v
		rs = true
		delete(i, k)
	}
	return rv, rs
}

// 复制map用于外部遍历
func (i DataInfo) ForEach() map[string]interface{} {
	mb := map[string]interface{}{}
	for k, v := range i {
		mb[k] = v
	}
	return mb
}

// 放回现在的个数
func (i DataInfo) Size() int {
	return len(i)
}

type KV struct {
	Key   string
	Value interface{}
}

// 排序
func (i DataInfo) Sort() []KV {
	var newm []KV
	var keyArray []string
	for k, _ := range i {
		keyArray = append(keyArray, k)
	}
	sort.Strings(keyArray)
	for _, v := range keyArray {
		kv := &KV{
			Key:   v,
			Value: i[v],
		}
		newm = append(newm, *kv)
	}
	return newm
}

type KVBytes struct {
	Key   string
	Value []byte
}

type Encoder interface {
	EncodeToBytes() ([]byte, error)
}

type Decoder interface {
	DecodeFromBytes([]byte) error
}

var (
	encoderInterface = reflect.TypeOf(new(Encoder)).Elem()
	decoderInterface = reflect.TypeOf(new(Decoder)).Elem()
)

func (i *DataInfo) EncodeToBytes() ([]KVBytes, error) {
	var data []KVBytes
	var kv KVBytes
	sort2 := i.Sort()

	for _, v := range sort2 {
		fmt.Println("EncodeToBytes v", v)
		kv = KVBytes{
			Key: v.Key,
		}
		valueOf := reflect.ValueOf(v.Value)
		if valueOf.Type().Implements(encoderInterface) {
			bytes, err := valueOf.Interface().(Encoder).EncodeToBytes()
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			kv.Value = bytes
		} else {
			bytes, err := json.Marshal(v.Value)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			kv.Value = bytes
		}
		data = append(data, kv)
	}
	return data, nil
}

func (i *DataInfo) DecodeFromBytes(data []KVBytes, val interface{}) error {
	valueOf := reflect.ValueOf(val)
	for _, v := range data {
		if val != nil {
			if valueOf.Type().Implements(decoderInterface) {
				val := valueOf.Interface().(Decoder)
				err := val.DecodeFromBytes(v.Value)
				if err != nil {
					return err
				}
			} else {
				err := json.Unmarshal(v.Value, val)
				if err != nil {
					return err
				}
			}
		}
		i.Put(v.Key, val)
	}
	return nil
}
