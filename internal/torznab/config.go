package torznab

import (
	"encoding/json"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Profile struct {
	OrderBy        search.TorrentContentOrderBy
	OrderDirection search.OrderDirection
	Tags           []string
}

type Config struct {
	Hostname *string
	Profiles map[string]Profile
}

type UntypedConfig struct {
	Hostname *string
	Profiles map[string]map[string]string
}

func NewDefaultUntypedConfig() UntypedConfig {
	return UntypedConfig{}
}

type Params struct {
	fx.In
	UntypedConfig UntypedConfig
	Log           *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Profiles lazy.Lazy[*Config]
}

// lazy validation of strongly typed config
func New(p Params) Result {
	return Result{
		Profiles: lazy.New[*Config](func() (*Config, error) {
			// make sure default profile has been defined.
			_, ok := p.UntypedConfig.Profiles[ProfileDefault]
			if !ok {
				p.UntypedConfig.Profiles[ProfileDefault] = make(map[string]string, 0)
			}
			config := Config{
				Hostname: p.UntypedConfig.Hostname,
				Profiles: make(map[string]Profile, len(p.UntypedConfig.Profiles)),
			}
			for name, profile := range p.UntypedConfig.Profiles {
				orderbyRaw, ok := profile[ProfileItemOrderBy]
				if !ok {
					orderbyRaw = string(search.TorrentContentOrderByRelevance)
				}
				orderby, err := search.ParseTorrentContentOrderBy(orderbyRaw)
				if err != nil {
					return nil, err
				}
				dirRaw, ok := profile[ProfileItemOrderDirection]
				if !ok {
					dirRaw = string(search.OrderDirectionDescending)
				}
				dir, err := search.ParseOrderDirection(dirRaw)
				if err != nil {
					return nil, err
				}
				tags := make([]string, 0)
				csvTags, ok := profile[ProfileItemTags]
				if ok && len(csvTags) > 0 {
					tags = strings.Split(csvTags, ",")
				}
				config.Profiles[name] = Profile{OrderBy: orderby, OrderDirection: dir, Tags: tags}

			}
			log, _ := json.MarshalIndent(config, "", "  ")
			p.Log.Infof("torznab profiles:\n%s\n", log)
			return &config, nil
		}),
	}
}
