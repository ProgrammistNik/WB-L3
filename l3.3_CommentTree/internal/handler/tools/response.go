package tools

import "github.com/wb-go/wbf/ginext"

func SendError(ctx *ginext.Context, code int, err string) {
	ctx.JSON(code, ginext.H{
		"error": err,
	})
}

func SendSuccess(ctx *ginext.Context, code int, result interface{}) {
	ctx.JSON(code, ginext.H{
		"result": result,
	})
}
