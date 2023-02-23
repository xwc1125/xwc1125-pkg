// Package protocol
//
// @author: xwc1125
package protocol

import (
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chain5j/chain5j-pkg/codec/json"
	"github.com/xwc1125/xwc1125-pkg/protocol/contextx"
	"github.com/xwc1125/xwc1125-pkg/types/response"

	"github.com/chain5j/chain5j-pkg/collection/maps/hashmap"
	"github.com/chain5j/logger"
	"github.com/xwc1125/xwc1125-pkg/database/redis"
	"github.com/xwc1125/xwc1125-pkg/protocol/pmodel"
	"github.com/xwc1125/xwc1125-pkg/utils/aesutil"
	"github.com/xwc1125/xwc1125-pkg/utils/md5util"
	"github.com/xwc1125/xwc1125-pkg/utils/rsautil/gorsa"
	"github.com/xwc1125/xwc1125-pkg/utils/tcputil"
)

// ApiAuthProtocol 解析协议
func ApiAuthProtocol(ctx contextx.Context, protocolConfig ProtocolConfig, next func(obj *RequestDataObj)) {
	// ======================公用部分=========================//
	path := ctx.Request().RequestURI
	logger.Info("request path", "path", path)
	// 过滤静态资源、login接口、首页等...不需要验证
	tcpInfo := tcputil.GetTcpInfo(ctx.Request())
	requestMap := GetRequestMap(ctx.Request())
	requestDataObj := new(RequestDataObj)
	requestDataObj.TcpInfo = tcpInfo
	requestDataObj.Map = requestMap
	if checkURL(protocolConfig, path) {
		// ctx.Values().Set(KEY_REQUEST_OBJ, requestDataObj)
		// ctx.Next()
		next(requestDataObj)
		return
	}
	logger.Info("parse protocol")
	var aesKey string
	var err error
	if protocolConfig.IsProtocol {
		obj, _, _ := requestMap.Get("rsa")
		if obj == nil {
			response.Error(ctx, 40001, "aesKey is empty", nil)
			// ctx.StopExecution()
			return
		}
		aesKey, err = gorsa.PriKeyDecrypt(obj.(string), protocolConfig.PrivateKey)
		if err != nil {
			response.Error(ctx, 40000, "decrypt error", nil)
			return
		}
	}

	if protocolConfig.IsProtocol {
		// 验签
		if !AuthBySign(protocolConfig, requestMap, getSignKey(tcpInfo.Api, aesKey)) {
			response.Error(ctx, 40012, "params is error", nil)
			return
		}
	}
	requestDataObj.AesKey = aesKey
	requestDataObj.OsType = pmodel.ParseOsTypeByName(tcpInfo.Os)
	appInfo := GetAppInfo(protocolConfig, aesKey, requestMap.GetObj("app").(string))
	requestDataObj.App = appInfo
	sdkInfo := GetSdkInfo(protocolConfig, aesKey, requestMap.GetObj("sdk").(string))
	requestDataObj.Sdk = sdkInfo
	deviceInfo := GetDeviceInfo(protocolConfig, aesKey, requestMap.GetObj("device").(string))
	requestDataObj.Device = deviceInfo
	clientInfo := GetClientInfo(protocolConfig, aesKey, requestMap.GetObj("client").(string))
	requestDataObj.Client = clientInfo
	coreDataInfo := GetCoreDataInfo(protocolConfig, aesKey, requestMap.GetObj("data").(string))
	requestDataObj.Data = coreDataInfo

	if protocolConfig.Limit.AntiBrushFlag {
		if isAntiBrush(protocolConfig, requestDataObj) {
			response.Error(ctx, 40011, "operations are too frequent", nil)
			return
		}
	}

	// ctx.Values().Set(KEY_REQUEST_OBJ, requestDataObj)
	// ctx.Next()
	next(requestDataObj)
}

