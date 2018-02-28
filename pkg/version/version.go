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
func Assemble(version string, prerelease string, metadata string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}

	if prerelease != "" {
		withPrerelease, err := v.SetPrerelease(prerelease)
		if err != nil {
			return "", err
		}
		v = &withPrerelease
	}

	if metadata != "" {
		withMetadata, err := v.SetMetadata(metadata)
		if err != nil {
			return "", err
		}
		v = &withMetadata
	}

	return v.String(), nil
}
