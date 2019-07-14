FROM golang as builder
LABEL maintainer="Pritish Payaningal <findpritish@gmail.com>"
WORKDIR $GOPATH/src/github.com/findpritish/plogger
RUN go get -d -v k8s.io/apimachinery/pkg/api/errors
RUN go get -d -v k8s.io/apimachinery/pkg/labels
RUN go get -d -v github.com/spf13/cobra
RUN go get -d -v k8s.io/api/core/v1
RUN go get -d -v github.com/fatih/color
RUN go get -d -v k8s.io/client-go/kubernetes/typed/core/v1
RUN go get -d -v k8s.io/apimachinery/pkg/apis/meta/v1
RUN go get -d -v k8s.io/apimachinery/pkg/watch
RUN go get -d -v github.com/pkg/errors
RUN go get -d -v k8s.io/client-go/kubernetes
RUN go get -d -v k8s.io/client-go/rest
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/plogger
FROM alpine
#RUN adduser -S -D -H -h /app appuser
#USER appuser
COPY --from=builder /bin/plogger .
ENV NSINPUT "learning-ns"
ENV EXINPUT "plogger"
CMD ./plogger -n $NSINPUT --exclude-container $EXINPUT .