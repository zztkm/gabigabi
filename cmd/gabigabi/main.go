package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/zztkm/gabigabi"
)

func main() {
	inputPath := flag.String("in", "", "Input image file path (png or jpg)")
	outputPath := flag.String("out", "", "Output image file path (png or jpg)")
	scale := flag.Float64("scale", 0.1, "Scale factor for pixelation (e.g., 0.1 for 10% of original size)")
	sharpen := flag.Float64("sharpen", 1.0, "Sharpen intensity (e.g., 1.0 for default intensity)")
	flag.Parse()

	if *inputPath == "" || *outputPath == "" {
		fmt.Println("Usage: gabigabi -in input.jpg -out output.png")
		os.Exit(1)
	}

	// 入力ファイルを開く
	inputBytes, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Printf("Failed to open input file: %v\n", err)
		os.Exit(1)
	}

	// 出力形式
	outputFormat := strings.ToLower(strings.Split(*outputPath, ".")[1])

	// ガビガビに変換
	outputBytes, err := gabigabi.ToGabigabi(inputBytes, *scale, *sharpen, outputFormat)
	if err != nil {
		fmt.Printf("Failed to convert image: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(*outputPath, outputBytes, 0666)
	if err != nil {
		fmt.Printf("Failed to write output image file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("gabigabi successfully!")
}
