FROM golang:1.21.6-bullseye AS builder
COPY .gitconfig /root/.gitconfig
COPY .netrc /root/.netrc

ARG VERSION
WORKDIR /app

RUN git clone -b main https://github.com/intel/trustauthority-kbs.git
RUN cd trustauthority-kbs && GITTAG=$(git describe --tags --abbrev=0 2>/dev/null); \
        GITCOMMIT=$(git describe --always); \
                VERSION=${VERSION:-v0.0.0}; \
                BUILDDATE=$(TZ=UTC date +%Y-%m-%dT%H:%M:%S%z); \
        cd cmd && env GOOS=linux CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2" \
         go build -buildmode=pie \
                -ldflags "-linkmode=external -s -extldflags '-Wl,-z,relro,-z,now' -X intel/kbs/v1/version.BuildDate=${BUILDDATE} -X intel/kbs/v1/version.Version=${VERSION} -X intel/kbs/v1/version.GitHash=${GITCOMMIT}" \
                -o kbs

RUN git clone -b containerization-sample  https://github.com/pravinrajr9/applications.security.amber.trustauthority-samples.git trustauthority-samples

FROM ubuntu:22.04 AS final

RUN apt-get update && apt-get install -y --no-install-recommends python3-pip openssl

# Install PyKMIP and its dependencies
RUN pip3 install pykmip

# Generate SSL/TLS certificates for PyKMIP server
WORKDIR /etc/pykmip/
COPY --from=builder /app/trustauthority-samples/deployment/kbs/server.conf .
COPY --from=builder /app/trustauthority-samples/deployment/kbs/generate-pykmip-certs.sh .
RUN chmod +x generate-pykmip-certs.sh && ./generate-pykmip-certs.sh

# Run as non-root user
RUN groupadd --gid 1000 kbs \
    && useradd --uid 1000 --gid 1000 -m kbs

# Install ca-certificates package to get the system certificates
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /

# Create directories for KBS
RUN mkdir -p /opt/kbs/keys-transfer-policy \
    /opt/kbs/users \
    /opt/kbs/keys \
    /etc/kbs/certs/signing-keys \
    /etc/kbs/certs/trustedca \
    /etc/kbs/certs/trustedjwt \
    /etc/kbs/certs/tls

COPY --from=builder /app/trustauthority-kbs/cmd/kbs .
COPY --from=builder /app/trustauthority-samples/deployment/kbs/entrypoint.sh .

RUN chown 1000:1000 kbs && chmod +x entrypoint.sh &&\
    chown 1000:1000 -R /etc/pykmip /opt/kbs /etc/kbs && mkdir /var/log/pykmip && chown 1000:1000 /var/log/pykmip

USER kbs

EXPOSE 9443

ENTRYPOINT ["/entrypoint.sh"]