# ไฟล์นี้รวมคำสั่งพัฒนาและตรวจสอบ; .PHONY ทำให้ target ทำงานแม้มีไฟล์ชื่อเดียวกัน
.PHONY: run test integration check build migrate seed

# เปิด HTTP API ด้วย environment ที่ export ไว้
run:
	go run . serve

# รันทุกชุดทดสอบพร้อม race detector และไม่ใช้ผลทดสอบจาก cache
test:
	go test -race -count=1 ./...

# บังคับระบุฐานข้อมูลทดสอบก่อนรัน store และ HTTP integration tests แบบแสดงรายละเอียด
integration:
	@test -n "$$TEST_DATABASE_URL" || (echo 'Set TEST_DATABASE_URL to a dedicated PostgreSQL test database'; exit 1)
	go test -race -count=1 -v ./store ./integration

# ตรวจรูปแบบโค้ด วิเคราะห์ด้วย vet และรันทดสอบก่อนส่งงาน
check:
	test -z "$$(gofmt -l $$(find . -name '*.go'))"
	go vet ./...
	go test -race -count=1 ./...

# สร้าง binary สำหรับรันจริง โดยไม่ฝัง path ของเครื่องที่ build
build:
	go build -trimpath -o bin/thaifood .

# สร้าง schema จาก SQL ที่ฝังในโปรแกรม โดยใช้ DATABASE_URL ปัจจุบัน
migrate:
	go run . migrate

# เติมจังหวัด ประเภทอาหาร และร้านเริ่มต้นด้วย SQL ที่ออกแบบให้รันซ้ำได้
seed:
	go run . seed
