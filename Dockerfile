FROM golang:1.11.0
ADD main main
EXPOSE 8000
LABEL name=GenaretorServer
ENTRYPOINT ["./main"]