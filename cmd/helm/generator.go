package helm

import (
	"fmt"

	hp "github.com/konveyor/asset-generation/pkg/providers/generators/helm"
)

type Generator struct {
}

func (g *Generator) Generate(templateDir string) (files map[string]string, err error) {
	return
}

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
