FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY pkg ./pkg

RUN go mod download

COPY . .
RUN go build -o /gitlab-review-bot github.com/spatecon/gitlab-review-bot/cmd/gitlab-review-bot

CMD ["/gitlab-review-bot"]