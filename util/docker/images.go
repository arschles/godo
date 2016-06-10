package docker

import (
	"io"
	"strings"

	dlib "github.com/fsouza/go-dockerclient"
)

// GetImages gets images matching img from dockerCl. Returns an empty slice and a non-nil error if there was an error talking to the daemon
func GetImages(dockerCl *dlib.Client, img *Image) ([]dlib.APIImages, error) {
	return dockerCl.ListImages(dlib.ListImagesOptions{
		All:    false,
		Filter: img.String(),
	})
}

// PullImageStatus represents the status of an image pull
type PullImageStatus interface {
	String() string
}

// EnsureImage ensures that image is on the docker daemon pointed to by dockerCl. If it doesn't, then it attempts to pull the image. If the image doesn't exist, calls ifNotExists before proceeding to download the image. If ifNotExists returns an error, immediately returns that error. Otherwise, returns any error pulling the image.
func EnsureImage(dockerCl *dlib.Client, image string, ifNotExists func() (io.Writer, error)) error {
	if _, err := dockerCl.InspectImage(image); err != nil {
		statusOut, err := ifNotExists()
		if err != nil {
			return err
		}
		spl := strings.Split(image, ":")
		repo := spl[0]
		tag := ""
		if len(spl) > 1 {
			tag = spl[1]
		}
		pullOpts := dlib.PullImageOptions{Repository: repo, Tag: tag, OutputStream: statusOut}
		authConf := dlib.AuthConfiguration{}
		if err := dockerCl.PullImage(pullOpts, authConf); err != nil {
			return err
		}
		return nil
	}
	return nil
}
