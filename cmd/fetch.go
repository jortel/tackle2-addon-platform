package main

import (
	cf "github.com/konveyor/tackle2-addon-platform/cmd/cloudfoundry"
	"github.com/konveyor/tackle2-hub/api"
)

type Fetch struct {
	BaseAction
}

func (r *Fetch) Run(d *Data) (err error) {
	var manifest *api.Manifest
	err = r.setApplication()
	if err != nil {
		return
	}
	switch r.platform.Kind {
	default:
		p := cf.Provider{}
		manifest, err = p.Fetch(&r.application)
		if err != nil {
			return
		}
	}
	mapi := addon.Application.Manifest(r.application.ID)
	err = mapi.Get()
	return
}
