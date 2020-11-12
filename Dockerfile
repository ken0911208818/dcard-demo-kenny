FROM golang:1.13-alpine
WORKDIR /dcard-demo-kenny
ADD . /dcard-demo-kenny
RUN cd /dcard-demo-kenny && go build
RUN mv ecc-private-key.pem /opt
RUN mv ecc-public-key.pem /opt

ENTRYPOINT ["./dcard-demo-kenny"]