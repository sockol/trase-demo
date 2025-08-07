FROM golang:alpine
WORKDIR /api
RUN apk add --no-cache make
COPY . .
RUN make build
EXPOSE 8080
CMD ["make", "run/live"]