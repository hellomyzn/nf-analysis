# 基本設計 — nf-analysis

## 1. 概要
Netflix の視聴履歴エクスポート（CSV）を正規化し、アーカイブ用途の `history.csv` に追記するためのツールです。1 回の実行で以下を行います。

1. `src/csv/netflix/` 配下から最新の生 CSV を読み込む。
2. 既存の `history.csv` と突き合わせて未登録レコードのみ抽出。
3. 新規レコードに連番 ID を付与して履歴に追記。

## 2. コンポーネント
| レイヤ | 役割 | 主な型 |
| ------ | ---- | ------ |
| Controller | 入力ファイルの探索とユースケースの組み立て。 | `internal/controller.NetflixController` |
| Service | 日付変換、重複排除、ID 採番などドメインロジック。 | `internal/service.NetflixService` |
| Repository | CSV の読み書きと整形。 | `internal/repository.NetflixRepository` と実装 `netflixRepositoryImpl` |
| Utility | 日付フォーマット変換。 | `internal/util.ConvertDate` |

各レイヤは依存方向が一方向（Controller → Service → Repository）になるように構成します。

## 3. 入出力仕様
### 生 CSV (`ViewingActivity.csv`)
- ヘッダー：`Title,Date,...`
- 日付フォーマット：`M/D/YY`
- 1 行目（ヘッダー）を除いて読み込む。

### 履歴 CSV (`history.csv`)
- ヘッダー：`id,date,title`
- ID は任意の接頭辞＋数字。例：`vid-0042`、`104`。
- 日付は ISO 形式（`YYYY-MM-DD`）。

## 4. 変換ルール
- 日付は `ConvertDate` で `M/D/YY` → `YYYY-MM-DD` に変換する。変換できない場合はエラーを返す。
- シグネチャ（`date + title`）で既存履歴と新規レコードを照合し、重複を除外する。
- 1 回の実行内で同一シグネチャが複数存在した場合も 1 件だけ採用する。

## 5. ID 採番
1. 履歴 CSV から既存 ID をすべて読み込み、最大値を解析する。
2. `prefix + zero-padded number` を維持したまま次の番号を算出。
3. 接頭辞が判別できない ID の場合は、数字部分のみを比較し昇順になるよう採番する。

`internal/service.newIDGenerator` がロジックを保持し、既存 ID から初期値を決定します。

## 6. CSV 書き出し
- ヘッダー行を常に書き出す。
- フィールドにカンマ・改行が含まれる場合はダブルクォートで囲む。
- ダブルクォートは `"` → `""` へエスケープする。

## 7. エラーハンドリング
- 入力 CSV が見つからない場合は Controller でエラーを返す。
- CSV 読込・書込で発生したエラーは上位に伝播し、プロセスを異常終了させる。
- 履歴ファイルが存在しない場合は空の配列として扱い、新規作成する。

## 8. 拡張ポイント
- 複数の生 CSV を連結処理する際の優先順位制御。
- タイトル正規化（シリーズ名のトリミングなど）の追加ルール。
- 視聴時間やプロフィールなど他カラムの保存。
