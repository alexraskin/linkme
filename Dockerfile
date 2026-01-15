FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG BUILD_TIME

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    BUILD_TIME="${BUILD_TIME:-$(date "+%b %d, %Y")}" && \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -ldflags "-X 'main.buildTime=${BUILD_TIME}'" -o linkme github.com/alexraskin/linkme

FROM alpine

RUN apk --no-cache add ca-certificates

COPY --from=build /build/linkme /bin/linkme

EXPOSE 8080

CMD ["/bin/linkme"]