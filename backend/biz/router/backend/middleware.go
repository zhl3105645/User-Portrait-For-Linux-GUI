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

func _rulesMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _basicbehaviorinpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _genebasicbehaviorMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _generuleMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _ruledatainpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _datasourcesMw() []app.HandlerFunc {
	// your code...
	return nil
}

func __ddmodelMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _modelMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _deletemodelMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _modelinpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _genemodelMw() []app.HandlerFunc {
	// your code...
	return nil
}

func __ddlabelMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _labelMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _deletelabelMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _labelinpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _genelabelMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _usersMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _deleteuserMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _labelsMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _treelabelsMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _profileMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _profile0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func __ddcrowdMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _crowdinpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _crowdMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _genecrowdMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _deletecrowdMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _group_profileMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _groupprofileMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _crowdsMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _singlelabelMw() []app.HandlerFunc {
	// your code...
	return nil
}
