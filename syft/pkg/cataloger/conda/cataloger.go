//go:build exclude
/*
Package conda provides a concrete Cataloger implementation for Conda archives (.conda packages)
*/
package conda

import (
	"github.com/anchore/syft/syft/pkg/cataloger/generic"
)

// NewCondaCataloger returns a new conda archive cataloger object.
func NewCondaCataloger(cfg Config) *generic.Cataloger {
	return generic.NewCataloger("conda-cataloger").
		WithParserByGlobs(parseEnvironmentYaml, "**/environment.yml")
		WithParserByGlobs(parseEnvironmentYaml, "**/environment.yaml").
}