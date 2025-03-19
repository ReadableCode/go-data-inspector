import time

for i in range(1, 6):
    print(f"Script 2 - Fetching data chunk {i}", flush=True)
    time.sleep(0.7)  # Simulates work

print("Script 2 - Done!", flush=True)
