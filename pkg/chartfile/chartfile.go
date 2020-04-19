package chartfile

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/Masterminds/semver"

	yamlpatch "github.com/krishicks/yaml-patch"
	"k8s.io/helm/pkg/chartutil"
)

type ChartFile struct {
	path    string
	version string
}

// Path of a chart's manifest file
func Path(chartDir string) string {
	return filepath.Join(chartDir, chartutil.ChartfileName)
}

// Open a chart's manifest
func Open(chartDir string) (*ChartFile, error) {
	path := Path(chartDir)
	metadata, err := chartutil.LoadChartfile(path)

	if err != nil {
		return nil, err
	}

	return &ChartFile{
		path:    path,
		version: metadata.Version,
	}, nil
}

// Version returns the chart version
func (c *ChartFile) Version() string {
	return c.version
}

// SetVersion patches the chartfile with the new chart version
func (c *ChartFile) SetVersion(version string) error {
	_, err := semver.NewVersion(version)
	if err != nil {
		return err
	}

	c.version = version

	patch, err := replaceOperation("/version", version)
	if err != nil {
		return err
	}
	return c.apply(patch)
}

func replaceOperation(path string, value string) (yamlpatch.Patch, error) {
	op := fmt.Sprintf(`---
- op: replace
  path: %s
  value: %s
`, path, value)

	ops := []byte(op)

	return yamlpatch.DecodePatch(ops)
}

func (c *ChartFile) apply(patch yamlpatch.Patch) error {
	doc, err := ioutil.ReadFile(c.path)
	if err != nil {
		return err
	}

	bs, err := patch.Apply(doc)
	if err != nil {
		return fmt.Errorf("applying patch failed: %s", err)
	}

	if err := ioutil.WriteFile(c.path, bs, 0777); err != nil {
		return err
	}

	return nil
}
