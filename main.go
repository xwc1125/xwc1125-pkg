// Package xwc1125_pkg
//
// @author: xwc1125
package main

import (
	_ "github.com/xwc1125/xwc1125-pkg/base"
	_ "github.com/xwc1125/xwc1125-pkg/cache/cachefree"
	_ "github.com/xwc1125/xwc1125-pkg/captcha"
	_ "github.com/xwc1125/xwc1125-pkg/copier"
	_ "github.com/xwc1125/xwc1125-pkg/database/db_gorm"
	_ "github.com/xwc1125/xwc1125-pkg/database/db_xorm"
	_ "github.com/xwc1125/xwc1125-pkg/database/es"
	_ "github.com/xwc1125/xwc1125-pkg/database/kafka"
	_ "github.com/xwc1125/xwc1125-pkg/database/redis"
	_ "github.com/xwc1125/xwc1125-pkg/database/search"
	_ "github.com/xwc1125/xwc1125-pkg/database/sqlite"
	_ "github.com/xwc1125/xwc1125-pkg/di"
	_ "github.com/xwc1125/xwc1125-pkg/email"
	_ "github.com/xwc1125/xwc1125-pkg/goinject"
	_ "github.com/xwc1125/xwc1125-pkg/hashring"
	_ "github.com/xwc1125/xwc1125-pkg/jsonp"
	_ "github.com/xwc1125/xwc1125-pkg/jwtauth"
	_ "github.com/xwc1125/xwc1125-pkg/language"
	_ "github.com/xwc1125/xwc1125-pkg/ldap"
	_ "github.com/xwc1125/xwc1125-pkg/lock/mock"
	_ "github.com/xwc1125/xwc1125-pkg/lock/redis"
	_ "github.com/xwc1125/xwc1125-pkg/middleware/circuitbreaker"
	_ "github.com/xwc1125/xwc1125-pkg/middleware/config"
	_ "github.com/xwc1125/xwc1125-pkg/middleware/limiter"
	_ "github.com/xwc1125/xwc1125-pkg/middleware/metrics/promotheus"
	_ "github.com/xwc1125/xwc1125-pkg/middleware/registry/consul"
	_ "github.com/xwc1125/xwc1125-pkg/middleware/registry/etcd"
	_ "github.com/xwc1125/xwc1125-pkg/middleware/registry/nacos"
	_ "github.com/xwc1125/xwc1125-pkg/middleware/shutdown"
	_ "github.com/xwc1125/xwc1125-pkg/middleware/tracer/jaeger"
	_ "github.com/xwc1125/xwc1125-pkg/middleware/tracer/plugins"
	_ "github.com/xwc1125/xwc1125-pkg/middleware/tracer/provider"
	_ "github.com/xwc1125/xwc1125-pkg/otp"
	_ "github.com/xwc1125/xwc1125-pkg/pool/groupwork"
	_ "github.com/xwc1125/xwc1125-pkg/pprof"
	_ "github.com/xwc1125/xwc1125-pkg/protocol"
	_ "github.com/xwc1125/xwc1125-pkg/protocol/contextx"
	_ "github.com/xwc1125/xwc1125-pkg/rbac"
	_ "github.com/xwc1125/xwc1125-pkg/resourcetree"
	_ "github.com/xwc1125/xwc1125-pkg/secure/checkurl"
	_ "github.com/xwc1125/xwc1125-pkg/secure/ddm"
	_ "github.com/xwc1125/xwc1125-pkg/secure/ipfilter"
	_ "github.com/xwc1125/xwc1125-pkg/secure/password"
	_ "github.com/xwc1125/xwc1125-pkg/secure/safecurl"
	_ "github.com/xwc1125/xwc1125-pkg/secure/shift"
	_ "github.com/xwc1125/xwc1125-pkg/secure/ssrf"
	_ "github.com/xwc1125/xwc1125-pkg/snowflake"
	_ "github.com/xwc1125/xwc1125-pkg/ssh/scp"
	_ "github.com/xwc1125/xwc1125-pkg/ssh/sftpx"
	_ "github.com/xwc1125/xwc1125-pkg/types"
	_ "github.com/xwc1125/xwc1125-pkg/utils/aesutil"
	_ "github.com/xwc1125/xwc1125-pkg/utils/base64util"
	_ "github.com/xwc1125/xwc1125-pkg/utils/excelutil"
	_ "github.com/xwc1125/xwc1125-pkg/utils/fileutil"
	_ "github.com/xwc1125/xwc1125-pkg/utils/iputil"
	_ "github.com/xwc1125/xwc1125-pkg/utils/jsonutil"
	_ "github.com/xwc1125/xwc1125-pkg/utils/kvutil"
	_ "github.com/xwc1125/xwc1125-pkg/utils/md5util"
	_ "github.com/xwc1125/xwc1125-pkg/utils/randutil"
	_ "github.com/xwc1125/xwc1125-pkg/utils/reflectutil"
	_ "github.com/xwc1125/xwc1125-pkg/utils/regexutil"
	_ "github.com/xwc1125/xwc1125-pkg/utils/rsautil/gorsa"
	_ "github.com/xwc1125/xwc1125-pkg/utils/stringutil"
	_ "github.com/xwc1125/xwc1125-pkg/utils/tcputil"
	_ "github.com/xwc1125/xwc1125-pkg/utils/ziputil"
	_ "github.com/xwc1125/xwc1125-pkg/validator"
	_ "github.com/xwc1125/xwc1125-pkg/version"
	_ "github.com/xwc1125/xwc1125-pkg/watermark"
)

func main() {

}
