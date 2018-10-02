FROM golang:1.11.0
ADD main main
LABEL name=AnimeSource
ENTRYPOINT ["./main"]