package filetype_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/filetype"
)

func TestFileImageType(t *testing.T) {
	f, err := os.OpenFile("./fixture/images.jpeg", os.O_RDONLY, 0o444)
	if err != nil {
		t.Fatal(err)
	}

	buff := make([]byte, 512)
	if _, err = f.Read(buff); err != nil {
		t.Fatal(err)
	}

	fileCt := http.DetectContentType(buff)
	if !filetype.IsTypeAllowed(fileCt) {
		t.Fatal("expected true")
	}
}

func TestFilePdfType(t *testing.T) {
	f, err := os.OpenFile("./fixture/pdf.pdf", os.O_RDONLY, 0o444)
	if err != nil {
		t.Fatal(err)
	}

	buff := make([]byte, 512)
	if _, err = f.Read(buff); err != nil {
		t.Fatal(err)
	}

	fileCt := http.DetectContentType(buff)
	if filetype.IsTypeAllowed(fileCt) {
		t.Fatal("expected false")
	}
}

func TestFileMkvType(t *testing.T) {
	f, err := os.OpenFile("./fixture/mkv.mkv", os.O_RDONLY, 0o444)
	if err != nil {
		t.Fatal(err)
	}

	buff := make([]byte, 512)
	if _, err = f.Read(buff); err != nil {
		t.Fatal(err)
	}

	fileCt := http.DetectContentType(buff)
	if filetype.IsTypeAllowed(fileCt) {
		t.Fatal("expected false")
	}
}
