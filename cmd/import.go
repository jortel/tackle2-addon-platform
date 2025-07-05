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
	addon.Activity(
		"[Import] Using platform (id=%d): %s",
		a.platform.ID,
		a.platform.Name)
	var list []api.Application
	switch a.platform.Kind {
	default:
		p := cf.Provider{
			URL: a.platform.URL,
		}
		if a.platform.Identity.ID != 0 {
			p.Identity, err = addon.Identity.Get(a.platform.Identity.ID)
			if err == nil {
				addon.Activity(
					"[Import] Using credentials (id=%d): %s",
					p.Identity.ID,
					p.Identity.Name)
			} else {
				return
			}
		}
		list, err = p.Find(d.Filter)
		if err != nil {
			return
		}
	}
	addon.Activity(
		"[Import] Found %d applications.",
		len(list))
	for _, app := range list {
		nErr := addon.Application.Create(&app)
		if nErr == nil {
			addon.Activity(
				"[Import] Application: %s, created.",
				app.Name)
			continue
		}
		addon.Errorf(
			"warn",
			"[Import] Application: %s, create failed: %s",
			app.Name,
			nErr.Error())
	}
	return
}
