# Copyright (C) 2024 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause

FROM golang:1.21.6-bullseye AS builder

WORKDIR /app

RUN git clone -b main https://github.com/intel/trustauthority-samples.git
RUN cd trustauthority-samples && GITCOMMIT=$(git describe --always); \
        VERSION=${VERSION:-v0.0.0}; \
        BUILDDATE=$(TZ=UTC date +%Y-%m-%dT%H:%M:%S%z); \
        cd tdx-ml-workload && env CGO_CFLAGS_ALLOW="-f.*" CGO_CFLAGS="-O2 -D_FORTIFY_SOURCE=2" \
        go build -buildmode=pie -trimpath \
        -ldflags "-linkmode=external -s -extldflags '-Wl,-z,relro,-z,now' -X github.com/intel/trustauthority-samples/tdxexample/version.BuildDate=${BUILDDATE} -X github.com/intel/trustauthority-samples/tdxexample/version.Version=${VERSION} -X github.com/intel/trustauthority-samples/tdxexample/version.GitHash=${GITCOMMIT}" \
        -o trustauthority-demo

RUN cd /app/trustauthority-samples/encryptor && go build encrypt.go

FROM ubuntu:22.04 AS final

RUN apt update && apt install -y --no-install-recommends openssl jq curl gnupg wget

# Install ca-certificates package to get the system certificates
RUN apt update && apt install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ARG PLATFORM
WORKDIR /app

# Create keypair for workload model encryption
RUN openssl genrsa -out keypair_tmp.pem 2048 && \
    openssl rsa -in keypair_tmp.pem -outform PEM -out keypair.pem && \
    openssl rsa -in keypair.pem -pubout -out public.crt

COPY --from=builder /app/trustauthority-samples/tdx-ml-workload/trustauthority-demo .
COPY --from=builder /app/trustauthority-samples/encryptor/encrypt /usr/local/bin/
COPY --from=builder /app/trustauthority-samples/deployment/sample-workload/encrypt-model.sh .
COPY --from=builder /app/trustauthority-samples/deployment/sample-workload/entrypoint.sh .
COPY --from=builder /app/trustauthority-samples/encryptor/diabetes-linreg.model .

# Install tdx-cli
RUN apt update && curl -L https://raw.githubusercontent.com/intel/trustauthority-client-for-go/main/release/install-tdx-cli-${PLATFORM}.sh | bash -

ARG USER_UID=1000
ARG USER_GID=$USER_UID
ARG USERNAME=workload

# Run as non-root user
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME

# Create directories for Workload
RUN mkdir -p /opt/trustauthority-demo/ && \
    chmod +x /app/encrypt-model.sh /app/entrypoint.sh

RUN touch /etc/model.enc

RUN chown -R $USER_UID:$USER_GID /usr/local/bin/encrypt /app/trustauthority-demo /opt/trustauthority-demo/ /app/ /etc/model.enc /usr/bin/trustauthority-cli

USER $USERNAME
ENTRYPOINT ["/app/entrypoint.sh"]
