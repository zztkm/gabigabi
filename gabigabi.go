package gabigabi

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
)

// ToGabigabi は入力画像をガビガビ処理します。
func ToGabigabi(input []byte, scale float64, sharpen float64, format string) ([]byte, error) {
	// 入力画像をデコード
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("failed to decode input image: %v", err)
	}

	// ここからガビガビ処理
	// ロジック参考: https://moshi-nara.com/18695/
	// 縮小処理
	smallWidth := int(float64(img.Bounds().Dx()) * scale)
	smallHeight := int(float64(img.Bounds().Dy()) * scale)
	smallImg := imaging.Resize(img, smallWidth, smallHeight, imaging.NearestNeighbor)

	// 縮小処理したものを元画像のサイズに戻す
	// こうすることでピクセルを荒くしたことがわかりやすくなる
	pixelatedImg := imaging.Resize(smallImg, img.Bounds().Dx(), img.Bounds().Dy(), imaging.NearestNeighbor)

	// シャープフィルタ
	// refs: https://github.com/disintegration/imaging?tab=readme-ov-file#sharpening
	sharpenedImg := imaging.Sharpen(pixelatedImg, sharpen)

	// 出力画像をエンコード
	var outputBuffer bytes.Buffer
	switch format {
	case "png":
		err = png.Encode(&outputBuffer, sharpenedImg)
	case "jpg", "jpeg":
		err = jpeg.Encode(&outputBuffer, sharpenedImg, &jpeg.Options{Quality: 90})
	default:
		err = fmt.Errorf("unsupported output format: %s", format)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to encode output image: %v", err)
	}

	return outputBuffer.Bytes(), nil
}
