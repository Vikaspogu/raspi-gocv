#FROM docker.io/vikaspogu/rpi-gocv:4.2.0
FROM golang@sha:aebb0b5f2c05fc84e9d85ac2dbd7ab5f33c8cd96cc4679983fa1d648f3ef3552
RUN mkdir /app
ADD . /app/ 
WORKDIR /app 
#RUN go build -o main .
CMD ["/app/main"]
