package version

import (
	"fmt"

	"github.com/Masterminds/semver"
)

// Increment a version's semver segment
func Increment(version string, segment string) (string, error) {
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

// Assemble a semver version from its parts
func Assemble(version string, prerelease string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}

	withPrerelease, err := v.SetPrerelease(prerelease)
	if err != nil {
		return "", err
	}
	v = &withPrerelease

	return v.String(), nil
}

// Get a specific segment of a semver version
func Get(version string, segment string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}

	switch segment {
	case "":
		return v.String(), nil
	case "major":
		return fmt.Sprintf("%v", v.Major()), nil
	case "minor":
		return fmt.Sprintf("%v", v.Minor()), nil
	case "patch":
		return fmt.Sprintf("%v", v.Patch()), nil
	case "prerelease":
		return fmt.Sprintf("%s", v.Prerelease()), nil
	default:
		return "", fmt.Errorf("Unknown version segment %s", segment)
	}
}
