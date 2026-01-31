package rbac

import (
	_ "embed"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

//go:embed casbin_model.conf
var embedModel string

func newCasbinEnforcer(adp persist.Adapter) (*casbin.Enforcer, error) {
	e, err := casbin.NewEnforcer()
	if err != nil {
		return nil, err
	}

	mdl, err := model.NewModelFromString(embedModel)
	if err != nil {
		return nil, err
	}

	err = e.InitWithModelAndAdapter(mdl, adp)
	if err != nil {
		return nil, err
	}

	return e, nil
}
