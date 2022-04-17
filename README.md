# RssReader

<img src="https://raw.githubusercontent.com/johnpoint/RssReader/v2/doc/img/logo.png" width="256px" height="256px"/>

vue + go 的在线 rss 阅读器

## Continuous build 持续构建

![](https://github.com/johnpoint/RssReader/workflows/RssReader_Auto_Build/badge.svg) ![](https://github.com/johnpoint/RssReader/workflows/RssReader_Web_Build/badge.svg)

## Getting Started 快速部署

### backend 后端
```
mkdir rssreader/config -p && cd rssreader
docker run --name cp_config johnpoint/rssreader:latest
docker cp cp_config:/root/service/config/config_local.yaml config/
docker container rm cp_config
docker run -d \
    --name rssreader-api \
    -v $(pwd)/config:/root/service/config \
    -e CMDNAME="api" \
    -e CONFIGNAME="config_local.yaml" \
    -p 8888:80 \
    johnpoint/rssreader:latest
docker run -d \
    --name rssreader-spider \
    -v $(pwd)/config:/root/service/config \
    -e CMDNAME="spider" \
    -e CONFIGNAME="config_local.yaml" \
    johnpoint/rssreader:latest
```

### frontend 前端