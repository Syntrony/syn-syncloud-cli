FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV GO_VERSION=1.22.5

RUN apt-get update && apt-get install -y \
    curl \
    git \
    sudo \
    ca-certificates \
    gnupg \
    lsb-release \
    build-essential \
    iproute2 \
    && rm -rf /var/lib/apt/lists/*

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

RUN go version

# Install delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

ENV PATH="/root/go/bin:${PATH}"

CMD ["sleep", "infinity"]