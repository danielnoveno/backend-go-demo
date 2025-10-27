#!/usr/bin/env python3
import os
import random
import time
from datetime import datetime

import requests

USE_SIMULATION = os.getenv("USE_SIMULATED_SENSOR", "").lower() in {"1", "true", "yes"}

GPIO = None
SENSOR_PIN = 17

if not USE_SIMULATION:
    try:
        import RPi.GPIO as GPIO  # type: ignore

        GPIO.setmode(GPIO.BCM)
        GPIO.setup(SENSOR_PIN, GPIO.IN)
    except (ImportError, RuntimeError) as err:
        print(f"⚠️  GPIO module not available ({err}), falling back to simulation mode.")
        USE_SIMULATION = True

API_URL = os.getenv("API_BASE_URL", "http://127.0.0.1:8080/api") + "/sensor"
POST_INTERVAL = float(os.getenv("POST_INTERVAL_SECONDS", "3"))


def read_sensor() -> int:
    """Read sensor state. On Windows this will always simulate."""
    if USE_SIMULATION or GPIO is None:
        return random.choice([0, 1])
    return GPIO.input(SENSOR_PIN)


def generate_value() -> float:
    """Generate a sample value for the sensor measurement."""
    return random.uniform(20.0, 30.0)


def send_to_api(status: str, value: float) -> None:
    try:
        data = {"status": status, "value": value}
        response = requests.post(API_URL, json=data, timeout=5)
        if response.status_code == 200:
            print(f"✅ [{datetime.now().strftime('%H:%M:%S')}] Sent: {status} | {value:.2f}")
        else:
            print(f"⚠️  Server responded with status code {response.status_code}")
    except requests.exceptions.RequestException as exc:
        print(f"❌ Error sending data: {exc}")


def loop_once() -> None:
    sensor_input = read_sensor()
    status = "DETECTED" if sensor_input == 1 else "IDLE"
    value = generate_value()
    send_to_api(status, value)


def main() -> None:
    print("🚀 ECB Sensor Reader Started")
    print("=" * 50)
    print(f"Simulation mode: {'ON' if USE_SIMULATION else 'OFF (real GPIO)'}")
    print(f"Posting interval: {POST_INTERVAL:.1f}s")
    print(f"Target API: {API_URL}")
    print("=" * 50)

    try:
        while True:
            loop_once()
            time.sleep(POST_INTERVAL)
    except KeyboardInterrupt:
        print("\n👋 Stopping sensor reader...")
    finally:
        if GPIO is not None and not USE_SIMULATION:
            GPIO.cleanup()


if __name__ == "__main__":
    main()
