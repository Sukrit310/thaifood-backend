# ไฟล์นี้ build Go binary ใน stage แรก แล้วคัดลอกเฉพาะโปรแกรมไปยัง runtime image
FROM golang:1.27.4-alpine AS build
WORKDIR /src
# คัดลอก dependency manifests ก่อน เพื่อใช้ Docker cache ได้เมื่อแก้เฉพาะ source code
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# ปิด CGO เพื่อสร้าง binary ที่ไม่พึ่ง C library พร้อมตัด path และ debug symbols ออกจากไฟล์
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/thaifood .

# runtime มี CA certificates สำหรับ TLS และใช้บัญชี app ที่ไม่มีสิทธิ์ root
FROM alpine:3.22
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 app
WORKDIR /app
COPY --from=build /out/thaifood /app/thaifood
USER app
EXPOSE 8080
# เปิด server เป็นค่าเริ่มต้น และเปลี่ยน argument เป็น migrate/seed/admin ได้เมื่อเรียก container
ENTRYPOINT ["/app/thaifood"]
CMD ["serve"]
