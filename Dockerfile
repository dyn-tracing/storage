FROM golang:1.14.3 as bd
WORKDIR storage
ADD . .
RUN GOPROXY=direct GOSUMDB=off go build -a -o /server .
CMD ["/server"]

