ARG FDB_VERSION=6.2.11
FROM foundationdb/foundationdb:${FDB_VERSION} as fdb
FROM golang:1.13.6-stretch
ARG FDB_VERSION

RUN apt update
# dnsutils is needed to have dig installed to create cluster file
RUN apt install -y dnsutils

WORKDIR /tmp

ADD /foundationdb-clients_6.2.11-1_amd64.deb /tmp/
RUN dpkg -i /tmp/foundationdb-clients_${FDB_VERSION}-1_amd64.deb

COPY --from=fdb /var/fdb/scripts/create_cluster_file.bash /

WORKDIR /go/src/app
COPY . .


RUN go get -d -v ./...
RUN go install -v ./...

CMD ["/go/bin/fdb"]