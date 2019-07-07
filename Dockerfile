FROM golang as builder
RUN go get -d -v k8s.io/apimachinery/pkg/api/errors
RUN go get -d -v k8s.io/client-go/kubernetes
RUN go get -d -v k8s.io/client-go/rest 
RUN mkdir /build 
COPY main.go /build/
WORKDIR /build 
RUN go build -o main .
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]