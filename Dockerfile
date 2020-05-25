FROM registry.access.redhat.com/ubi8/go-toolset AS builder
COPY . .
RUN go get -d -v
RUN go build -o vct

FROM registry.access.redhat.com/ubi8/ubi-minimal
RUN mkdir -p /opt/vct
COPY --from=builder /opt/app-root/src/vct /opt/vct
COPY static /opt/vct/static
WORKDIR /opt/vct
ENTRYPOINT ["./vct"]
