FROM golang:1.6-onbuild

WORKDIR /go/src/app
ENV ALLOW_CONTAINER_ROOT=1

COPY . /go/src/app
RUN \
	go-wrapper download && \
	go-wrapper install && \
	mkdir -p /export/docker

EXPOSE 9000
ENTRYPOINT ["go-wrapper", "run", "server"]
VOLUME ["/export"]
CMD ["/export"]
