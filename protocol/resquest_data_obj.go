// Package protocol
//
// @author: xwc1125
package protocol

import (
	"github.com/chain5j/chain5j-pkg/codec/json"

	"github.com/chain5j/chain5j-pkg/collection/maps/hashmap"
	"github.com/xwc1125/xwc1125-pkg/protocol/pmodel"
	"github.com/xwc1125/xwc1125-pkg/utils/tcputil"
)

type RequestDataObj struct {
	TcpInfo *tcputil.TcpInfo
	App     *pmodel.AppInfo
	Sdk     *pmodel.SdkInfo
	Phone   *pmodel.PhoneInfo
	Device  *pmodel.DeviceInfo
	Client  *pmodel.ClientInfo
	Data    *pmodel.CoreDataInfo
	Map     *hashmap.HashMap
	OsType  pmodel.ClientOsType
	AesKey  string
}

func (data RequestDataObj) String() string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}
