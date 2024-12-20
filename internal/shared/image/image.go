package image

import (
	"fmt"
	"strings"

	versionobject "github.com/ci-space/version-object"
)

type Image struct {
	Name    string
	Version *versionobject.Version
}

func ParseImage(img string) (*Image, error) {
	const imageMinParts = 2

	imageParts := strings.Split(img, ":")
	if len(imageParts) < imageMinParts {
		return nil, fmt.Errorf("expected image to have at least two parts")
	}

	version, err := versionobject.ParseVersion(imageParts[len(imageParts)-1])
	if err != nil {
		return nil, err
	}

	return &Image{
		Name:    strings.Join(imageParts[0:len(imageParts)-1], ":"),
		Version: version,
	}, nil
}

func (i *Image) String() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Version.String())
}
