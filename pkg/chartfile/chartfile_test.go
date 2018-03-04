package chartfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func TestLoad(t *testing.T) {
	givenNewChart(t, "chartfile_TestLoad", "0.0.1", shouldLoadChartFile)
}

func TestPath(t *testing.T) {
	givenNewChart(t, "chartfile_TestPath", "0.0.1", shouldCalculateChartFilePath)
}
func TestSave(t *testing.T) {
	givenNewChart(t, "chartfile_TestSave", "0.0.1", shouldSaveChartFile)
}

func shouldLoadChartFile(t *testing.T, chartName string, chartVersion string, chartDir string) {
	c, err := Load(chartDir)
	if err != nil {
		t.Fatal(err)
	}
	if c.Name != chartName {
		t.Fatal("incorrect chart name")
	}
	if c.Version != chartVersion {
		t.Fatal("incorrect chart version")
	}
}

func shouldCalculateChartFilePath(t *testing.T, chartName string, chartVersion string, chartDir string) {
	path := Path(chartDir)
	if path != filepath.Join(chartDir, "Chart.yaml") {
		t.Fatal("incorrect chartfile path")
	}
}

func shouldSaveChartFile(t *testing.T, chartName string, chartVersion string, chartDir string) {
	chartFile, err := Load(chartDir)
	if err != nil {
		t.Fatal(err)
	}

	chartFile.Version = chartVersion + ".3"
	if err := Save(chartFile, chartDir); err != nil {
		t.Fatal(err)
	}

	updatedChartFile, err := Load(chartDir)
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
