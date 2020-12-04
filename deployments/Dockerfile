FROM golang:1.15.5-alpine3.12 AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /executionengine cmd/executionengine.go

FROM scratch AS bin
COPY --from=build /executionengine /
ENTRYPOINT ["/executionengine"]