#include <M5EPD.h>
#include <WiFi.h>
#include <HTTPClient.h>

// WiFi settings
const char* WIFI_SSID     = "YOUR_WIFI_SSID";
const char* WIFI_PASSWORD = "YOUR_WIFI_PASSWORD";

// Dashboard image URL (CloudFront or S3 presigned URL)
const char* IMAGE_URL = "https://your-cloudfront-domain.cloudfront.net/your-key/dashboard.png";

// Sleep interval in minutes
const int SLEEP_MINUTES = 10;

// Display
M5EPD_Canvas canvas(&M5.EPD);

void setup() {
    M5.begin();
    M5.EPD.SetRotation(0);
    M5.EPD.Clear(true);
    M5.RTC.begin();

    // Connect to WiFi
    if (!connectWiFi()) {
        Serial.println("WiFi connection failed, going to sleep");
        goToSleep();
        return;
    }

    // Fetch and display dashboard image
    if (!fetchAndDisplay()) {
        Serial.println("Failed to fetch dashboard image");
    }

    // Disconnect WiFi and go to sleep
    WiFi.disconnect(true);
    WiFi.mode(WIFI_OFF);
    goToSleep();
}

void loop() {
    // Not reached - device sleeps after setup
}

bool connectWiFi() {
    Serial.printf("Connecting to %s", WIFI_SSID);
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

    int retry = 0;
    while (WiFi.status() != WL_CONNECTED && retry < 30) {
        delay(500);
        Serial.print(".");
        retry++;
    }
    Serial.println();

    if (WiFi.status() == WL_CONNECTED) {
        Serial.printf("Connected! IP: %s\n", WiFi.localIP().toString().c_str());
        return true;
    }
    return false;
}

bool fetchAndDisplay() {
    HTTPClient http;
    http.begin(IMAGE_URL);
    http.setTimeout(15000);

    int httpCode = http.GET();
    if (httpCode != HTTP_CODE_OK) {
        Serial.printf("HTTP GET failed: %d\n", httpCode);
        http.end();
        return false;
    }

    int contentLength = http.getSize();
    Serial.printf("Image size: %d bytes\n", contentLength);

    // Allocate buffer for PNG data
    WiFiClient* stream = http.getStreamPtr();
    uint8_t* buf = (uint8_t*)malloc(contentLength);
    if (!buf) {
        Serial.println("Failed to allocate memory");
        http.end();
        return false;
    }

    int bytesRead = 0;
    while (bytesRead < contentLength) {
        if (stream->available()) {
            int read = stream->readBytes(buf + bytesRead, contentLength - bytesRead);
            bytesRead += read;
        }
        delay(1);
    }
    http.end();

    Serial.printf("Downloaded %d bytes\n", bytesRead);

    // Draw PNG to canvas
    canvas.createCanvas(540, 960);
    canvas.drawPng(buf, bytesRead, 0, 0);
    canvas.pushCanvas(0, 0, UPDATE_MODE_GC16);

    free(buf);

    Serial.println("Dashboard displayed");
    return true;
}

void goToSleep() {
    Serial.printf("Sleeping for %d minutes\n", SLEEP_MINUTES);
    M5.shutdown(SLEEP_MINUTES * 60);
}
