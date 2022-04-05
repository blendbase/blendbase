FROM golang:1.17.8-alpine3.15

RUN mkdir /code
WORKDIR /code
COPY . /code/

RUN adduser -D -h /home/blendbase blendbase && mv /code /home/blendbase && chown -R blendbase:1000 /home/blendbase/code

WORKDIR /home/blendbase/code
USER blendbase

RUN go build -o bin/blendbase


EXPOSE 8080
CMD ["./bin/blendbase-server.sh"]
