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
	addon.Activity(
		"[Fetch] Fetch manifest for application (id=%d): %s",
		a.application.ID,
		a.application.Name)
	addon.Activity(
		"[Fetch] Using platform (id=%d): %s",
		a.platform.ID,
		a.platform.Name)
	var manifest *api.Manifest
	switch a.platform.Kind {
	default:
		p := cf.Provider{
			URL: a.platform.URL,
		}
		if a.platform.Identity.ID != 0 {
			p.Identity, err = addon.Identity.Get(a.platform.Identity.ID)
			if err == nil {
				addon.Activity(
					"[Fetch] Using identity (id=%d): %s",
					p.Identity.ID,
					p.Identity.Name)
			} else {
				return
			}
		}
		manifest, err = p.Fetch(&a.application)
		if err != nil {
			return
		}
	}
	manifest.Application.ID = a.application.ID
	err = addon.Manifest.Create(manifest)
	if err == nil {
		addon.Activity(
			"Manifest (id=%d) created.",
			manifest.ID)
	}
	return
}
