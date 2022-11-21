FROM golang:1.19.3-alpine AS build
WORKDIR /app
COPY . .
ENV GOPROXY=https://goproxy.cn,direct
RUN GOOS=linux go build -o main main.go

FROM alpine:latest
WORKDIR /app

COPY --from=build /app/main .
COPY application.yaml .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]