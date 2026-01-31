package target

import "github.com/bitmagnet-io/bitmagnet/pkg/i18n"

func I18NMessages() []*i18n.Message {
	return []*i18n.Message{
		{
			ID:          "category",
			Description: "Label for category field in add torrent dialog",
			Other:       "Category",
		},
		{
			ID:          "stopped",
			Description: "Label for stopped field in add torrent dialog",
			Other:       "Add as stopped",
		},
	}
}
