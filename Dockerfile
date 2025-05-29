FROM node:lts-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY ./ui/ .
RUN npm install
RUN npm run build

FROM golang:tip-alpine
RUN apk add --update ca-certificates
WORKDIR /go
COPY ./api .
RUN rm -r pkg/embeded/ui/dist
COPY --from=0 /app/dist ./pkg/embeded/ui/dist
RUN go build -o /go-app

FROM alpine
COPY --from=1 /etc/ssl /etc/ssl
COPY --from=1 /etc/ca-certificates /etc/cat-certificates
COPY --from=1 /go-app /go-app
