# GoImageOptimaze
画像を圧縮，リサイズ，クリッピングすることができます。  
jpg,png,gif(アニメーションgif)に対応しています。  

## インストール
```bash
$ go get github.com/so-ta/imageoptimaze
```

## 使い方
```bash
import "github.com/so-ta/imageoptimaze"
```

`imageoptimize.GenerateOriginalImage(file []byte) (OriginalImage, error)` または  
`imageoptimize.OpenFile(filePath string) (OriginalImage, error)` を使用して、  
圧縮、リサイズ、クリッピングしたい画像を読み込んでください。  
  
その後、  
`imageoptimize.ResizeAndCompress(originalImage OriginalImage, width uint, height uint, contentMode ContentMode, verticalAlignment VerticalAlignment, horizontalAlignment HorizontalAlignment) ([]byte, error)`  
とすることで、圧縮、リサイズ、クリッピングされた画像が返り値として使用できます。  

## サンプル
```go
originalImage, _ := imageoptimize.OpenFile("sample.gif")
resizedImage, _ := imageoptimize.ResizeAndCompress(originalImage, 500, 500, AspectFit, VerticalAlignmentCenter, HorizontalAlignmentCenter)
```

`contentMode` には `ScaleToFill` , `AspectFit` , `AspectFill` が指定できます。  
`VerticalAlignment` には `VerticalAlignmentTop` , `VerticalAlignmentBottom` , `VerticalAlignmentCenter` が指定できます。  
`HorizontalAlignment` には `HorizontalAlignmentLeft` , `HorizontalAlignmentRight` , `HorizontalAlignmentCenter` が指定できます。  

## おまけ
リサイズ後ファイルを保存したい  
```go
originalImage, _ := imageoptimize.OpenFile("sample.png")
resizedImage, _ := imageoptimize.ResizeAndCompress(originalImage, 500, 500, AspectFit, VerticalAlignmentCenter, HorizontalAlignmentCenter)
file, _ := os.Create(`sample-resized.png`)
defer file.Close()
file.Write(resizedImage)
```



## その他
以下を参考にしました。  
https://github.com/nfnt/resize  
https://qiita.com/from_Unknown/items/40d1947292c53fe7ea74
