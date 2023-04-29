FROM golang:1.20 as app-builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' -tags timetzdata


FROM scratch

COPY --from=app-builder /go/bin/git-estimate /git-estimate

ENTRYPOINT ["/git-estimate"]
