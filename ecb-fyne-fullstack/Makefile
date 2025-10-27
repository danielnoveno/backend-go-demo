.PHONY: help install install-windows install-pi run-backend run-frontend run-sensor run-sensor-pi build-all build-pi lint fmt

help:
	@echo "Available commands:"
	@echo "  make install           - Install Go and Python deps for local demo (Windows friendly)"
	@echo "  make install-windows   - Same as install (alias)"
	@echo "  make install-pi        - Install dependencies on Raspberry Pi"
	@echo "  make run-backend       - Run backend server"
	@echo "  make run-frontend      - Run Fyne GUI"
	@echo "  make run-sensor        - Run sensor reader in simulation mode"
	@echo "  make run-sensor-pi     - Run sensor reader with GPIO (Pi)"
	@echo "  make build-all         - Build Windows/desktop binaries into ./bin"
	@echo "  make build-pi          - Cross-compile ARM binaries for Raspberry Pi"

install: install-windows

install-windows:
	cd backend && go mod tidy
	cd frontend && go mod tidy
	cd sensor && python -m pip install -r requirements.txt

install-pi:
	cd backend && go mod tidy
	cd frontend && go mod tidy
	cd sensor && python3 -m pip install -r requirements.pi.txt

run-backend:
	cd backend && go run .

run-frontend:
	cd frontend && go run .

run-sensor:
	cd sensor && powershell -ExecutionPolicy Bypass -Command "set USE_SIMULATED_SENSOR=1; python sensor_reader.py"

run-sensor-pi:
	cd sensor && USE_SIMULATED_SENSOR=0 python3 sensor_reader.py

build-all:
	mkdir -p bin
	cd backend && go build -o ../bin/ecb-backend.exe
	cd frontend && go build -o ../bin/ecb-frontend.exe
	cd sensor && python -m compileall sensor_reader.py

build-pi:
	mkdir -p bin
	cd backend && GOOS=linux GOARCH=arm GOARM=7 go build -o ../bin/ecb-backend-pi
	cd frontend && GOOS=linux GOARCH=arm GOARM=7 go build -o ../bin/ecb-frontend-pi
