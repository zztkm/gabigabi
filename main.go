package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/disintegration/imaging"
)

func main() {
	inputPath := flag.String("in", "", "Input image file path (png or jpg)")
	outputPath := flag.String("out", "", "Output image file path (png or jpg)")
	scale := flag.Float64("scale", 0.1, "Scale factor for pixelation (e.g., 0.1 for 10% of original size)")
	sharpen := flag.Float64("sharpen", 1.0, "Sharpen intensity (e.g., 1.0 for default intensity)")
	flag.Parse()

	if *inputPath == "" || *outputPath == "" {
		fmt.Println("Usage: gabigabi -in=in.jpg -out=out.png")
		os.Exit(1)
	}

	// 入力ファイルを開く
	inputFile, err := os.Open(*inputPath)
	if err != nil {
		fmt.Printf("Failed to open input file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	// 入力画像をデコード
	img, format, err := image.Decode(inputFile)
	if err != nil {
		fmt.Printf("Failed to decode input image: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Input format: %s\n", format)

	// ここからガビガビ処理
	// ロジック参考: https://moshi-nara.com/18695/
	// 縮小処理
	smallWidth := int(float64(img.Bounds().Dx()) * *scale)
	smallHeight := int(float64(img.Bounds().Dy()) * *scale)
	smallImg := imaging.Resize(img, smallWidth, smallHeight, imaging.NearestNeighbor)

	// 縮小処理したものを元画像のサイズに戻す
	// こうすることでピクセルを荒くしたことがわかりやすくなる
	pixelatedImg := imaging.Resize(smallImg, img.Bounds().Dx(), img.Bounds().Dy(), imaging.NearestNeighbor)

	// シャープフィルタ
	// refs: https://github.com/disintegration/imaging?tab=readme-ov-file#sharpening
	sharpenedImg := imaging.Sharpen(pixelatedImg, *sharpen)

	outputFile, err := os.Create(*outputPath)
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// 出力形式を判定して保存
	outputFormat := strings.ToLower(strings.Split(*outputPath, ".")[1])
	switch outputFormat {
	case "jpg", "jpeg":
		err = jpeg.Encode(outputFile, sharpenedImg, &jpeg.Options{Quality: 90})
	case "png":
		err = png.Encode(outputFile, sharpenedImg)
	default:
		err = fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	if err != nil {
		fmt.Printf("Failed to save output image: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("gabigabi successfully!")
}
