package enums

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/search"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
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
	newEnum("FacetLogic", model.FacetLogicNames()),
	newEnum("FacetKey", search.FacetNames()),
	newEnum("ContentType", model.ContentTypeNames()),
	newEnum("FileType", model.FileTypeNames()),
	newEnum("FilesStatus", model.FilesStatusNames()),
	newEnum("Language", model.LanguageValueStrings()),
	newEnum("Video3D", model.Video3DNames()),
	newEnum("VideoCodec", model.VideoCodecNames()),
	newEnum("VideoModifier", model.VideoModifierNames()),
	newEnum("VideoResolution", model.VideoResolutionNames()),
	newEnum("VideoSource", model.VideoSourceNames()),
	newEnum("TorrentContentOrderByField", search.TorrentContentOrderByNames()),
	newEnum("TorrentFilesOrderByField", search.TorrentFilesOrderByNames()),
	newEnum("QueueJobsOrderByField", search.QueueJobsOrderByNames()),
	newEnum("WorkerState", worker.StateNames()),
}
