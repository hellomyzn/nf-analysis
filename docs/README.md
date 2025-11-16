# 📚 ドキュメント概要

nf-analysis のドキュメント一式です。Netflix の視聴履歴 CSV を正規化するための仕様、設計方針、開発ルールをまとめています。

## プロダクト概要
- Netflix からダウンロードした `ViewingActivity.csv` を `src/csv/netflix/` に配置します。
- コマンドを実行すると、新規視聴分のみを抽出し `src/csv/history.csv` に追記します。
- タイトルと視聴日をキーに重複を排除し、連番 ID を維持しながら履歴を作成します。

## ドキュメント構成
| パス | 内容 |
| ---- | ---- |
| [`design/README.md`](design/README.md) | コンポーネント設計、CSV スキーマ、ID 採番ロジックなどの基本設計。 |
| [`dev/README.md`](dev/README.md) | コーディング規約、テスト指針、CSV 配置ルールなどの開発ルール。 |
| [`sequences/overview.md`](sequences/overview.md) | 実行フロー全体のシーケンス図。 |
| [`sequences/history_merge.md`](sequences/history_merge.md) | `SaveHistory` が履歴ファイルを更新する際の詳細シーケンス。 |

## 用語
- **生 CSV（raw）**：Netflix のエクスポートファイル。`Title,Date,...` のカラムを含みます。
- **履歴 CSV（history）**：本ツールが生成する正規化済みファイル。`id,date,title` の3列で管理します。
- **シグネチャ**：`date + title` の組み合わせ。重複判定に利用します。

## 参照
- 実装コード：[`src/internal`](../src/internal)
- テストコード：[`src/test`](../src/test)
- 実行手順：リポジトリ直下の [`README.md`](../README.md)
