package version

import "testing"

func TestDoesNotIncrementNonSemver(t *testing.T) {
	if incremented, err := Increment("---", "major"); err == nil {
		t.Fatalf("shouldn't have incremented a non semver version. got '%s' instead of an error.", incremented)
	}
}

func TestDoesNotIncrementUnknownSegment(t *testing.T) {
	if incremented, err := Increment("0.0.1", "something"); err == nil {
		t.Fatalf("shouldn't have incremented an unknown segment. got '%s' instead of an error.", incremented)
	}
}
func TestIncrementMajorVersion(t *testing.T) {
	expected := "1.0.0"

	if incremented, err := Increment("0.0.1", "major"); err != nil || incremented != expected {
		t.Fatalf("failed to increment the major version. got '%s' instead of '%s'. Error: %v", incremented, expected, err)
	}
}
func TestIncrementMinorVersion(t *testing.T) {
	expected := "0.1.0"
	if incremented, err := Increment("0.0.1", "minor"); err != nil || incremented != expected {
		t.Fatalf("failed to increment the minor version. got '%s' instead of '%s'. Error: %v", incremented, expected, err)
	}
}

func TestIncrementPatchVersion(t *testing.T) {
	expected := "0.0.2"
	if incremented, err := Increment("0.0.1", "patch"); err != nil || incremented != expected {
		t.Fatalf("failed to increment the patch version. got '%s' instead of '%s'. Error: %v", incremented, expected, err)
	}
}

func TestDoesNotAssembleNonSemver(t *testing.T) {
	if assembled, err := Assemble("---", "prerelease", "metadata"); err == nil {
		t.Fatalf("shouldn't have assembled a non semver version. got '%s' instead of an error.", assembled)
	}
}

func TestAddPrerelease(t *testing.T) {
	expected := "0.0.1-pre"
	if assembled, err := Assemble("0.0.1", "pre", ""); err != nil || assembled != expected {
		t.Fatalf("failed to add a prerelease. got '%s' instead of '%s'. Error: %v", assembled, expected, err)
	}
}

func TestChangePrerelease(t *testing.T) {
	expected := "0.0.1-rc"
	if assembled, err := Assemble("0.0.1-pre", "rc", ""); err != nil || assembled != expected {
		t.Fatalf("failed to change the prerelease. got '%s' instead of '%s'. Error: %v", assembled, expected, err)
	}
}

func TestDoesNotSetInvalidPrerelease(t *testing.T) {
	if assembled, err := Assemble("0.0.1-pre", "_", ""); err == nil {
		t.Fatalf("shouldn't have changed the prerelease to an invalid one. got '%s' instead of an error.", assembled)
	}
}

func TestRemovePrerelease(t *testing.T) {
	expected := "0.0.1"
	if assembled, err := Assemble("0.0.1-pre", "", ""); err != nil || assembled != expected {
		t.Fatalf("failed to remove the prerelease. got '%s' instead of '%s'. Error: %v", assembled, expected, err)
	}
}

func TestAddMetadata(t *testing.T) {
	expected := "0.0.1+md"
	if assembled, err := Assemble("0.0.1", "", "md"); err != nil || assembled != expected {
		t.Fatalf("failed to add a metadata. got '%s' instead of '%s'. Error: %v", assembled, expected, err)
	}
}

func TestChangeMetadata(t *testing.T) {
	expected := "0.0.1+md"
	if assembled, err := Assemble("0.0.1+metadata", "", "md"); err != nil || assembled != expected {
		t.Fatalf("failed to change the metadata. got '%s' instead of '%s'. Error: %v", assembled, expected, err)
	}
}

func TestDoesNotSetInvalidMetadata(t *testing.T) {
	if assembled, err := Assemble("0.0.1+md", "", "_"); err == nil {
		t.Fatalf("shouldn't have changed the metadata to an invalid one. got '%s' instead of an error.", assembled)
	}
}

func TestRemoveMetadata(t *testing.T) {
	expected := "0.0.1"
	if assembled, err := Assemble("0.0.1+md", "", ""); err != nil || assembled != expected {
		t.Fatalf("failed to remove the metadata. got '%s' instead of '%s'. Error: %v", assembled, expected, err)
	}
}

func TestShouldNotGetSegmentFromNonSemverVersion(t *testing.T) {
	if segment, err := Get("nonsemver", ""); err == nil {
		t.Fatalf("shouldn't have gotten a segment from a non semver version. got '%s' instead of an error.", segment)
	}
}
func TestShouldNotGetUnknownSegment(t *testing.T) {
	if segment, err := Get("0.0.1-pre+md", "whatever"); err == nil {
		t.Fatalf("shouldn't have gotten an unknown segment. got '%s' instead of an error.", segment)
	}
}

func TestGetNoSpecificSegment(t *testing.T) {
	version := "0.0.1-pre+md"
	if v, err := Get(version, ""); err != nil || v != version {
		t.Fatalf("failed to get the raw version. got '%s' instead of '%s'. Error: %v", v, version, err)
	}
}

func TestGetVersionMajor(t *testing.T) {
	expected := "1"
	if segment, err := Get("1.2.3-pre+md", "major"); err != nil || segment != expected {
		t.Fatalf("failed to get the major version. got '%s' instead of '%s'. Error: %v", segment, expected, err)
	}
}

func TestGetVersionMinor(t *testing.T) {
	expected := "2"
	if segment, err := Get("1.2.3-pre+md", "minor"); err != nil || segment != expected {
		t.Fatalf("failed to get the minor version. got '%s' instead of '%s'. Error: %v", segment, expected, err)
	}
}

func TestGetVersionPatch(t *testing.T) {
	expected := "3"
	if segment, err := Get("1.2.3-pre+md", "patch"); err != nil || segment != expected {
		t.Fatalf("failed to get the patch version. got '%s' instead of '%s'. Error: %v", segment, expected, err)
	}
}

func TestGetVersionPrerelease(t *testing.T) {
	expected := "pre"
	if segment, err := Get("1.2.3-pre+md", "prerelease"); err != nil || segment != expected {
		t.Fatalf("failed to get the prerelease version. got '%s' instead of '%s'. Error: %v", segment, expected, err)
	}
}

func TestGetVersionMetadata(t *testing.T) {
	expected := "md"
	if segment, err := Get("1.2.3-pre+md", "metadata"); err != nil || segment != expected {
		t.Fatalf("failed to get the metadata version. got '%s' instead of '%s'. Error: %v", segment, expected, err)
	}
}
