package main

import (
	"os"
	"path"
	"strconv"

	"github.com/konveyor/tackle2-addon-platform/cmd/helm"
	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-hub/api"
)

type Generate struct {
	BaseAction
}

func (a *Generate) Run(d *Data) (err error) {
	err = a.setApplication()
	if err != nil {
		return
	}
	if a.application.Assets == nil {
		return
	}
	assetRp, err := repository.New(
		AssetDir,
		a.application.Assets,
		a.application.Identities)
	if err != nil {
		return
	}
	err = assetRp.Fetch()
	if err != nil {
		return
	}
	generators, err := a.generators()
	if err != nil {
		return
	}
	paths := []string{}
	for _, generator := range generators {
		var templateDir string
		templateDir, err = a.fetchTemplates(generator)
		if err != nil {
			return
		}
		var names []string
		names, err = a.generate(generator, templateDir)
		if err != nil {
			return
		}
		paths = append(
			paths,
			names...)
	}
	err = assetRp.Commit(paths, "Generated.")
	return
}

func (a *Generate) generate(generator *api.Generator, templateDir string) (paths []string, err error) {
	var generated map[string]string
	switch generator.Kind {
	default:
		h := helm.Generator{}
		generated, err = h.Generate(templateDir)
		if err != nil {
			return
		}
	}
	for name, content := range generated {
		var f *os.File
		f, err = os.Create(path.Join(AssetDir, name))
		if err != nil {
			return
		}
		_, err = f.Write([]byte(content))
		if err != nil {
			return
		}
		paths = append(paths, name)
	}
	return
}

func (a *Generate) fetchTemplates(g *api.Generator) (dir string, err error) {
	dir = path.Join(
		TemplateDir,
		strconv.Itoa(int(g.ID)))
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}
	var identities []api.Ref
	if g.Identity != nil {
		identities = append(identities, *g.Identity)
	}
	template, err := repository.New(
		TemplateDir,
		g.Repository,
		identities)
	if err != nil {
		return
	}
	err = template.Fetch()
	if err != nil {
		return
	}
	return
}

func (a *Generate) generators() (list []*api.Generator, err error) {
	for _, ref := range a.application.Archetypes {
		var arch *api.Archetype
		arch, err = addon.Archetype.Get(ref.ID)
		if err != nil {
			return
		}
		for _, p := range arch.Profiles {
			var g *api.Generator
			for _, ref = range p.Generators {
				g, err = addon.Generator.Get(ref.ID)
				list = append(list, g)
			}
		}
	}
	return
}
