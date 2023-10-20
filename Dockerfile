FROM golang:1.21 as builder

WORKDIR /workspace

# Cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer.
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy go source
# COPY internal/ internal/
COPY main.go main.go
COPY internal/ internal/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/main .
USER 65532:65532

EXPOSE 8090
ENTRYPOINT ["./main"]
