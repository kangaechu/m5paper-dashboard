# m5paper-dashboard

M5Paper (540x960 E-Ink) 用のダッシュボード画像生成ツール。

Goで画像を直接生成し、AWS Lambda + S3経由でM5Paperに配信します。

## 表示内容

- 天気（気象庁API + Open-Meteo）
  - 今日の天気・気温・降水確率
  - 明日・明後日の予報
  - 24時間予報（2時間ごと、風向・風速付き）
- 鉄道運行情報（Yahoo!乗換案内）
  - 山手線、都営地下鉄、東京メトロ
- 今日・明日の予定（Google Calendar API）

## アーキテクチャ

```
EventBridge (10分間隔) → Lambda (Go) → S3 → CloudFront → M5Paper
                            │
                    ┌───────┴───────┐
                    │               │
              気象庁API /        Google Calendar
              Open-Meteo /       Yahoo!乗換案内
```

## セットアップ

### 1. 環境変数の設定

```shell
cp .env.sample .env
# .env を編集
```

### 2. Google Calendar の設定

1. [GCPコンソール](https://console.cloud.google.com/)でプロジェクトを作成
2. Calendar API を有効化
3. サービスアカウントを作成し、JSONキーをダウンロード
4. Googleカレンダーの設定で、サービスアカウントのメールアドレスにカレンダーを共有
5. JSONキーをbase64エンコードして `.env` の `GOOGLE_CREDENTIALS_JSON` に設定

```shell
GOOGLE_CREDENTIALS_JSON=$(base64 < service-account.json)
```

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
| `LOCATION_CODE` | 気象庁の地域コード | `130000`（東京） |
| `LOCATION_LAT` | 緯度（時間別天気用） | `35.6895` |
| `LOCATION_LON` | 経度（時間別天気用） | `139.6917` |
| `GOOGLE_CREDENTIALS_JSON` | サービスアカウントJSONキー（base64） | - |
| `GOOGLE_CALENDAR_IDS` | カレンダーID（カンマ区切り） | `primary` |
| `TRAIN_LINES` | 監視路線（`名前:コード`のカンマ区切り） | 山手線、都営、東京メトロ全線 |
| `S3_BUCKET` | S3バケット名 | - |
| `S3_OBJECT_KEY` | S3オブジェクトキー | - |
| `TZ` | タイムゾーン | `Asia/Tokyo` |

## 技術スタック

- Go 1.23+
- [fogleman/gg](https://github.com/fogleman/gg) - 2D画像生成
- NotoSansJP / Weather Icons / Material Design Icons - フォント（go:embed）
- 気象庁API - 日別天気予報
- Open-Meteo JMA API - 時間別天気・風向風速
- Yahoo!乗換案内 - 鉄道運行情報
- Google Calendar API - 予定取得
- AWS Lambda + S3 + CloudFront + EventBridge - インフラ
- Terraform - IaC
