#include <M5EPD.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include "config.h"

M5EPD_Canvas canvas(&M5.EPD);

bool connectWiFi();
bool fetchAndDisplay();
void drawBatteryBar();
void goToSleep();

void setup() {
    M5.begin();
    M5.EPD.SetRotation(0);
    M5.EPD.Clear(true);
    M5.RTC.begin();

    if (!connectWiFi()) {
        Serial.println("WiFi connection failed, going to sleep");
        goToSleep();
        return;
    }

    if (!fetchAndDisplay()) {
        Serial.println("Failed to fetch dashboard image");
    }

    WiFi.disconnect(true);
    WiFi.mode(WIFI_OFF);
    goToSleep();
}

void loop() {
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

    WiFiClient* stream = http.getStreamPtr();
    uint8_t* buf = (uint8_t*)ps_malloc(contentLength);
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

    canvas.createCanvas(960, 540);
    canvas.drawJpg(buf, bytesRead, 0, 0);
    drawBatteryBar();
    canvas.pushCanvas(0, 0, UPDATE_MODE_GC16);

    free(buf);

    Serial.println("Dashboard displayed");
    return true;
}

void drawBatteryBar() {
    uint32_t voltage = M5.getBatteryVoltage();
    int percent = (voltage - 3200) * 100 / (4200 - 3200);
    if (percent < 0) percent = 0;
    if (percent > 100) percent = 100;

    int barWidth = 960 * percent / 100;
    canvas.fillRect(0, 0, barWidth, 1, 0);  // black
    Serial.printf("Battery: %dmV (%d%%)\n", voltage, percent);
}

void goToSleep() {
    Serial.printf("Sleeping for %d minutes\n", SLEEP_MINUTES);
    M5.shutdown(SLEEP_MINUTES * 60);
}
