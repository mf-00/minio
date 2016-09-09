FROM golang:1.6

WORKDIR /go/src/app
ENV ALLOW_CONTAINER_ROOT=1

COPY . /go/src/app
RUN \
	go-wrapper download && \
	go-wrapper install && \
	mkdir -p /export/docker

EXPOSE 9000
ENTRYPOINT ["go-wrapper", "run", "server" "-e \"MINIO_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE\"" "-e \"MINIO_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY\""]
VOLUME ["/export"]
CMD ["/export"]
