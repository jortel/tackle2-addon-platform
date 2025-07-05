package cloudfoundry

import (
	"fmt"
	"path/filepath"

	cf "github.com/cloudfoundry/go-cfclient/v3/config"
	cfp "github.com/konveyor/asset-generation/pkg/providers/discoverers/cloud_foundry"
	"github.com/konveyor/tackle2-hub/addon"
	"github.com/konveyor/tackle2-hub/api"
	"gopkg.in/yaml.v3"
)

type Provider struct {
	URL      string
	Identity api.Identity
}

func (p *Provider) Fetch(application *api.Application) (m *api.Manifest, err error) {
	if application.Coordinates == nil {
		return
	}
	coordinates := Coordinates{}
	err = application.Coordinates.As(&coordinates)
	if err != nil {
		return
	}
	ref := cfp.AppReference{
		SpaceName: coordinates.Space,
		AppName:   coordinates.Name,
	}
	client, err := p.client(ref.SpaceName)
	if err != nil {
		return
	}
	manifest, nErr := client.Discover(ref)
	if nErr != nil {
		err = nErr
		return
	}
	m = &api.Manifest{}
	m.Content = manifest.Content
	m.Secret = manifest.Secret
	return
}

func (p *Provider) Find(filter api.Map) (found []api.Application, err error) {
	f := Filter{}
	err = filter.As(&f)
	if err != nil {
		return
	}
	client, err := p.client(f.Spaces...)
	if err != nil {
		return
	}
	spaces, err := client.ListApps()
	if err != nil {
		return
	}
	for space, applications := range spaces {
		if !f.MatchSpace(space) {
			continue
		}
		for _, ref := range applications {
			appRef, cast := ref.(cfp.AppReference)
			if !cast {
				continue
			}
			if !f.MatchName(appRef.AppName) {
				continue
			}
			r := api.Application{}
			r.Name = appRef.AppName
			found = append(found, r)
		}
	}
	return
}

func (p *Provider) client(spaces ...string) (client *cfp.CloudFoundryProvider, err error) {
	cfConfig, err := cf.New(
		p.URL,
		cf.UserPassword(
			p.Identity.User,
			p.Identity.Password),
		cf.SkipTLSValidation())
	if err != nil {
		return
	}
	pConfig := &cfp.Config{
		CloudFoundryConfig: cfConfig,
		SpaceNames:         spaces,
	}
	client, err = cfp.New(pConfig, &addon.Log)
	if err != nil {
		return
	}
	return
}

type Coordinates struct {
	Space string `json:"space"`
	Name  string `json:"name"`
}

type Filter struct {
	Spaces []string `json:"spaces"`
	Names  []string `json:"names"`
}

func (f *Filter) MatchSpace(name string) (match bool) {
	for _, pattern := range f.Spaces {
		match, _ = filepath.Match(pattern, name)
		if match {
			break
		}
	}
	return
}

func (f *Filter) MatchName(name string) (match bool) {
	for _, pattern := range f.Names {
		match, _ = filepath.Match(pattern, name)
		if match {
			break
		}
	}
	return
}

//
//

func (p *Provider) Test() (err error) {
	provider, err := p.testProvider()
	if err != nil {
		return
	}
	spaces, err := provider.ListApps()
	if err != nil {
		return
	}
	//
	//
	for _, refs := range spaces {
		for _, ref := range refs {
			manifest, nErr := provider.Discover(ref)
			if nErr != nil {
				err = nErr
				return
			}
			s, _ := yaml.Marshal(manifest)
			fmt.Printf("%s\n", s)
		}
	}
	ref := cfp.AppReference{
		SpaceName: "space",
		AppName:   "nginx",
	}
	manifest, err := provider.Discover(ref)
	if err != nil {
		return
	}
	s, _ := yaml.Marshal(manifest)
	fmt.Printf("%s\n", s)

	return
}

func (p *Provider) testProvider() (provider *cfp.CloudFoundryProvider, err error) {
	user := "admin"
	password := "dtuqBCRms14buxCnCVy2J7g2n8GVHs"
	cfConfig, err := cf.New(
		"https://api.bosh-lite.com",
		cf.UserPassword(user, password),
		cf.SkipTLSValidation())
	if err != nil {
		return
	}
	pConfig := &cfp.Config{
		CloudFoundryConfig: cfConfig,
		SpaceNames:         []string{"space"},
	}
	provider, err = cfp.New(pConfig, &addon.Log)
	if err != nil {
		return
	}
	return
}
