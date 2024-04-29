FROM golang:1.22-bookworm as build
ENV GOOS=linux \
  GO111MODULE="on"

RUN mkdir -m 700 /root/.ssh; \
  touch -m 600 /root/.ssh/known_hosts; \
  ssh-keyscan github.com > /root/.ssh/known_hosts; \
  git config --global url."git@github.com:".insteadOf "https://github.com/"

WORKDIR /code

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w -X 'main.version=$(VERSION)'" -o app-binary ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /code/docs/swagger.yaml /app/docs/swagger.yaml
COPY --from=build /code/app-binary /app/app-binary

ENTRYPOINT [ "/app/app-binary" ]