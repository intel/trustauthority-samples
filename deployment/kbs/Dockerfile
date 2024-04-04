# Copyright (C) 2024 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause

FROM golang:1.21.6-bullseye AS builder

WORKDIR /app

RUN git clone -b main https://github.com/intel/trustauthority-kbs.git
RUN cd trustauthority-kbs && GITCOMMIT=$(git describe --always); \
        VERSION=${VERSION:-v0.0.0}; \
        BUILDDATE=$(TZ=UTC date +%Y-%m-%dT%H:%M:%S%z); \
        cd cmd && env GOOS=linux CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2" \
        go build -buildmode=pie -trimpath \
        -ldflags "-linkmode=external -s -extldflags '-Wl,-z,relro,-z,now' -X intel/kbs/v1/version.BuildDate=${BUILDDATE} -X intel/kbs/v1/version.Version=${VERSION} -X intel/kbs/v1/version.GitHash=${GITCOMMIT}" \
        -o kbs

RUN git clone -b task/dockerfile-deployment https://github.com/arvind5/trustauthority-samples.git trustauthority-samples

FROM ubuntu:22.04 AS final

RUN apt-get update && apt-get install -y --no-install-recommends python3-pip openssl

# Install ca-certificates package to get the system certificates
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Install PyKMIP and its dependencies
RUN pip3 install pykmip

# Generate SSL/TLS certificates for PyKMIP server
WORKDIR /etc/pykmip/
COPY --from=builder /app/trustauthority-samples/deployment/kbs/server.conf .
COPY --from=builder /app/trustauthority-samples/deployment/kbs/generate-pykmip-certs.sh .
RUN chmod +x generate-pykmip-certs.sh && ./generate-pykmip-certs.sh

WORKDIR /

COPY --from=builder /app/trustauthority-kbs/cmd/kbs .
COPY --from=builder /app/trustauthority-samples/deployment/kbs/entrypoint.sh .

ARG USER_UID=1000
ARG USER_GID=$USER_UID
ARG USERNAME=kbs

# Run as non-root user
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME

# Create directories for KBS
RUN mkdir -p /opt/kbs/keys-transfer-policy \
    /opt/kbs/users \
    /opt/kbs/keys \
    /etc/kbs/certs/signing-keys \
    /etc/kbs/certs/trustedca \
    /etc/kbs/certs/trustedjwt \
    /etc/kbs/certs/tls && \
    chmod +x entrypoint.sh

RUN chown -R $USER_UID:$USER_GID kbs /etc/pykmip /opt/kbs /etc/kbs && mkdir /var/log/pykmip && chown $USER_GID:$USER_GID /var/log/pykmip

USER $USERNAME
EXPOSE 9443
ENTRYPOINT ["/entrypoint.sh"]