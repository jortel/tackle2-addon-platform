package main

import (
	"os"
	"path"
	"strconv"

	"github.com/konveyor/tackle2-addon-platform/cmd/helm"
	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-hub/api"
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
	addon.Activity(
		"[Gen] Fetch manifest for application (id=%d): %s",
		a.application.ID,
		a.application.Name)
	if a.application.Assets == nil {
		addon.Failed("[Gen] not asset repository defined.")
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
	assetDir := path.Join(
		AssetDir,
		a.application.Assets.Path)
	generators, err := a.generators()
	if err != nil {
		return
	}
	paths := []string{}
	for _, gen := range generators {
		addon.Activity(
			"[Gen] Using generator (id=%d): %s.",
			gen.ID,
			gen.Name)
		var templateDir string
		templateDir, err = a.fetchTemplates(gen)
		if err != nil {
			return
		}
		var names []string
		names, err = a.generate(gen, templateDir, assetDir)
		if err != nil {
			return
		}
		paths = append(
			paths,
			names...)
	}
	err = assetRepo.Commit(paths, "Generated.")
	if err != nil {
		return
	}
	return
}

func (a *Generate) generate(
	gen *api.Generator,
	templateDir string,
	assetDir string) (paths []string, err error) {
	//
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
		assetPath := path.Join(
			assetDir,
			path.Base(name))
		err = a.write(assetPath, content)
		if err == nil {
			paths = append(paths, assetPath)
		} else {
			return
		}
	}
	return
}

func (a *Generate) write(assetPath, content string) (err error) {
	f, err := os.Create(assetPath)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = f.Write([]byte(content))
	if err != nil {
		return
	}
	addon.Activity(
		"[Gen] created: %s",
		f.Name())
	return
}

func (a *Generate) values(gen *api.Generator) (values api.Map, err error) {
	tags := []api.Tag{}
	for _, ref := range a.application.Tags {
		var tag *api.Tag
		tag, err = addon.Tag.Get(ref.ID)
		if err != nil {
			return
		}
		tags = append(tags, *tag)
	}
	mapi := addon.Application.Manifest(a.application.ID)
	manifest, err := mapi.Get()
	if err != nil {
		return
	}
	values = api.Map{
		"manifest": manifest.Content,
		"tags":     tags,
	}
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
		templateDir,
		gen.Repository,
		identities)
	if err != nil {
		return
	}
	err = template.Fetch()
	if err != nil {
		return
	}
	templateDir = path.Join(
		templateDir,
		gen.Repository.Path)
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
