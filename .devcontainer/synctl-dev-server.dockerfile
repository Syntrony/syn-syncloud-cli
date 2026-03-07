FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV GO_VERSION=1.22.5

# -----------------------------
# Base tools
# -----------------------------
RUN apt-get update && apt-get install -y \
    sudo \
    curl \
    git \
    ca-certificates \
    iproute2 \
    gnupg \
    lsb-release \
    build-essential \
    docker.io \
    jq \
    vim \
    net-tools \
    iputils-ping \
    dnsutils \
    && apt-get clean

# -----------------------------
# Install Go (multi-arch)
# -----------------------------
ARG TARGETARCH

RUN if [ "$TARGETARCH" = "arm64" ]; then \
        GO_ARCH="arm64"; \
    else \
        GO_ARCH="amd64"; \
    fi && \
    curl -LO https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-${GO_ARCH}.tar.gz && \
    rm go${GO_VERSION}.linux-${GO_ARCH}.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# verify go
RUN go version

# -----------------------------
# Install Delve (debugger)
# -----------------------------
RUN go install github.com/go-delve/delve/cmd/dlv@latest

ENV PATH="/root/go/bin:${PATH}"

# -----------------------------
# Install kubectl
# -----------------------------
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
 && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl \
 && rm kubectl

# -----------------------------
# Install k3d (for dev cluster)
# -----------------------------
RUN curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash

# -----------------------------
# Workspace
# -----------------------------
WORKDIR /workspace

CMD [ "sleep", "infinity" ]