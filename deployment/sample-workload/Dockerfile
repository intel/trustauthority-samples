FROM golang:1.21.6-bullseye AS builder
COPY .gitconfig /root/.gitconfig
COPY .netrc /root/.netrc

ARG VERSION

RUN apt-get update && apt-get install -y libssl-dev

# Build the trustauthority-demo
WORKDIR /app

RUN git clone -b containerization-sample  https://github.com/pravinrajr9/applications.security.amber.trustauthority-samples.git trustauthority-samples

RUN cd trustauthority-samples && GITCOMMIT=$(git describe --always); \
        VERSION=${VERSION:-v0.0.0}; \
        BUILDDATE=$(TZ=UTC date +%Y-%m-%dT%H:%M:%S%z); \
        cd tdx-ml-workload && env  CGO_CFLAGS_ALLOW="-f.*"  CGO_CFLAGS="-O2 -D_FORTIFY_SOURCE=2" \
	go build -buildmode=pie -trimpath -ldflags "-s -linkmode=external -extldflags '-Wl,-O1,-z,relro,-z,now' -X github.com/intel/trustauthority-samples/tdxexample/version.BuildDate=${BUILDDATE} -X github.com/intel/trustauthority-samples/tdxexample/version.Version=${VERSION} -X github.com/intel/trustauthority-samples/tdxexample/version.GitHash=${GITCOMMIT}" \
    -o trustauthority-demo

RUN cd /app/trustauthority-samples/encryptor && go build encrypt.go

FROM ubuntu:20.04 AS final

RUN apt update && apt install openssl jq curl gnupg wget -y

WORKDIR /app

# Create keypair for workload model encryption
RUN openssl genrsa -out keypair_tmp.pem 2048 &&\
    openssl rsa -in keypair_tmp.pem -outform PEM -out keypair.pem &&\
    openssl rsa -in keypair.pem -pubout -out public.crt

COPY --from=builder /app/trustauthority-samples/tdx-ml-workload/trustauthority-demo .
COPY --from=builder /app/trustauthority-samples/encryptor/encrypt /usr/local/bin/
COPY --from=builder /app/trustauthority-samples/deployment/sample-workload/encrypt-model.sh .
COPY --from=builder /app/trustauthority-samples/deployment/sample-workload/entrypoint.sh .

ENV CLI_VERSION=v1.2.0
# Download tdx-cli
RUN curl -L https://raw.githubusercontent.com/intel/trustauthority-client-for-go/main/release/install-tdx-cli-dcap.sh |  bash -

RUN mkdir -p /opt/trustauthority-demo/ && \
    chmod +x /app/encrypt-model.sh /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
