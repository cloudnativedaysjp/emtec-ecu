development
===

## 事前準備

手元で Websocket を有効化した OBS を立ち上げ、検証用にシーン・ソースを用意します。

### Download

* [Download OBS](https://obsproject.com/ja/download)
* [Download obs-websocket](https://github.com/obsproject/obs-websocket/releases/)

### 検証用シーン・ソースの準備

TBW

## 検証環境の起動

Docker Compose を用い、cnd-operation-server, seaman (SlackBot), dk-mock-server (Dreamkast API の Mock Server) を立ち上げます。

* `.env` ファイルの作成

```
SLACK_BOT_TOKEN=...
SLACK_APP_TOKEN=...
SWITHCER01_SLACK_CHANNEL_ID=...
GITHUB_ACCESS_TOKEN=...
SWITCHER01_HOST=...
SWITCHER01_PASSWORD=...
```

* コンテナの起動

```bash
sed -e "s/DATE/`date "+%Y-%m-%d/"`" \
-e "s/TIME1/"2000-01-01T`date "+%H:%M:%S.%3N" -d "1 minute"`+09:00"/" \
-e "s/TIME2/"2000-01-01T`date "+%H:%M:%S.%3N" -d "4 minute"`+09:00"/" \
-e "s/TIME3/"2000-01-01T`date "+%H:%M:%S.%3N" -d "7 minute"`+09:00"/" \
-e "s/TIME4/"2000-01-01T`date "+%H:%M:%S.%3N" -d "9 minute"`+09:00"/" \
dk-mock-server-base.yaml > dk-mock-server.yaml; docker compose up -d
```
