# nf-analysis

Go 製のコマンドラインツールで、Netflix の「視聴履歴」CSV を正規化して履歴ファイルに追記します。生のエクスポートデータから重複視聴を取り除き、日付や ID を整理した `history.csv` を生成することを目的にしています。

## 主な機能
- **視聴履歴の正規化**：Netflix からエクスポートした `ViewingActivity.csv` を読み取り、日付を `YYYY-MM-DD` 形式に変換します。
- **重複排除と履歴突合**：既存の `history.csv` を読み込み、タイトル＋日付の組み合わせで照合して未登録のエピソードだけを抽出します。
- **連番 ID の自動採番**：最新の ID から連番を継続し、`vid-0001` のような接頭辞付きの形式も維持します。
- **CSV 出力の整形**：カンマや改行を含むタイトルを適切にクォートし、UTF-8 のヘッダー付き CSV を書き出します。

## ディレクトリ構成
```
.
├── README.md              # このファイル
├── docs/                  # 仕様・設計ドキュメント
├── src/
│   ├── cmd/main.go        # エントリーポイント
│   ├── internal/
│   │   ├── controller/    # 入力・出力経路を制御
│   │   ├── service/       # 重複排除や採番などのドメインロジック
│   │   ├── repository/    # CSV の読込・書込
│   │   └── util/          # 日付変換ユーティリティ
│   └── test/              # サービス / リポジトリのユニットテスト
└── infra/                 # 開発補助の Docker ファイルなど
```

## 必要要件
- Go 1.24 以降
- Netflix からエクスポートした `ViewingActivity.csv`
- 既存履歴を保存する `src/csv/history.csv`（初回は空ファイルでも可）

## 使い方
1. 必要なディレクトリとファイルを用意します。
   ```bash
   mkdir -p src/csv/netflix
   cp /path/to/ViewingActivity.csv src/csv/netflix/viewing_activity.csv
   # 履歴ファイルが無ければヘッダーのみで作成
   printf "id,date,title\n" > src/csv/history.csv
   ```
2. ツールを実行します。
   ```bash
   cd src
   go run ./cmd
   ```
3. 処理が完了すると、`src/csv/history.csv` に新しいレコードが追記されます。

## テスト
ユニットテストは Go Modules 配下で実行します。
```bash
cd src
go test ./...
```

## 補足
- 入力 CSV にはヘッダー行が必要です（Netflix 公式エクスポートのデフォルト形式）。
- 既存履歴は ID/日付/タイトルの3列構成で管理します。ID に接頭辞が含まれている場合は、最も大きい値を元に連番が継続されます。
- 変換ルールや内部構造の詳細は [`docs/`](docs/README.md) を参照してください。
