package docker

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errInvalidImageName = errors.New("invalid image name")
)

// Image represents a single image name, including all information about its registry, repository and tag
type Image struct {
	registry string
	repo     string
	name     string
	tag      string
}

// ParseImageFromName parses a raw image name string into an Image
func ParseImageFromName(name string) (*Image, error) {
	spl := strings.Split(name, "/")
	splLast := strings.Split(spl[len(spl)-1], ":")
	tag := "latest"
	if len(splLast) > 1 {
		tag = splLast[1]
		spl[len(spl)-1] = splLast[0]
	}
	if len(spl) == 1 {
		// dockerhub trusted image
		return &Image{
			registry: "",
			repo:     "",
			name:     spl[0],
			tag:      tag,
		}, nil
	} else if len(spl) == 2 {
		// dockerhub image
		return &Image{
			registry: "",
			repo:     spl[0],
			name:     spl[1],
			tag:      tag,
		}, nil
	} else if len(spl) == 3 {
		// non-dockerhub image
		return &Image{
			registry: spl[0],
			repo:     spl[1],
			name:     spl[2],
			tag:      tag,
		}, nil
	}
	return nil, errInvalidImageName
}

// FullWithoutTag returns the full image name without its tag
func (i Image) FullWithoutTag() string {
	return strings.Split(i.String(), ":")[0]
}

// String is the fmt.Stringer interface implementation. It returns the full image name and its tag
func (i Image) String() string {
	if i.registry != "" {
		return fmt.Sprintf("%s/%s/%s:%s", i.registry, i.repo, i.name, i.tag)
	} else if i.repo != "" {
		return fmt.Sprintf("%s/%s:%s", i.repo, i.name, i.tag)
	}
	return fmt.Sprintf("%s:%s", i.name, i.tag)
}
