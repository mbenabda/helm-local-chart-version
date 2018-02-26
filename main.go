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

type getVersionCommand struct {
	chart string
	out   io.Writer
}

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

func (c *getVersionCommand) run() error {
	chart, err := chartutil.Load(c.chart)
	if err != nil {
		return err
	}

	fmt.Fprint(c.out, chart.Metadata.Version)

	return nil
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

func newGetVersionCommand(out io.Writer) *cobra.Command {
	gv := &getVersionCommand{out: out}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Print a local chart's version number",
		RunE: func(cmd *cobra.Command, args []string) error {
			return gv.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&gv.chart, "chart", "c", "", "Path to a local chart's root directory")

	cmd.MarkFlagRequired("chart")

	return cmd
}

func newSetVersionCommand(out io.Writer) *cobra.Command {
	sc := &setVersionCommand{out: out}

	cmd := &cobra.Command{
		Use:   "set",
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
		Use:   "bump",
		Short: "Increment the desired segment of a local chart's version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bc.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&bc.chart, "chart", "c", "", "Path to a local chart's root directory")
	f.StringVarP(&bc.segment, "version-segment", "s", "", "Segment of the chart's version to bump (major|minor|patch)")

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

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version of the local-chart-version helm plugin",
		Long:  "All software has versions. This is helm-local-chart-version's",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(out, Version)
		},
	})

	rootCmd.AddCommand(newGetVersionCommand(out))
	rootCmd.AddCommand(newSetVersionCommand(out))
	rootCmd.AddCommand(newBumpVersionCommand(out))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
