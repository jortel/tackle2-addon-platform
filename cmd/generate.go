package main

import (
	"os"
	"path"
	"strconv"

	"github.com/konveyor/tackle2-addon-platform/cmd/helm"
	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-hub/api"
	"github.com/konveyor/tackle2-hub/migration/json"
)

type Files = map[string]string

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
	assetRepo, err := repository.New(
		AssetDir,
		a.application.Assets,
		a.application.Identities)
	if err != nil {
		return
	}
	err = assetRepo.Fetch()
	if err != nil {
		return
	}
	generators, err := a.generators()
	if err != nil {
		return
	}
	paths := []string{}
	for _, gen := range generators {
		var templateDir string
		templateDir, err = a.fetchTemplates(gen)
		if err != nil {
			return
		}
		var names []string
		names, err = a.generate(gen, templateDir)
		if err != nil {
			return
		}
		paths = append(
			paths,
			names...)
	}
	err = assetRepo.Commit(paths, "Generated.")
	return
}

func (a *Generate) generate(gen *api.Generator, templateDir string) (paths []string, err error) {
	values, err := a.values(gen)
	if err != nil {
		return
	}
	var files Files
	switch gen.Kind {
	default:
		h := helm.Generator{}
		files, err = h.Generate(templateDir, values)
		if err != nil {
			return
		}
	}
	for name, content := range files {
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

func (a *Generate) values(gen *api.Generator) (vMap api.Map, err error) {
	v := Values{}
	for _, ref := range a.application.Tags {
		var tag *api.Tag
		tag, err = addon.Tag.Get(ref.ID)
		if err != nil {
			return
		}
		v.Tags = append(v.Tags, *tag)
	}
	mapi := addon.Application.Manifest(a.application.ID)
	manifest, err := mapi.Get()
	if err != nil {
		return
	}
	v.Manifest = *manifest
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &vMap)
	return
}

func (a *Generate) fetchTemplates(gen *api.Generator) (templateDir string, err error) {
	genId := strconv.Itoa(int(gen.ID))
	templateDir = path.Join(
		TemplateDir,
		genId)
	err = os.MkdirAll(templateDir, 0755)
	if err != nil {
		return
	}
	var identities []api.Ref
	if gen.Identity != nil {
		identities = append(identities, *gen.Identity)
	}
	template, err := repository.New(
		TemplateDir,
		gen.Repository,
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
			var gen *api.Generator
			for _, ref = range p.Generators {
				gen, err = addon.Generator.Get(ref.ID)
				list = append(list, gen)
			}
		}
	}
	return
}

type Values struct {
	Tags     []api.Tag    `json:"tags"`
	Manifest api.Manifest `json:"manifest"`
}
