FROM docker.io/vikaspogu/rpi-gocv:4.2.0
COPY qemu-arm-static /usr/bin/qemu-arm-static
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main .
CMD ["/app/main"]
