//go:build js || wasm

package main

import (
	"encoding/base64"
	"fmt"
	"syscall/js"

	"github.com/zztkm/gabigabi"
)

func main() {
	// JavaScript から呼び出せる関数を登録
	js.Global().Set("processImage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		input := args[0].String() // Base64 入力データ
		scale := args[1].Float()
		sharpen := args[2].Float()
		format := args[3].String() // "png" or "jpg"

		// 入力データをデコード
		inputBytes, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			return fmt.Sprintf("Failed to decode input: %v", err)
		}

		// ガビガビ処理
		outputBytes, err := gabigabi.ToGabigabi(inputBytes, scale, sharpen, format)
		if err != nil {
			return fmt.Sprintf("Error processing image: %v", err)
		}

		// Base64 エンコードして返す
		return base64.StdEncoding.EncodeToString(outputBytes)
	}))

	// 無限ループで終了しないようにする
	select {}
}
