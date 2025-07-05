package main

import (
	cf "github.com/konveyor/tackle2-addon-platform/cmd/cloudfoundry"
	"github.com/konveyor/tackle2-hub/api"
)

type Fetch struct {
	BaseAction
}

func (a *Fetch) Run(d *Data) (err error) {
	err = a.setApplication()
	if err != nil {
		return
	}
	var manifest *api.Manifest
	switch a.platform.Kind {
	default:
		p := cf.Provider{}
		manifest, err = p.Fetch(&a.application)
		if err != nil {
			return
		}
	}
	manifest.Application.ID = a.application.ID
	err = addon.Manifest.Create(manifest)
	return
}
