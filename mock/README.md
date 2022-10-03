cnd-operation-server mock
=========================

手元で cnd-operation-server を動作させるために、Dreamkast の talks, tracks API を提供するためのmock。  
mock server は https://github.com/Thiht/smocker を利用する。

## mock server 利用手順
以下、smocker を docker で起動した場合の手順を記載する。

1. smocker を docker で起動  
smocker の [document](https://smocker.dev/guide/installation.html#with-docker) に記載の方法で起動する。
```
docker run -d \
  --restart=always \
  -p 8080:8080 \
  -p 8081:8081 \
  --name smocker \
  thiht/smocker
```

`8080` が mock server port で、`8081` が設定用のport。
```
8080 is the mock server port. It will expose the routes you register through the configuration port.
8081 is the configuration port. It's the port you will use to register new mocks. This port also exposes a user interface.
```

2. `mock.yaml`を利用して、起動した mock server に mock シナリオを追加する。
```
$ curl -XPOST --header "Content-Type: application/x-yaml" --data-binary "@mock.yaml" localhost:8081/mocks

{"message":"Mocks registered successfully"}
```

## mock response
`mock.yaml` に定義されている API エンドポイントは以下の2つ。
- `/api/v1/tracks?eventAbbr=cnsec2022`
- `/api/v1/talks?eventAbbr=cnsec2022&trackId=29`

`/api/v1/talks?eventAbbr=cnsec2022&trackId=29` は、APIをたたいた時刻から1分後、2分後、3分後に開始するSession情報を返却する。

### Response Example
- `/api/v1/tracks?eventAbbr=cnsec2022`
```
$ curl "localhost:8080/api/v1/tracks?eventAbbr=cnsec2022"

[
  {
    "id": 29,
    "name": "A",
    "videoPlatform": "ivs",
    "videoId": null,
    "channelArn": null,
    "onAirTalk": {
      "id": 645,
      "talk_id": 1503,
      "site": null,
      "url": null,
      "on_air": true,
      "created_at": "2022-08-04T03:20:57.000+09:00",
      "updated_at": "2022-08-05T18:45:51.000+09:00",
      "video_id": "",
      "slido_id": null,
      "video_file_data": null
    }
  },
  {
    "id": 30,
    "name": "B",
    "videoPlatform": "ivs",
    "videoId": null,
    "channelArn": null,
    "onAirTalk": {
      "id": 646,
      "talk_id": 1504,
      "site": null,
      "url": null,
      "on_air": true,
      "created_at": "2022-08-04T03:20:58.000+09:00",
      "updated_at": "2022-08-05T18:47:50.000+09:00",
      "video_id": "",
      "slido_id": null,
      "video_file_data": null
    }
  },
  {
    "id": 31,
    "name": "C",
    "videoPlatform": "ivs",
    "videoId": null,
    "channelArn": null,
    "onAirTalk": {
      "id": 647,
      "talk_id": 1505,
      "site": null,
      "url": null,
      "on_air": true,
      "created_at": "2022-08-04T03:20:58.000+09:00",
      "updated_at": "2022-08-05T18:45:38.000+09:00",
      "video_id": "",
      "slido_id": null,
      "video_file_data": null
    }
  }
]
```

- /api/v1/talks?eventAbbr=cnsec2022&trackId=29`
```
$ curl "localhost:8080/api/v1/talks?eventAbbr=cnsec2022&trackId=29"

[
  {
    "id": 1433,
    "conferenceId": 6,
    "trackId": 29,
    "videoPlatform": null,
    "videoId": "https://d3pun3ptcv21q4.cloudfront.net/mediapackage/cnsec2022/talks/1433a/12/playlist.m3u8",
    "title": "Simplify Cloud Native Security with Trivy",
    "abstract": "クラウド環境への移行に伴い必要なセキュリティ対策も大きく変化しました。しかしこれらの対策には多くのツールを組み合わせて使う必要があり、導入・学習コストが高くなっています。そこで本発表では、OSS であるTrivyを用いて特に攻撃へと繋がりやすい依存ライブラリの脆弱性や脆弱なインフラ設定、誤ってコミットされたパスワード等の検知を一括で行う方法について説明します。また、実際にCloudFormationやHelmチャートをスキャンすることでデプロイ前に危険な設定を検知するデモを行います。",
    "speakers": [
      {
        "id": 1283,
        "name": "Teppei Fukuda"
      }
    ],
    "dayId": 16,
    "showOnTimetable": true,
    "startTime": "2022-10-03T18:49:47.939+09:00",
    "endTime": "2022-10-03T18:50:47.939+09:00",
    "talkDuration": 0,
    "talkDifficulty": "初級者",
    "talkCategory": "",
    "onAir": false,
    "documentUrl": "https://speakerdeck.com/knqyf263/simplify-cloud-native-security-with-trivy/",
    "conferenceDayId": 16,
    "conferenceDayDate": "2022-08-05",
    "startOffset": 0,
    "endOffset": 0,
    "actualStartTime": "2022-10-03T18:49:47.939+09:00",
    "actualEndTime": "2022-10-03T18:50:47.939+09:00",
    "presentationMethod": "オンライン登壇"
  },
  {
    "id": 1438,
    "conferenceId": 6,
    "trackId": 29,
    "videoPlatform": null,
    "videoId": "https://dreamkast-ivs-stream-archive-prd.s3.amazonaws.com/medialive/cnsec2022/talks/1438/1438.m3u8",
    "title": "Sysdig Secure/Falcoの活用術！〜Kubernetes基盤の脅威モデリングとランタイムセキュリティの強化〜",
    "abstract": "本講演では、IIJのマルチテナントなKubernetes基盤であるIKE（IIJ Kubernetes Engine）上の実運用システムやフローに対して脅威モデリングを行い、Sysdig Secure/Falcoを用いてIKEのランタイムセキュリティ を強化するために実施したPoCについて紹介します。どのようにしてKubernetesを基盤とした実運用システムやフローの潜在的な脅威を見つけるのか。また、見つかった脅威についてどのような対策をしていくのか。Sysdig Secure/Falcoを活用しながら試行錯誤したプロセスを解説していきます。",
    "speakers": [
      {
        "id": 1290,
        "name": "Chihiro Hasegawa / Han Li"
      }
    ],
    "dayId": 16,
    "showOnTimetable": true,
    "startTime": "2022-10-03T18:50:47.939+09:00",
    "endTime": "2022-10-03T18:51:47.939+09:00",
    "talkDuration": 0,
    "talkDifficulty": "中級者",
    "talkCategory": "",
    "onAir": false,
    "documentUrl": "https://speakerdeck.com/owlinux1000/falcofalsehuo-yong-shu-kubernetesji-pan-falsexie-wei-moderingutorantaimusekiyuriteifalseqiang-hua",
    "conferenceDayId": 16,
    "conferenceDayDate": "2022-08-05",
    "startOffset": 0,
    "endOffset": 0,
    "actualStartTime": "2022-10-03T18:50:47.939+09:00",
    "actualEndTime": "2022-10-03T18:51:47.939+09:00",
    "presentationMethod": "事前収録"
  },
  {
    "id": 1439,
    "conferenceId": 6,
    "trackId": 29,
    "videoPlatform": null,
    "videoId": "https://d3pun3ptcv21q4.cloudfront.net/mediapackage/cnsec2022/talks/1439/19/playlist.m3u8",
    "title": "Kubernetes Service Account As Multi-Cloud Identity",
    "abstract": "多くのクラウドプロバイダの提供するManaged Kubernetesサービスでは、Kubernetes Service Accountから、クラウドAPIを自動的に認証する機能を提供しています。EKSであればIAM Roles for Service Account(IRSA), GKEであればWorkload Ideneity等です。これらの機能を使えば、アプリケーションはクラウドSDKを呼び出すだけで、適切な権限をもつ、自動的に払い出された一時的なクレデンシャルを使って、シームレスかつ安全にクラウドAPIにアクセスすることが可能です。\r\n\r\nしかし、これがオンプレミスで稼働しているようなUnmanagedなKubernetesであればどうでしょうか。永続的なAPI キーをクラウドプロバイダで発行し、それらを各アプリケーションPodに マウントすることでクラウドAPIの認証を実現していることが多いのではないでしょうか。永続的なAPIキーはそれ自体がセキュリティリスクとなりますし、リスク低減のためにはローテーションの運用負荷も存在します。アプリケー ションが利用するAPIキーが複数存在したり、複数のクラウドプロバイダとやり取りする環境だと、これらのリスク・運用負荷はより大きくなります。\r\n\r\n本セッションでは、KubernetesのServiceAccountIssuerDiscoveryの機能 と、各種クラウドプロバイダが提供しているIdentity Federationの機能を活用して、永続的なAPI Keyを利用せず、KubernetesのService Accountを複数クラウド共通のIdentityとして利用可能にする方法について共有します。これによって、クラウドプロバイダ側でKubernetes Service Accountへの権限を付与するだけで、シームレスかつ安全に、複数クラウドAPIにアクセスできる環境を実現できます。",
    "speakers": [
      {
        "id": 1291,
        "name": "Shingo Omura"
      }
    ],
    "dayId": 16,
    "showOnTimetable": true,
    "startTime": "2022-10-03T18:51:47.939+09:00",
    "endTime": "2022-10-03T18:52:47.939+09:00",
    "talkDuration": 0,
    "talkDifficulty": "中級者",
    "talkCategory": "",
    "onAir": false,
    "documentUrl": "https://www.slideshare.net/pfi/kubernetes-service-account-as-multicloud-identity-cloud-native-security-conference-2022-cnsec2022",
    "conferenceDayId": 16,
    "conferenceDayDate": "2022-08-05",
    "startOffset": 0,
    "endOffset": 0,
    "actualStartTime": "2022-10-03T18:51:47.939+09:00",
    "actualEndTime": "2022-10-03T18:52:47.939+09:00",
    "presentationMethod": "事前収録"
  }
]
```
