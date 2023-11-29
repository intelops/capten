# Build stage
FROM golang AS builder

WORKDIR /go/src/app

COPY . /go/src/app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o capten ./cmd/main.go

# Final stage
FROM alpine

WORKDIR /app

COPY --from=builder /go/src/app /app/

RUN apk add --no-cache bash

# Clone the Git repository for Terraform modules
RUN apk add --no-cache git && \
    git clone https://github.com/kube-tarian/controlplane-dataplane.git /app/terraform_modules && \
    rm -rf /app/terraform_modules/.git

# Copy necessary files from the build context into the image
COPY --from=builder /go/src/app/config /app/config
COPY --from=builder /go/src/app/templates /app/templates
COPY --from=builder /go/src/app/apps /app/apps
COPY --from=builder /go/src/app/README.md /app/README.md

# Download and extract Terraform binary for Linux
RUN wget https://releases.hashicorp.com/terraform/0.12.31/terraform_0.12.31_linux_amd64.zip && \
    unzip terraform_0.12.31_linux_amd64.zip -d /app/ && \
    chmod +x /app/terraform && \
    rm terraform_0.12.31_linux_amd64.zip

# Download and extract Talosctl binary for Linux
RUN wget https://github.com/siderolabs/talos/releases/download/v1.4.8/talosctl-linux-amd64 -O /app/terraform_modules/talosctl && \
    chmod +x /app/terraform_modules/talosctl

# Make all scripts executable
RUN find /app/ -type f -name "*.sh" -exec chmod +x {} \;

ENTRYPOINT  ["./capten"]