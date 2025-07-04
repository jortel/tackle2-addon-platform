package main

import (
	"errors"

	"github.com/konveyor/tackle2-hub/api"
)

type Action interface {
	Run(*Data) error
}

func NewAction(d *Data) (a Action, err error) {
	switch d.Action {
	case "fetch":
		a = &Fetch{}
	case "import":
		a = &Import{}
	case "generate":
	default:
		err = errors.New("action not supported")
	}
	return
}

type BaseAction struct {
	application api.Application
	platform    api.Platform
}

func (r *BaseAction) setApplication() (err error) {
	defer func() {
		if err != nil {
			if errors.Is(err, &api.NotFound{}) {
				err = nil
			}
		}
	}()
	app, err := addon.Task.Application()
	if err == nil {
		r.application = *app
	} else {
		return
	}
	if app.Platform == nil {
		return
	}
	p, err := addon.Platform.Get(app.Platform.ID)
	if err == nil {
		r.platform = *p
	}
	return
}

func (r *BaseAction) setPlatform() (err error) {
	defer func() {
		if err != nil {
			if errors.Is(err, &api.NotFound{}) {
				err = nil
			}
		}
	}()
	p, err := addon.Task.Platform()
	if err == nil {
		r.platform = *p
	}
	return
}
