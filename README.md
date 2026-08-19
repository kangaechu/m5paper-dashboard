# m5paper-dashboard

M5Paper (960x540 E-Ink, 横向き) 用の宇連ダム貯水率ダッシュボード。

Goで画像を直接生成し、AWS Lambda + S3経由でM5Paperに配信します。

## 表示内容

- 宇連ダム貯水率（水資源機構中部支社リアルタイム情報）
  - 現在の貯水率（大きく表示）
  - 貯水位・貯水量・流入量・放流量
  - 24時間の1時間ごと貯水量変化
  - 年間貯水率グラフ（今年＋過去3年比較）

## アーキテクチャ

```
EventBridge (10分間隔) → Lambda (Go) → S3 → CloudFront → M5Paper
                            │
                     水資源機構中部支社
                     リアルタイム情報
```

## セットアップ

### 1. 環境変数の設定

```shell
cp .env.sample .env
# .env を編集
```

### 2. 過去データの取得（オプション）

年間グラフに過去データを表示するため、opengov.jp から過去の貯水率データを取得できます。

```shell
go run ./cmd/fetch-history --cache dam_history.json
```

2005年〜現在までの日次貯水率データがキャッシュファイルに保存されます。

### 3. ローカル実行

```shell
make run
# output.jpg が生成される
```

### 4. Lambda デプロイ

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
| `DAM_URL` | 水資源機構データページURL | 宇連ダムのURL |
| `DAM_CACHE_FILE` | 年間データキャッシュファイルパス | `dam_history.json` |
| `S3_BUCKET` | S3バケット名 | - |
| `S3_OBJECT_KEY` | S3オブジェクトキー | - |
| `TZ` | タイムゾーン | `Asia/Tokyo` |

## 技術スタック

- Go 1.23+
- [fogleman/gg](https://github.com/fogleman/gg) - 2D画像生成
- NotoSansJP / Weather Icons / Material Design Icons - フォント（go:embed）
- 水資源機構中部支社リアルタイム情報 - ダムデータ取得
- [opengov.jp](https://opengov.jp/geo/dam-reservoir/ure/) - 過去データ取得
- AWS Lambda + S3 + CloudFront + EventBridge - インフラ
