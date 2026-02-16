FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o preview-generator .

FROM quay.io/argoproj/argocd:v2.14.3
USER root
COPY --from=builder /build/preview-generator /usr/local/bin/preview-generator
USER 999
