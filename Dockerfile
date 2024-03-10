# Stage  1: Build the Go app
FROM golang:alpine AS build

RUN apk update && \
  apk add --no-cache build-base nodejs npm

WORKDIR  /build

COPY . .

RUN go env -w GOCACHE=/go-cache

RUN --mount=type=cache,target=/go-cache npm install -g esbuild tailwindcss

RUN esbuild --version

RUN esbuild web/src/index.ts --bundle --minify --format=esm --target=esnext --outfile=scripts/index.js

RUN go mod tidy

RUN --mount=type=cache,target=/go-cache go install github.com/a-h/templ/cmd/templ@latest

RUN templ version

RUN templ generate

RUN npx tailwindcss -i ./static/input.css -o ./static/output.css --minify

RUN CGO_ENABLED=1 go build -v -o . ./cmd/goth-starter-template/main.go

# Stage  2: Run the Go app
FROM alpine:latest

WORKDIR /app

EXPOSE 8080

ENV PORT=8080
ENV CSRF_SECRET=
ENV RUNTIME=production

# Copy the Go app binary from the builder stage
COPY --from=build /build/main .
COPY static/ static/
COPY scripts/ scripts/

CMD [ "./main" ]
