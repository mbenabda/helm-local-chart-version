package chartfile

import (
	"path/filepath"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

// Save a chart's manifest file
func Save(cm *chart.Metadata, chartDir string) error {
	return chartutil.SaveChartfile(Path(chartDir), cm)
}

// Path of a chart's manifest file
func Path(chartDir string) string {
	return filepath.Join(chartDir, chartutil.ChartfileName)
}

// Load a chart's manifest file
func Load(chartDir string) (*chart.Metadata, error) {
	return chartutil.LoadChartfile(Path(chartDir))
}
