FROM docker.io/vikaspogu/rpi-gocv:4.2.0
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN GOARCH=arm go build -o main .
CMD ["/app/main"]
