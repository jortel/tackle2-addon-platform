package cloudfoundry

import (
	"fmt"

	cf "github.com/cloudfoundry/go-cfclient/v3/config"
	cfp "github.com/konveyor/asset-generation/pkg/providers/discoverers/cloud_foundry"
	"github.com/konveyor/tackle2-hub/addon"
	"github.com/konveyor/tackle2-hub/api"
	"gopkg.in/yaml.v3"
)

type Provider struct {
}

func (p *Provider) Fetch(application *api.Application) (m *api.Manifest, err error) {
	return
}

func (p *Provider) Import(platform *api.Platform, filter api.Map) (found []api.Application, err error) {
	return
}

func (p *Provider) Test() (err error) {
	provider, err := p.provider()
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

func (p *Provider) provider() (provider *cfp.CloudFoundryProvider, err error) {
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
