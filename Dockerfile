#FROM docker.io/vikaspogu/rpi-gocv:4.2.0
FROM arm/v7/golang:1.13.7
RUN mkdir /app
ADD . /app/ 
WORKDIR /app 
#RUN go build -o main .
CMD ["/app/main"]
