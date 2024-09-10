package conv

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"spotify-stats-cli/cnsts"
	"strconv"

	"github.com/nfnt/resize"
)

func RGBtoAnsi(r, g, b int) string {
	return "\x1b[38;2;" + strconv.FormatInt(int64(r), 10) + ";" + strconv.FormatInt(int64(g), 10) + ";" + strconv.FormatInt(int64(b), 10) + "m"
}

func Convert(img image.Image) string {
	// assume image is square so aspect ratio won't be distorted
	h := cnsts.DIM * cnsts.FONT_ASPECT_RATIO
	img = resize.Resize(cnsts.DIM, uint(h), img, resize.Lanczos3)
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var outputString string
	for y := 0; y < height; y += 1 {
		for x := 0; x < width; x += 1 {
			pixel := img.At(x, y)
			r, g, b, _ := pixel.RGBA()
			r, g, b = r>>8, g>>8, b>>8 // convert to 8-bit color

			colorStr := RGBtoAnsi(int(r), int(g), int(b))

			outputString += colorStr + cnsts.CHAR
		}
		outputString += cnsts.ANSI_RESET + "\n"
	}

	return outputString
}
