package main

import (
	cf "github.com/konveyor/tackle2-addon-platform/cmd/cloudfoundry"
	"github.com/konveyor/tackle2-hub/api"
)

type Import struct {
	BaseAction
}

func (r *Import) Run(d *Data) (err error) {
	var found []api.Application
	err = r.setPlatform()
	if err != nil {
		return
	}
	switch r.platform.Kind {
	default:
		p := cf.Provider{}
		found, err = p.Import(&r.platform, d.Filter)
		if err != nil {
			return
		}
	}
	for _, app := range found {
		err = addon.Application.Create(&app)
		if err != nil {
			return
		}
	}
	return
}
