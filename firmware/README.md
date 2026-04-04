# M5Paper Firmware

M5Paper用のArduinoファームウェア。CloudFront/S3からダッシュボード画像を取得して表示します。

## 動作

1. WiFi接続
2. 画像URL (PNG) をHTTP GETで取得
3. E-Inkディスプレイに表示
4. ディープスリープ（デフォルト10分）
5. 1に戻る

## セットアップ

### Arduino IDE

1. Arduino IDE をインストール
2. ボードマネージャで `M5Stack` を追加
   - URL: `https://m5stack.oss-cn-shenzhen.aliyuncs.com/resource/arduino/package_m5stack_index.json`
3. ボードから `M5Paper` を選択
4. ライブラリマネージャで `M5EPD` をインストール

### 設定

`firmware.ino` の以下を編集：

```cpp
const char* WIFI_SSID     = "YOUR_WIFI_SSID";
const char* WIFI_PASSWORD = "YOUR_WIFI_PASSWORD";
const char* IMAGE_URL     = "https://your-cloudfront-domain.cloudfront.net/your-key/dashboard.png";
const int   SLEEP_MINUTES = 10;
```

### 書き込み

1. M5PaperをUSB-Cで接続
2. Arduino IDEでポートを選択
3. アップロード
