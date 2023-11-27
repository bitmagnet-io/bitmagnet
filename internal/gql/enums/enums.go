package enums

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type enum struct {
	Name   string
	Values []string
}

func newEnum(name string, values []string) enum {
	return enum{
		Name:   name,
		Values: values,
	}
}

var Enums = []enum{
	newEnum("ContentType", model.ContentTypeNames()),
	newEnum("FacetLogic", model.FacetLogicNames()),
	newEnum("FileType", model.FileTypeNames()),
	newEnum("FilesStatus", model.FilesStatusNames()),
	newEnum("Language", model.LanguageValueStrings()),
	newEnum("Video3d", model.Video3dNames()),
	newEnum("VideoCodec", model.VideoCodecNames()),
	newEnum("VideoModifier", model.VideoModifierNames()),
	newEnum("VideoResolution", model.VideoResolutionNames()),
	newEnum("VideoSource", model.VideoSourceNames()),
}