func ApiNormalProtocol(ctx contextx.Context, protocolConfig ProtocolConfig, next func(obj *RequestDataObj)) {
	// ======================公用部分=========================//
	path := ctx.Request().RequestURI
	logger.Info("request path", "path", path, "method", ctx.Request().Method)
	// 过滤静态资源、login接口、首页等...不需要验证
	tcpInfo := tcputil.GetTcpInfo(ctx.Request())
	requestMap := GetRequestMap(ctx.Request())
	requestDataObj := new(RequestDataObj)
	requestDataObj.TcpInfo = tcpInfo
	requestDataObj.Map = requestMap
	if checkURL(protocolConfig, path) {
		// ctx.Values().Set(KEY_REQUEST_OBJ, requestDataObj)
		// ctx.Next()
		next(requestDataObj)
		return
	}
	logger.Info("parse protocol")
	requestDataObj.OsType = pmodel.ParseOsTypeByName(tcpInfo.Os)

	// if antiBrushFlag {
	//	if isAntiBrush(requestDataObj) {
	//		coderesult.Error(ctx, 40011, "operations are too frequent", nil)
	//		return
	//	}
	// }

	// ctx.Values().Set(KEY_REQUEST_OBJ, requestDataObj)
	// ctx.Next()
	next(requestDataObj)
}
func checkURL(protoclConfig ProtocolConfig, reqPath string) bool {
	if protoclConfig.WhiteApiList == nil || len(protoclConfig.WhiteApiList) == 0 {
		return false
	}
	for _, v := range protoclConfig.WhiteApiList {
		if reqPath == v {
			return true
		}
	}
	return false
}

// GetRequestMap 获取请求中的参数
func GetRequestMap(request *http.Request) *hashmap.HashMap {
	reqMap := hashmap.NewHashMap(true)
	body := request.Body
	s, _ := ioutil.ReadAll(body) // 把body 内容读入字符串 s
	mapStr := string(s)
	if mapStr == "" {
		query := request.URL.Query()
		for k, v := range query {
			if len(v) == 1 {
				reqMap.Put(k, v[0])
			} else {
				reqMap.Put(k, v)
			}
		}
		return reqMap
	}
	mapStr, _ = url.QueryUnescape(mapStr)
	logger.Info("请求参数", "map", mapStr)
	split := strings.Split(mapStr, "&")
	for _, v := range split {
		key := v[0:strings.Index(v, "=")]
		value := v[strings.Index(v, "=")+1:]
		reqMap.Put(key, value)
	}
	return reqMap
}

func GetAppInfo(protocolConfig ProtocolConfig, aseKey, enStr string) *pmodel.AppInfo {
	appInfo := &pmodel.AppInfo{}
	if protocolConfig.IsProtocol {
		enStr = aesutil.AesDecryptECBFromBase64(enStr, aseKey)
	}
	json.Unmarshal([]byte(enStr), appInfo)
	return appInfo
}

func GetSdkInfo(protocolConfig ProtocolConfig, aseKey, enStr string) *pmodel.SdkInfo {
	sdkInfo := &pmodel.SdkInfo{}
	if protocolConfig.IsProtocol {
		enStr = aesutil.AesDecryptECBFromBase64(enStr, aseKey)
	}
	json.Unmarshal([]byte(enStr), sdkInfo)
	return sdkInfo
}

func GetDeviceInfo(protocolConfig ProtocolConfig, aseKey, enStr string) *pmodel.DeviceInfo {
	deviceInfo := &pmodel.DeviceInfo{}
	if protocolConfig.IsProtocol {
		enStr = aesutil.AesDecryptECBFromBase64(enStr, aseKey)
	}
	json.Unmarshal([]byte(enStr), deviceInfo)
	return deviceInfo
}

func GetClientInfo(protocolConfig ProtocolConfig, aseKey, enStr string) *pmodel.ClientInfo {
	clientInfo := &pmodel.ClientInfo{}
	if protocolConfig.IsProtocol {
		enStr = aesutil.AesDecryptECBFromBase64(enStr, aseKey)
	}
	json.Unmarshal([]byte(enStr), clientInfo)
	return clientInfo
}

func GetCoreDataInfo(protocolConfig ProtocolConfig, aseKey, enStr string) *pmodel.CoreDataInfo {
	coreDataInfo := &pmodel.CoreDataInfo{}
	if protocolConfig.IsProtocol {
		enStr = aesutil.AesDecryptECBFromBase64(enStr, aseKey)
	}
	json.Unmarshal([]byte(enStr), coreDataInfo)
	return coreDataInfo
}

