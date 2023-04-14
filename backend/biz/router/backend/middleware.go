// Code generated by hertz generator.

package backend

import (
	"backend/biz/mw"
	"github.com/cloudwego/hertz/pkg/app"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _registerMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{mw.JwtMiddleware.LoginHandler}
}

func __pplistMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _apiMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{mw.JwtMiddleware.MiddlewareFunc()}
}

func __ccountMw() []app.HandlerFunc {
	// your code...
	return nil
}

func __dduserMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userinpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _uploadMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userdatauploadMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _componentinpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _genecomponentMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _ruleinpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func __ddruleMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _updateruleMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _ruleMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _elementinpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _deleteruleMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _elementMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _updateelementMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _deleteelementMw() []app.HandlerFunc {
	// your code...
	return nil
}
