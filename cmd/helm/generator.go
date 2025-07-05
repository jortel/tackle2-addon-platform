package helm

import (
	"fmt"

	hp "github.com/konveyor/asset-generation/pkg/providers/generators/helm"
	"github.com/konveyor/tackle2-hub/api"
)

type Files = map[string]string

type Generator struct {
}

func (g *Generator) Generate(templateDir string, values api.Map) (files Files, err error) {
	files = make(Files)
	config := hp.Config{
		ChartPath: templateDir,
		Values:    values,
	}
	provider := hp.New(config)
	files, err = provider.Generate()
	if err != nil {
		return
	}
	return
}

//
//

func (g *Generator) Test() (err error) {
	config := hp.Config{
		ChartPath: "/tmp/asset-generation/pkg/providers/generators/helm/test_data/k8s_only",
		Values: map[string]any{
			"foo": map[string]any{
				"bar": "baz",
			},
		},
	}
	provider := hp.New(config)
	files, err := provider.Generate()
	if err != nil {
		return
	}
	for name, content := range files {
		fmt.Printf("%s\n", name)
		fmt.Printf("%s\n", content)
	}
	return
}
