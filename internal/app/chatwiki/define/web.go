// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

import (
	"net/http"
)

var WebService *http.Server

const DefaultMultipartMemory = 32 << 20 // 32 MB

const (
	StatusOK                        = 0
	ErrorCodeContainsSensitiveWords = 10001
	ErrorCodeNeedLogin              = 10002
	ErrorCodeNeedNoPermissionLogin  = 10003
)
