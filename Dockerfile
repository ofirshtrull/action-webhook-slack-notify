FROM golang:1.11-alpine3.9@sha256:7a0bf914dd581a35afb054bc02c6b7a3fa31ed6398adf95ac88fb1efffe89cf6 AS builder

LABEL "com.github.actions.icon"="code"
LABEL "com.github.actions.color"="red"
LABEL "com.github.actions.name"="Webhook Slack Notify"
LABEL "com.github.actions.description"="This action will send notification to Slack"


WORKDIR ${GOPATH}/src/github.com/partnerhero/action-webhook-slack-notify
COPY main.go ${GOPATH}/src/github.com/partnerhero/action-webhook-slack-notify

ENV CGO_ENABLED 0
ENV GOOS linux

RUN go get -v ./...
RUN go build -a -installsuffix cgo -ldflags '-w  -extldflags "-static"' -o /go/bin/webhook-slack-notify .

FROM alpine:3.18.2

COPY --from=builder /go/bin/webhook-slack-notify /usr/bin/webhook-slack-notify

ENV VAULT_VERSION 1.0.2
ENV PYTHONUNBUFFERED=1

RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python

RUN apk update \
	&& apk upgrade \
	&& apk add \
	bash \
	jq \
	ca-certificates \
	rsync


RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools shyaml # py2-pip
#  && \
# 	pip install shyaml && \
# 	rm -rf /var/cache/apk/*
# 	python
# \
# py2-pip \



# Setup Vault
RUN wget https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip && \
	unzip vault_${VAULT_VERSION}_linux_amd64.zip && \
	rm vault_${VAULT_VERSION}_linux_amd64.zip && \
	mv vault /usr/local/bin/vault

# fix the missing dependency - https://stackoverflow.com/a/35613430
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY *.sh /

RUN chmod +x /*.sh

ENTRYPOINT ["/entrypoint.sh"]
