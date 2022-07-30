package filetype

var AllowedType = []string{
	"image/apng",
	"image/bmp",
	"image/gif",
	"image/jpeg",
	"image/pjpeg",
	"image/png",
	"image/svg+xml",
	"image/tiff",
	"image/webp",
	"image/x-icon",
}

func IsTypeAllowed(typ string) bool {
	for _, v := range AllowedType {
		if v == typ {
			return true
		}
	}

	return false
}
