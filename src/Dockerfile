FROM alpine:latest

MAINTAINER johnpoint
ENV TZ=Asia/Shanghai

WORKDIR /root/service
COPY RssReader bin/
COPY config_local.yaml config/

ENV CMDNAME=""
ENV CONFIGNAME="config_local.yaml"

CMD ["sh","-c","/root/service/bin/RssReader $CMDNAME --config=/root/service/config/$CONFIGNAME"]