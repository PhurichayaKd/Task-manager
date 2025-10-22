# ---- build stage ----
FROM golang:1.23 AS build
WORKDIR /app

# โหลดโมดูลก่อน เพื่อ cache
COPY go.mod go.sum ./
RUN go mod download

# คัดลอกซอร์สและคอมไพล์
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# ---- runtime stage ----
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /app/server /app/server
COPY --from=build /app/frontend /app/frontend

EXPOSE 8080
USER nonroot:nonroot
CMD ["/app/server"]
