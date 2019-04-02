/*
 *
 * api.go
 * apis
 *
 * Created by lin on 2018/12/10 5:16 PM
 * Copyright © 2017-2018 PYL. All rights reserved.
 *
 */

package apis

type ApiServer interface {
	LoadMiddleware()
	RegisterRouter()
	Run(port string)
}
