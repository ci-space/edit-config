package image

import (
	"fmt"
	"strings"
)

type Image struct {
	Name    string
	Version *Version
}

func ParseImage(img string) (*Image, error) {
	const imageMinParts = 2

	imageParts := strings.Split(img, ":")
	if len(imageParts) < imageMinParts {
		return nil, fmt.Errorf("expected image to have at least two parts")
	}

	version, err := ParseVersion(imageParts[len(imageParts)-1])
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
