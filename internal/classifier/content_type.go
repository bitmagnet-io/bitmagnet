package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

var contentTypeSpec = json_spec.Transformer[string, model.NullContentType]{
	Typed: json_spec.Enum[string]{Values: append(model.ContentTypeNames(), "unknown")},
	Transform: func(str string, _ json_spec.ParseContext) (model.NullContentType, error) {
		if str == "unknown" {
			return model.NullContentType{}, nil
		}
		contentType, err := model.ParseContentType(str)
		if err != nil {
			return model.NullContentType{}, err
		}
		return model.NullContentType{ContentType: contentType, Valid: true}, nil
	},
}
