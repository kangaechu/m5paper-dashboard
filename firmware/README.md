# M5Paper Firmware

M5Paper用ファームウェア。CloudFront/S3からダッシュボード画像を取得してE-Inkに表示します。

## 動作

1. WiFi接続
2. 画像URL (PNG) をHTTP GETで取得
3. E-Inkディスプレイに表示
4. ディープスリープ（デフォルト60分）
5. 1に戻る

## セットアップ

### 1. PlatformIO のインストール

```shell
# CLI
pip install platformio

# または VSCode 拡張 "PlatformIO IDE" をインストール
```

### 2. 設定

```shell
cd firmware
cp include/config.h.sample include/config.h
```

`include/config.h` を編集：

```cpp
#define WIFI_SSID     "YOUR_WIFI_SSID"
#define WIFI_PASSWORD "YOUR_WIFI_PASSWORD"
#define IMAGE_URL     "https://your-cloudfront-domain.cloudfront.net/your-key/dashboard.png"
#define SLEEP_MINUTES 60
```

### 3. ビルド・書き込み

```shell
cd firmware

# ビルド
pio run

# M5PaperをUSB-Cで接続して書き込み
pio run -t upload

# シリアルモニター
pio device monitor
```
