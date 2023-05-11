# Base image from https://repo1.dso.mil/platform-one/distros/red-hat/ocp4/openshift4-mirror/-/blob/main/Dockerfile,
# updated to use golang not python. Also added support for different archs

FROM registry.access.redhat.com/ubi8/ubi:latest

ARG arch="amd64"
ENV arch=${arch}

ENV PYCURL_SSL_LIBRARY=openssl

ENV LC_CTYPE=en_US.UTF-8
ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US.UTF-8

LABEL \
    name="openshift4-mirror-go" \
    description="Utility for mirroring OpenShift 4 content" \
    maintainer="Shanedell"

USER root

RUN \
    yum install -y \
        wget \
        vim \
        which \
        make \
    && yum clean all \
    && wget https://go.dev/dl/go1.19.8.linux-${arch}.tar.gz \
    && tar -xzf go*.linux-${arch}.tar.gz -C /usr/local \
    && rm -rf go*.linux-${arch}.tar.gz \
    && echo 'export PATH=/usr/local/go/bin:$PATH' >> /root/.bashrc \
    && echo 'export PS1="\n\[\e[34m\]\u\[\e[m\] at \[\e[32m\]\h\[\e[m\] in \[\e[33m\]\w\[\e[m\] \[\e[31m\]\n\\$\[\e[m\] "' >> /root/.bashrc \
    && mkdir -p /app/app

COPY . /app

# Install Golang dependencies and build executable
WORKDIR /app

ENV PATH=$PATH:/usr/local/go/bin
RUN go build -o openshift_mirror .

ENTRYPOINT ["/app/entrypoint.sh"]