FROM golang:1.25 AS builder

ARG VERSION=dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build \
    -ldflags="-X 'github.com/relychan/relixy/types.BuildVersion=${VERSION}'" \
    -v -o relixy ./server/rest

# stage 2: production image
FROM gcr.io/distroless/static-debian13:nonroot

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/relixy /relixy

USER 65532

# Run the web service on container startup.
ENTRYPOINT ["/relixy"]