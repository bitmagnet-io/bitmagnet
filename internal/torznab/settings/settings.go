package settings

import (
	"encoding/json"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Profile struct {
	OrderBy        search.TorrentContentOrderBy
	OrderDirection search.OrderDirection
	Tags           []string
}

type Settings struct {
	Hostname *string
	Profiles map[string]Profile
}

type Params struct {
	fx.In
	Config torznab.Config
	Log    *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Profiles lazy.Lazy[*Settings]
}

// lazy validation of settings
func New(p Params) Result {
	return Result{
		Profiles: lazy.New[*Settings](func() (*Settings, error) {
			// make sure default profile has been defined.
			_, ok := p.Config.Profiles[torznab.ProfileDefault]
			if !ok {
				p.Config.Profiles[torznab.ProfileDefault] = make(map[string]string, 0)
			}
			settings := Settings{
				Hostname: p.Config.Hostname,
				Profiles: make(map[string]Profile, len(p.Config.Profiles)),
			}
			for name, profile := range p.Config.Profiles {
				orderbyRaw, ok := profile[torznab.ProfileItemOrderBy]
				if !ok {
					orderbyRaw = string(search.TorrentContentOrderByRelevance)
				}
				orderby, err := search.ParseTorrentContentOrderBy(orderbyRaw)
				if err != nil {
					return nil, err
				}
				dirRaw, ok := profile[torznab.ProfileItemOrderDirection]
				if !ok {
					dirRaw = string(search.OrderDirectionDescending)
				}
				dir, err := search.ParseOrderDirection(dirRaw)
				if err != nil {
					return nil, err
				}
				tags := make([]string, 0)
				csvTags, ok := profile[torznab.ProfileItemTags]
				if ok && len(csvTags) > 0 {
					tags = strings.Split(csvTags, ",")
				}
				settings.Profiles[name] = Profile{OrderBy: orderby, OrderDirection: dir, Tags: tags}

			}
			log, _ := json.MarshalIndent(settings, "", "  ")
			p.Log.Infof("torznab profiles:\n%s\n", log)
			return &settings, nil
		}),
	}
}
