FROM alpine:3.2

RUN apk --update install ca-certificates

ADD ./dockerlogstream /dockerlogstream

ENTRYPOINT ["/dockerlogstream"]
CMD ["--testing"]
