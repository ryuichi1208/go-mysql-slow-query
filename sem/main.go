package main

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type person struct {
	name       string
	order, age int
}

var persons = []person{
	{name: "A", order: 1, age: 24},
	{name: "B", order: 2, age: 29},
	{name: "C", order: 3, age: 20},
	{name: "D", order: 4, age: 21},
	{name: "E", order: 5, age: 29},
	{name: "F", order: 6, age: 25},
	{name: "G", order: 7, age: 45},
	{name: "H", order: 8, age: 19},
	{name: "I", order: 9, age: 36},
	{name: "J", order: 10, age: 29},
}

// goroutine safe
func execute() error {
	ctx := context.Background()

	// 並列処理を開始
	eg, ctx := errgroup.WithContext(context.Background())
	// 同時実行できるゴルーチンを設定する。この場合は3個まで同時に並行処理走らせます
	// 4個目からは実行待ちにはいる。
	sem := semaphore.NewWeighted(3)

	for _, aPerson := range persons {
		// 無名関数にする意図は受け取るデータが変わらないようにするための実装です。
		// そうしないと通常処理前にaPersonが入れ替わってしまいます。
		func(p person) {
			// Goメソッドでgoroutine化します
			eg.Go(func() error {
				if err := sem.Acquire(ctx, 1); err != nil {
					// semaphore取得エラー
					return err
				}
				defer sem.Release(1)

				select {
				case <-ctx.Done():
					// エラーが発生した場合は後続処理をキャンセルして終了する
					println("cancel")
					return nil
				default:
					// 通常時の処理
					fmt.Printf("名前：%s 番号：%d 年齢：%d\n", p.name, p.order, p.age)
					return longProcess(p)
				}
			})
		}(aPerson)
	}

	// errgroupは全ての処理が終わるまたはエラーが返るまで 待ち合わせします
	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	}

	return nil
}

// 今まで時間のかかっていた重い処理は変更せずメソッド化した
func longProcess(p person) error {
	if p.age == 20 {
		// 45歳の場合はエラーを返す
		return fmt.Errorf("error: %s", p.name)
	}
	return nil
}

// サンプルタスク
func main() {
	fmt.Println("Start")
	execute()
	fmt.Println("End")
}
