FROM alpine:latest

MAINTAINER johnpoint
ENV TZ=Asia/Shanghai

WORKDIR /root/service
COPY RssReader /root/service/bin/
COPY config_local.yaml /root/service/

RUN /root/service/bin/RssReader

ENV CMDNAME=""

CMD ["sh","-c","/root/service/bin/RssReader $CMDNAME --config=/root/service/config_local.yaml"]