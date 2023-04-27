// Code generated by hertz generator.

package main

import (
	"backend/biz/mq"
	"backend/biz/mw"
	"backend/cmd/dal"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	dal.Init()
	mw.InitJwt()
	mq.Init()

	h := server.Default()

	register(h)
	h.Spin()
}
