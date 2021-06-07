# m5paper-dashboard-vue

M5Paperをダッシュボード用の表示デバイスとして使用するNode.jsのアプリです。

# Demo

![DEMO](https://user-images.githubusercontent.com/989985/120999952-a7f06880-c7c4-11eb-82a8-f775fb6a5798.jpg)

[M5Paper](https://github.com/m5stack/M5EPD)はESP32を内蔵し、解像度540 x 960 4.7インチの電子ペーパーを搭載したデバイスです。
バッテリーを内蔵しているので、どこにでも取り付けができ、長時間稼働することができます。

M5Paperに天気や温度・湿度・ニュース・今日の予定などを表示し、それをリビングやキッチンの壁に貼り付けることによって情報端末として使用することを目的としています。

# 構成

![m5paper-dashboard-vue](https://user-images.githubusercontent.com/989985/121002160-028ac400-c7c7-11eb-8459-ced520afae4a.png)

ServerはWebサーバとして稼働します。自分自身やバックエンドのAPIを使用し、ダッシュボードに表示する情報を保持します。フロントエンドはVue.js、バックエンドはExpressで動いています。
ScraperはHeadless browserです。ServerにHTTPでアクセスし、取得したHTMLをレンダリングし、540 x 960のサイズで保存したスクリーンショットをServerの `/public/dashboard.png` に保存します。これは定期的に実行します。
M5Paperは定期的にServerにアクセスし、スクリーンショットを取得して表示します。
それ以外の時間はSleepします。

# Requirement

- Raspberry Pi
- Ansible

# Installation

Ansibleを使用します。

```shell
cd ansible
ansible-playbook -i inventory raspberrypi.yml
```