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
  
その後、以下のメソッドで、圧縮、リサイズ、クリッピングされた画像が返り値として使用できます。
### 画像を特定のサイズにクリッピング
`width` , `height` の画像を作成します。  
`contentMode` を `AspectFit` にした場合、余白ができることがあります。  
詳しい指定は引数の項に記述されています。  
`imageoptimize.ResizeAndCompress(originalImage OriginalImage, width uint, height uint, contentMode ContentMode, verticalAlignment VerticalAlignment, horizontalAlignment HorizontalAlignment) ([]byte, error)`

### 画像を拡大、縮小する
`maxWidth` , `maxHeight` いずれかの最大幅を指定して画像を作成します。  
`maxWidth` を 0に指定すると、縦幅を基準に拡大縮小されます。  
`maxHeight` を 0に指定すると、横幅を基準に拡大縮小されます。  
作成される画像はアスペクト比が保たれ、余白は発生しません。  
`imageoptimize.ThumbnailAndCompress(originalImage OriginalImage, maxWidth uint, maxHeight uint) ([]byte, error)`

### 画像の容量を圧縮する
画像の容量を削減します。横幅や高さは変わりません  
`imageoptimize.Compress(originalImage OriginalImage) ([]byte, error)`

## 引数
`contentMode` は画像のトリミング方法を指定できます  
`ScaleToFill` , `AspectFit` , `AspectFill` が指定できます。  

| contentMode | 効果 |
----|:----
| ScaleToFill | 指定されたサイズにぴったりになるように画像の横幅縦幅をそれぞれ拡大縮小します。 <br> アスペクト比は保証されません |
| AspectFit | アスペクト比を保ったまま、指定されたサイズに画像全体が収まるように拡大縮小します。 <br> 元の画像と指定する横縦幅のアスペクト比が違う場合、上下左右に余白が発生します |
| AspectFill | アスペクト比を保ったまま、 <br> 拡大縮小した結果、上下左右にはみ出る部分ある場合、画像の上下左右が切れます |


`VerticalAlignment` は上下に余白が生じた際にいずれに寄せるかを指定できます。  
`VerticalAlignmentTop` , `VerticalAlignmentBottom` , `VerticalAlignmentCenter` が指定できます。
  
| VerticalAlignment | 効果 |
----|:----
| VerticalAlignmentTop | 上寄せ |
| VerticalAlignmentBottom | 下寄せ |
| VerticalAlignmentCenter | 上下中央寄せ |

`HorizontalAlignment` は左右に余白が生じた際にいずれに寄せるかを指定できます。  
`HorizontalAlignmentLeft` , `HorizontalAlignmentRight` , `HorizontalAlignmentCenter` が指定できます。
  
| VerticalAlignment | 効果 |
----|:----
| HorizontalAlignmentLeft | 左寄せ |
| HorizontalAlignmentRight | 右寄せ |
| HorizontalAlignmentCenter | 左右中央寄せ |

## サンプル
```go
originalImage, _ := imageoptimize.OpenFile("sample.gif")
resizedImage, _ := imageoptimize.ResizeAndCompress(originalImage, 500, 500, AspectFit, VerticalAlignmentCenter, HorizontalAlignmentCenter)
```

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
