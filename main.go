package main

import (
	"io"
	"path/filepath"

	"k8s.io/helm/pkg/proto/hapi/chart"

	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/chartutil"

	"fmt"
	"os"

	"github.com/Masterminds/semver"
)

// Version identifier populated via the CI/CD process.
var Version = "HEAD"

type setVersionCommand struct {
	chart   string
	version string
	out     io.Writer
}

type bumpVersionCommand struct {
	chart   string
	segment string
	out     io.Writer
}

func (c *setVersionCommand) run() error {
	chart, err := chartutil.Load(c.chart)
	if err != nil {
		return err
	}

	chart.Metadata.Version = c.version

	return writeChartFile(chart, c.chart)
}

func writeChartFile(c *chart.Chart, dest string) error {
	return chartutil.SaveChartfile(filepath.Join(dest, "Chart.yaml"), c.Metadata)
}

func incrementVersion(version string, segment string) (string, error) {
	v1, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}

	var v2 semver.Version
	switch segment {
	case "patch":
		v2 = v1.IncPatch()
	case "minor":
		v2 = v1.IncMinor()
	case "major":
		v2 = v1.IncMajor()
	default:
		return "", fmt.Errorf("Unknown version segment %s", segment)
	}

	return v2.String(), nil
}

func (c *bumpVersionCommand) run() error {
	chart, err := chartutil.Load(c.chart)
	if err != nil {
		return err
	}

	incrementedVersion, err := incrementVersion(chart.Metadata.Version, c.segment)
	if err != nil {
		return err
	}

	chart.Metadata.Version = incrementedVersion

	return writeChartFile(chart, c.chart)
}

func newSetVersionCommand(out io.Writer) *cobra.Command {
	sc := &setVersionCommand{out: out}

	cmd := &cobra.Command{
		Use:   "set --chart [PATH_TO_CHART_DIRECTORY] --version [version]",
		Short: "Modify a local chart's version number in place",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sc.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&sc.chart, "chart", "c", "", "Path to a local chart's root directory")
	f.StringVarP(&sc.version, "version", "v", "", "New version of the chart")

	cmd.MarkFlagRequired("chart")
	cmd.MarkFlagRequired("version")

	return cmd
}

func newBumpVersionCommand(out io.Writer) *cobra.Command {
	bc := &bumpVersionCommand{out: out}

	cmd := &cobra.Command{
		Use:   "bump --chart [PATH_TO_CHART_DIRECTORY] --version-segment (major|minor|patch)",
		Short: "Increment the desired segment of a local chart's version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bc.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&bc.chart, "chart", "c", "", "Path to a local chart's root directory")
	f.StringVarP(&bc.segment, "version-segment", "s", "", "segment of the chart's version to bump (major|minor|patch)")

	cmd.MarkFlagRequired("chart")
	cmd.MarkFlagRequired("version-segment")

	return cmd
}

func main() {
	rootCmd := &cobra.Command{
		Use:  "local-chart-version",
		Long: "Modify the version number of a local helm chart",
	}

	out := rootCmd.OutOrStdout()

	fmt.Fprintln(out, "Helm local-chart-version Plugin --", Version)
	fmt.Fprintln(out, "")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version of the local-chart-version helm plugin",
		Long:  "All software has versions. This is helm-local-chart-version's",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(out, "Helm local-chart-version Plugin --", Version)
		},
	})

	rootCmd.AddCommand(newSetVersionCommand(out))
	rootCmd.AddCommand(newBumpVersionCommand(out))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
