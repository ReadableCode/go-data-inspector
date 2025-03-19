import time

for i in range(1, 6):
    print(f"Script 1 - Processing batch {i}", flush=True)
    time.sleep(0.5)  # Simulates work

print("Script 1 - Done!", flush=True)
