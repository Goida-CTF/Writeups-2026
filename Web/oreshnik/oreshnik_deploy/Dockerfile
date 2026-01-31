FROM golang:1.25-alpine

RUN apk add --no-cache openssl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY templates/ ./templates/
COPY static/ ./static/

RUN go build -o /oreshnik ./cmd/oreshnik

RUN mkdir public
RUN openssl genrsa -out public/private.pem 2048
RUN openssl rsa -in public/private.pem -pubout -out public/public.pem

EXPOSE 8080

CMD [ "/oreshnik" ]
