// Package pmodel
//
// @author: xwc1125
package pmodel

import "github.com/chain5j/chain5j-pkg/codec/json"

type PhoneInfo struct {
	Ma     string    `json:"ma"`     // mac地址
	Ies    []string  `json:"ies"`    // imei
	Smis   []SimInfo `json:"smis"`   // sim信息
	IsList int       `json:"isList"` // 是否使用的List：0:false，1:true
	R      int64     `json:"r"`      // 随机码
}

type SimInfo struct {
	Is   string `json:"is"`   // imsi
	Ic   string `json:"ic"`   // iccid
	M    string `json:"m"`    // mobile
	Cn   string `json:"cn"`   // CarrierName -中国联通 4G
	Sid  int    `json:"sid"`  // 从0开始，最大为卡槽数。 0 即表示卡1的数据
	T    int    `json:"t"`    // 类别 0:高通、联发科反射数据 1:兼容的4.0，5.0的反射数据 2:5.0以上的反射数据
	Idfd bool   `json:"idfd"` // 移动网络是否首选该卡
	Idfs bool   `json:"idfs"` // 短信是否首选该卡
	R    int64  `json:"r"`    // 随机码;
}

func (a PhoneInfo) String() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}

func (a SimInfo) String() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}
