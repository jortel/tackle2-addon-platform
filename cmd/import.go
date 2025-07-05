package main

import (
	cf "github.com/konveyor/tackle2-addon-platform/cmd/cloudfoundry"
	"github.com/konveyor/tackle2-hub/api"
)

type Import struct {
	BaseAction
}

func (a *Import) Run(d *Data) (err error) {
	err = a.setPlatform()
	if err != nil {
		return
	}
	var list []api.Application
	switch a.platform.Kind {
	default:
		p := cf.Provider{}
		list, err = p.Find(d.Filter)
		if err != nil {
			return
		}
	}
	for _, app := range list {
		err = addon.Application.Create(&app)
		if err != nil {
			return
		}
	}
	return
}
