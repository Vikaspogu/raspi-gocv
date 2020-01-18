FROM docker.io/vikaspogu/gocv:4.2.0
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main . 
CMD ["/app/main.go 0 0.0.0.0:8080"]
