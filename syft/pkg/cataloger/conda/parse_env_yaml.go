//go:build exclude
package conda

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"

	"github.com/anchore/syft/internal/log"
	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/pkg/cataloger/generic"
	"github.com/anchore/syft/syft/source"
)

var _ generic.Parser = parseEnvironmentYaml

// parseEnvYaml takes a Python requirements.txt file, returning all Python packages that are locked to a
// specific version.
func parseEnvironmentYaml(_ source.FileResolver, _ *generic.Environment, reader source.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	var packages []pkg.Package

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		line = trimEnvironmentYamlLine(line)

		if line == "" {
			// nothing to parse on this line
			continue
		}

		if strings.HasPrefix(line, "-e") {
			// editable packages aren't parsed (yet)
			continue
		}

		if !strings.Contains(line, "==") {
			// a package without a version, or a range (unpinned) which does not tell us
			// exactly what will be installed.
			continue
		}

		// parse a new requirement
		parts := strings.Split(line, "==")
		if len(parts) < 2 {
			// this should never happen, but just in case
			log.WithFields("path", reader.RealPath).Warnf("unable to parse environment.yaml line: %q", line)
			continue
		}

		// check if the version contains hash declarations on the same line
		version, _ := parseVersionAndHashes(parts[1])

		name := strings.TrimSpace(parts[0])
		version = strings.TrimFunc(version, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		if name == "" || version == "" {
			log.WithFields("path", reader.RealPath).Debugf("found empty package in environment.yaml line: %q", line)
			continue
		}
		packages = append(packages, newPackageForIndex(name, version, reader.Location))
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to parse environment.yaml file: %w", err)
	}

	return packages, nil, nil
}


// trimEnvYamlLine removes content from the given requirements.txt line
// that should not be considered for parsing.
func trimEnvironmentYamlLine(line string) string {
	line = strings.TrimSpace(line)
	line = removeTrailingComment(line)
	line = removeEnvironmentMarkers(line)

	return line
}
