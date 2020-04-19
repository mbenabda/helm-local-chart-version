package chartfile

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func TestOpen(t *testing.T) {
	givenNewChart(t, "chartfile_TestOpen", "0.0.1", shouldOpenChartFile)
}

func TestPath(t *testing.T) {
	givenNewChart(t, "chartfile_TestPath", "0.0.1", shouldCalculateChartFilePath)
}
func TestSetVersion(t *testing.T) {
	givenNewChart(t, "chartfile_TestSetVersion", "0.0", shouldSetChartFileVersion)
}

func TestDoesNotOverrideUnrelatedChartfieldKeys(t *testing.T) {
	chartName := "chartfile_TestDoesNotOverrideUnrelatedChartfieldKeys"
	chartVersion := "0.0.1"

	chartDir, err := ioutil.TempDir("", chartName)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(chartDir)

	contents, err := yaml.Marshal(map[string]string{
		"name":           chartName,
		"version":        chartVersion,
		"someOtherField": "aValue",
	})

	err = ioutil.WriteFile(Path(chartDir), contents, 0777)
	if err != nil {
		t.Fatal(err)
	}

	c, err := Open(chartDir)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.SetVersion(chartVersion); err != nil {
		t.Fatal(err)
	}

	read, err := ioutil.ReadFile(Path(chartDir))
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(read, contents) {
		t.Fatal(fmt.Errorf("SetVersion modified the contents of the charfile. expected\n %s\n got\n %s", contents, read))
	}
}

func shouldOpenChartFile(t *testing.T, chartName string, chartVersion string, chartDir string) {
	c, err := Open(chartDir)
	if err != nil {
		t.Fatal(err)
	}
	if c.Version() != chartVersion {
		t.Fatal("incorrect chart version")
	}
}

func shouldCalculateChartFilePath(t *testing.T, chartName string, chartVersion string, chartDir string) {
	path := Path(chartDir)
	if path != filepath.Join(chartDir, "Chart.yaml") {
		t.Fatal("incorrect chartfile path")
	}
}

func shouldSetChartFileVersion(t *testing.T, chartName string, chartVersion string, chartDir string) {
	chartFile, err := Open(chartDir)
	if err != nil {
		t.Fatal(err)
	}

	if err := chartFile.SetVersion(chartVersion + ".3"); err != nil {
		t.Fatal(err)
	}

	updatedChartFile, err := Open(chartDir)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chartFile, updatedChartFile) {
		t.Fatal("saved chart file does not match updated chart file")
	}
}

func givenNewChart(t *testing.T, chartName string, chartVersion string, testFunc func(*testing.T, string, string, string)) {
	dir, err := ioutil.TempDir("", chartName)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	chartfile := &chart.Metadata{
		Name:    chartName,
		Version: chartVersion,
	}

	chartDir, err := chartutil.Create(chartfile, dir)
	if err != nil {
		t.Fatal(err)
	}

	testFunc(t, chartName, chartVersion, chartDir)
}