func AuthBySign(protocolConfig ProtocolConfig, requestMap *hashmap.HashMap, _key string) bool {
	tempMap := requestMap
	logger.Info("签名key", "key", _key)
	// 获得参数秘钥 清除删除列表中sign
	sign := tempMap.GetObj("sign")
	if sign == nil {
		logger.Info("sign 值不存在")
		return false
	}
	tempMap.Remove("sign") // 清除删除列表中sign
	if protocolConfig.FilterSignList != nil && len(protocolConfig.FilterSignList) > 0 {
		for _, v := range protocolConfig.FilterSignList {
			tempMap.Remove(v)
		}
	}
	tempSign := createSign(tempMap, _key)
	_sign := md5util.Md5([]byte(tempSign))
	if strings.EqualFold(sign.(string), _sign) {
		logger.Info("sign 校验成功")
		return true
	}
	logger.Info("sign 校验失败")
	return false
}

func createSign(requestMap *hashmap.HashMap, _key string) string {
	builder := strings.Builder{}
	if _key != "" {
		builder.WriteString(_key)
	}
	sort := requestMap.Sort()
	kv := hashmap.KV{}
	for _, v := range sort {
		if v != kv {
			value := v.Value
			if value != nil && value.(string) != "null" {
				builder.WriteString(v.Key)
				builder.WriteString("=")
				builder.WriteString(v.Value.(string))
				builder.WriteString("&")
			}
		}
	}
	return builder.String()[:builder.Len()-1]
}

func getSignKey(api string, md5key string) string {
	api = strings.ReplaceAll(api, "//", "/")
	if strings.Contains(api, "//") {
		return getSignKey(api, md5key)
	}
	return md5key + api + "?"
}

// 防刷
// true：已经存在
// false：不存在
func isAntiBrush(protocolConfig ProtocolConfig, dataObj *RequestDataObj) bool {
	// 单次访问的数据存入Redis进行URl只能单词访问。防刷
	builder := strings.Builder{}
	builder.WriteString(dataObj.TcpInfo.Api)
	builder.WriteString("_")
	if dataObj.App != nil {
		builder.WriteString(dataObj.App.String())
		builder.WriteString("_")
	}
	if dataObj.Sdk != nil {
		builder.WriteString(dataObj.Sdk.String())
		builder.WriteString("_")
	}
	if dataObj.Phone != nil {
		builder.WriteString(dataObj.Phone.String())
		builder.WriteString("_")
	}
	if dataObj.Device != nil {
		builder.WriteString(dataObj.Device.String())
		builder.WriteString("_")
	}
	if dataObj.Data != nil {
		builder.WriteString(dataObj.Data.String())
		builder.WriteString("_")
	}

	urlParams := builder.String()
	var ping = false
	_, err := redis.Cache().Redis().Ping().Result()
	if err == nil {
		ping = true
		result, err := redis.Cache().Redis().Exists(urlParams).Result()
		if err != nil {
			logger.Error("protocol redis", "err", err)
		}
		if result == 1 {
			// 已经存在
			return true
		}
	} else {
		// 使用本地缓存

		return false
	}

	interMinu := getInterMinu(dataObj.App, dataObj.Sdk)
	if interMinu > float64(protocolConfig.Limit.InterTime) {
		logger.Info("时间间隔超出范围", "interMinu", interMinu)
		return true
	}
	if ping {
		err := redis.Cache().Redis().Set(urlParams, "", time.Duration(protocolConfig.Limit.InterTime)*time.Minute).Err()
		if err != nil {
			logger.Error("protocol redis error", "err", err)
		}
	} else {
		// 本地存储

	}
	return false
}

// 获取服务器和移动端的时间间隔
func getInterMinu(appInfo *pmodel.AppInfo, sdkInfo *pmodel.SdkInfo) float64 {
	// 手机端的时间和服务器端的时间间隔不能够超过10分钟，否则访问无效，提醒用户修改时间。
	// 同时一次访问在10分钟之类不能重复访问。
	appTime := appInfo.R / 1000
	serverTime := time.Now().Unix()

	interTimstap := math.Abs(float64(appTime - serverTime)) // 间隔的秒数
	interTimstap = interTimstap / 60                        // 间隔的分钟数

	timeTemplate1 := "2006-01-02 15:04:05" // 常规类型
	appTimeStr := time.Unix(appTime, 0).Format(timeTemplate1)
	serverTimeStr := time.Unix(serverTime, 0).Format(timeTemplate1)
	logger.Info("时间戳", "时间间隔(分钟)", interTimstap, "app", appTimeStr, "server", serverTimeStr)
	return interTimstap
}
