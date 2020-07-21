# go-remo

Nature Remo EおよびNature Remo E lite用のクライアントパッケージ。

## 使い方

```go
import "github.com/reeve0930/go-remoe"

// クライアントの作成
client := remoe.NewClient("ここにREMOのアクセストークンを記述")

// データの取得
data, err := client.GetRawData()

// 積算電力量の取得
p0 := remoe.GetPowerCunsumption(data[0])

time.Sleep(1 * time.Hour) //例えば、一時間の消費電力量

data, err := client.GetRawData()
// 2点間の消費電力量の取得
p := remoe.GetPowerCunsumptionDiff(data[0], p0)
```

## 参考資料

- [GODOC](https://godoc.org/github.com/reeve0930/go-remoe)
- [GO言語でNature Remo Eのデータを取得するパッケージを作った](https://fe-notes.work/posts/20200721_go-remoe/)
