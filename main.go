package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mbenabda/helm-local-chart-version/pkg/chartfile"
	"github.com/mbenabda/helm-local-chart-version/pkg/version"
	"github.com/spf13/cobra"
)

// Version identifier populated via the CI/CD process.
var Version = "HEAD"

type getVersionCommand struct {
	chart   string
	segment string
	out     io.Writer
}

type setVersionCommand struct {
	chart            string
	version          string
	prerelease       string
	updatePrerelease bool
	updateMetadata   bool
	metadata         string
	out              io.Writer
}

type bumpVersionCommand struct {
	chart   string
	segment string
	out     io.Writer
}

func (c *getVersionCommand) run() error {
	chart, err := chartfile.Load(c.chart)
	if err != nil {
		return err
	}

	segment, err := version.Get(chart.Version, c.segment)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.out, segment)

	return nil
}

func (c *setVersionCommand) run() error {
	chart, err := chartfile.Load(c.chart)
	if err != nil {
		return err
	}

	var baseVersion string
	if c.version != "" {
		baseVersion = c.version
	} else {
		baseVersion = chart.Version
	}

	prerelease := c.prerelease
	if !c.updatePrerelease {
		pre, err := version.Get(baseVersion, "prerelease")
		if err != nil {
			return err
		}
		prerelease = pre
	}

	metadata := c.metadata
	if !c.updateMetadata {
		md, err := version.Get(baseVersion, "metadata")
		if err != nil {
			return err
		}
		metadata = md
	}

	finalVersion, err := version.Assemble(baseVersion, prerelease, metadata)
	if err != nil {
		return err
	}

	chart.Version = finalVersion
	return chartfile.Save(chart, c.chart)
}

func (c *bumpVersionCommand) run() error {
	chart, err := chartfile.Load(c.chart)
	if err != nil {
		return err
	}

	incrementedVersion, err := version.Increment(chart.Version, c.segment)
	if err != nil {
		return err
	}

	chart.Version = incrementedVersion

	return chartfile.Save(chart, c.chart)
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
	f.StringVarP(&gv.segment, "version-segment", "s", "", "Specific segment of the chart's version to get (major|minor|patch|prerelease|metadata)")

	cmd.MarkFlagRequired("chart")

	return cmd
}

func newSetVersionCommand(out io.Writer) *cobra.Command {
	sc := &setVersionCommand{out: out}

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Modify a local chart's version number in place",
		RunE: func(cmd *cobra.Command, args []string) error {
			sc.updatePrerelease = cmd.Flags().Lookup("prerelease") != nil
			sc.updateMetadata = cmd.Flags().Lookup("metadata") != nil
			return sc.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&sc.chart, "chart", "c", "", "Path to a local chart's root directory")
	f.StringVarP(&sc.version, "version", "v", "", "New version of the chart")
	f.StringVarP(&sc.prerelease, "prerelease", "p", "", "")
	f.StringVarP(&sc.metadata, "metadata", "m", "", "")

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
