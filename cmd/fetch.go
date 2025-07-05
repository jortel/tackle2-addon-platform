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
	if a.platform.Identity.ID != 0 {

	}
	var manifest *api.Manifest
	switch a.platform.Kind {
	default:
		p := cf.Provider{
			URL: a.platform.URL,
		}
		if a.platform.Identity.ID != 0 {
			var idPtr *api.Identity
			idPtr, err = addon.Identity.Get(a.platform.Identity.ID)
			if err != nil {
				return
			}
			p.Identity = *idPtr
		}
		manifest, err = p.Fetch(&a.application)
		if err != nil {
			return
		}
	}
	manifest.Application.ID = a.application.ID
	err = addon.Manifest.Create(manifest)
	return
}
