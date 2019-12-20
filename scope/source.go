package scope

import (
	"strings"

	"gocv.io/x/gocv"
)

// Source is a wrapper around video source that can cleanly handle displaying images.
type Source struct {
	src         string
	video       *gocv.VideoCapture
	isImage     bool
	gotFirstImg bool
}

// OpenVideoCapture opens a source
func OpenVideoCapture(src string) (s *Source, err error) {
	s = &Source{src: src}

	if strings.Contains(src, ".jpg") {
		s.isImage = true
		s.gotFirstImg = false
	}
	s.video, err = gocv.OpenVideoCapture(s.src)
	return
}

// Read video
func (s *Source) Read(img *gocv.Mat) bool {
	if !s.isImage || !s.gotFirstImg {
		if ok := s.video.Read(img); !ok {
			return false
		}
		s.gotFirstImg = true
	}
	return true
}

// Close video
func (s *Source) Close() error {
	return s.video.Close()
}
