package image

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int

	prefix string
}

func ParseVersion(v string) (*Version, error) {
	const minVersionParts = 3

	prefix := ""
	if strings.HasPrefix(v, "v") {
		prefix = "v"
		v = strings.TrimPrefix(v, "v")
	}

	parts := strings.Split(v, ".")
	if len(parts) != minVersionParts {
		return nil, errors.New("invalid version")
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %w", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %w", err)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %w", err)
	}

	return &Version{
		Major:  major,
		Minor:  minor,
		Patch:  patch,
		prefix: prefix,
	}, nil
}

func (v *Version) UpMajor() {
	v.Major++
}

func (v *Version) UpMinor() {
	v.Minor++
}

func (v *Version) UpPatch() {
	v.Patch++
}

func (v *Version) String() string {
	return fmt.Sprintf("%s%d.%d.%d", v.prefix, v.Major, v.Minor, v.Patch)
}
