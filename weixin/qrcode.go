package weixin

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tuotoo/qrcode"
)

func PrintQrcode(matrix *qrcode.Matrix) error {
	w := bufio.NewWriterSize(os.Stdout, 1024)

	reset := "\033[0m"
	black := "\033[30;40m"
	white := "\033[30;47m"

	height := matrix.Size.Max.Y
	width := matrix.Size.Max.X
	line := white + fmt.Sprintf("%*s", width*2+2, "") + reset + "\n"

	fmt.Fprint(w, line)
	for y := 0; y < height; y++ {
		fmt.Fprint(w, white, " ")
		color_prev := white
		for x := 0; x < width; x++ {
			if matrix.Points[y][x] {
				if color_prev != black {
					fmt.Fprint(w, black)
					color_prev = black
				}
			} else {
				if color_prev != white {
					fmt.Fprint(w, white)
					color_prev = white
				}
			}
			fmt.Fprint(w, "  ")
		}
		fmt.Fprint(w, white, " ", reset, "\n")
		w.Flush()
	}
	fmt.Fprint(w, line)
	w.Flush()
	return nil
}
