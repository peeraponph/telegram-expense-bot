package service

import (
	"fmt"
	"regexp"
	"strconv"

	 "github.com/otiai10/gosseract/v2"

)

// ExtractAmountFromImage reads the image and extracts the highest number it finds (e.g. total price)
func ExtractAmountFromImage(path string) (int, error) {
	client := gosseract.NewClient()
	defer client.Close()

	if err := client.SetImage(path); err != nil {
		return 0, fmt.Errorf("set image failed: %v", err)
	}

	text, err := client.Text()
	if err != nil {
		return 0, fmt.Errorf("OCR failed: %v", err)
	}

	re := regexp.MustCompile(`\d+[,.]?\d*`)
	nums := re.FindAllString(text, -1)
	max := 0

	for _, s := range nums {
		s = regexp.MustCompile(`,`).ReplaceAllString(s, "")
		f, _ := strconv.ParseFloat(s, 64)
		if int(f) > max {
			max = int(f)
		}
	}

	if max == 0 {
		return 0, fmt.Errorf("no valid number found")
	}

	return max, nil
}
