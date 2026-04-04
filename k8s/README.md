# Kubernetes (microk8s) デプロイ

自宅のmicrok8sサーバーにデプロイする構成。S3不要で、CronJobが画像を生成し、nginxで配信します。

## アーキテクチャ

```
CronJob (10分間隔) → PVC に dashboard.jpg を保存
                            ↓
M5Paper → nginx (NodePort:30080) → PVC から配信
```

## セットアップ

### 1. Docker イメージのビルド

```shell
docker build -t ghcr.io/kangaechu/m5paper-dashboard:latest .
# microk8sの場合
docker save ghcr.io/kangaechu/m5paper-dashboard:latest | microk8s ctr image import -
```

### 2. Secret の作成

```shell
cp k8s/secret.yaml.sample k8s/secret.yaml
# secret.yaml の GOOGLE_CREDENTIALS_JSON を設定
```

### 3. デプロイ

```shell
microk8s kubectl apply -k k8s/
```

### 4. M5Paper の設定

M5Paper の `IMAGE_URL` を以下に設定：

```
http://<microk8s-server-ip>:30080/dashboard.jpg
```

## 確認

```shell
# CronJob の状態
microk8s kubectl -n m5paper-dashboard get cronjob

# 手動実行
microk8s kubectl -n m5paper-dashboard create job --from=cronjob/m5paper-dashboard manual-test

# ログ確認
microk8s kubectl -n m5paper-dashboard logs job/manual-test

# 画像確認
curl http://localhost:30080/dashboard.jpg -o test.jpg
```
