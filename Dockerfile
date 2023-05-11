FROM gcr.io/distroless/static-debian11:debug AS build

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /tmp

COPY openshift4_mirror /

ARG BUILD_DATE
ARG BUILD_VERSION
ARG VCS_REF
ARG VCS_URL

LABEL org.opencontainers.image.created=$BUILD_DATE
LABEL org.opencontainers.image.title="openshift4-mirror-go"
LABEL org.opencontainers.image.description="CLI tool for mirroring Openshift4 content"
LABEL org.opencontainers.image.source=$VCS_URL
LABEL org.opencontainers.image.revision=$VCS_REF
LABEL org.opencontainers.image.vendor="Shane Dell"
LABEL org.opencontainers.image.version=$BUILD_VERSION

ENTRYPOINT ["/openshift4_mirror"]
