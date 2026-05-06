# m5paper-dashboard

M5Paper (960x540 E-Ink, 横向き) 用の荒川水系ダム貯水率ダッシュボード。

Goで画像を直接生成し、AWS Lambda + S3経由でM5Paperに配信します。

## 表示内容

- 荒川水系（関東）4ダム合計の貯水率（国土交通省関東地方整備局リアルタイム情報）
  - 4ダム合計貯水率（大きく表示）
  - 合計貯水量 / 有効容量
  - 個別ダム（二瀬・滝沢・浦山・荒川貯水池）の貯水率
  - 年間貯水率グラフ（今年＋過去年比較。過去データは実行を重ねて自動蓄積）

## アーキテクチャ

```
EventBridge (10分間隔) → Lambda (Go) → S3 → CloudFront → M5Paper
                            │
                     国土交通省関東地方整備局
                     荒川4ダム貯水状況ページ
```

## セットアップ

### 1. 環境変数の設定

```shell
cp .env.sample .env
# .env を編集
```

### 2. ローカル実行

```shell
make run
# output.jpg が生成される
```

実行のたびに「4ダム合計の当日貯水率」が `dam_history.json` に追記されます。年間グラフは初回時点では当日 1 点のみとなり、運用を続けるほど線が伸びていきます。

### 3. Lambda デプロイ

```shell
make build-lambda
cd terraform && terraform apply
```

## 開発

```shell
# ビルド
go build ./...

# テスト
make test

# ローカルで画像生成
make run

# Lambda用バイナリのビルド (arm64)
make build-lambda
```

## 設定項目

| 環境変数 | 説明 | デフォルト |
|---------|------|----------|
| `DAM_URL` | 関東地方整備局 荒川4ダム貯水状況ページURL | 荒川4ダムページのURL |
| `DAM_CACHE_FILE` | 年間データキャッシュファイルパス | `dam_history.json` |
| `S3_BUCKET` | S3バケット名 | - |
| `S3_OBJECT_KEY` | S3オブジェクトキー | - |
| `TZ` | タイムゾーン | `Asia/Tokyo` |

## 技術スタック

- Go 1.23+
- [fogleman/gg](https://github.com/fogleman/gg) - 2D画像生成
- NotoSansJP / Weather Icons / Material Design Icons - フォント（go:embed）
- 国土交通省関東地方整備局 - 荒川水系ダムリアルタイム情報
- AWS Lambda + S3 + CloudFront + EventBridge - インフラ
