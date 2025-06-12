package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

// PreprocessImage uses ImageMagick to improve OCR accuracy
func PreprocessImage(inputPath string) (string, error) {
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + "_cleaned.png"

	cmd := exec.Command("convert", inputPath,
		"-resize", "300%", // Upscale for better OCR
		"-colorspace", "Gray", // Convert to grayscale
		"-brightness-contrast", "15x30", // Improve contrast
		"-sharpen", "0x1", // Slight sharpening
		"-threshold", "50%", // Binarization
		"-define", "png:compression-level=9", // Reduce file size
		outputPath,
	)

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("image preprocessing failed: %v", err)
	}

	return outputPath, nil
}

// ExtractAmountFromImage reads the image and extracts the highest probable amount
func ExtractAmountFromImage(path string) (float64, error) {
	processedPath, err := PreprocessImage(path)
	if err != nil {
		return 0, err
	}
	defer os.Remove(processedPath)

	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(processedPath)
	client.SetLanguage("tha+eng")
	client.SetWhitelist("0123456789.,")

	text, err := client.Text()
	if err != nil {
		return 0, fmt.Errorf("OCR failed: %v", err)
	}

	re := regexp.MustCompile(`\d{1,4}[.,]\d{1,2}`)

	matches := re.FindAllString(text, -1)

	max := 0.0
	for _, match := range matches {
		match = strings.ReplaceAll(match, ",", "")
		fmt.Println("🔍 Found match:", match)
		if strings.HasSuffix(match, ".") {
			match += "00" // "50." → "50.00"
		}
		val, err := strconv.ParseFloat(match, 64)
		if err == nil && val > max {
			max = val
		}
	}

	if max == 0 {
		fmt.Println("🛑 Raw OCR text:", text)
		return 0, fmt.Errorf("no valid amount found from OCR")
	}

	return max, nil

}
