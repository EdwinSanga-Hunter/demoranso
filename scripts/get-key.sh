#!/bin/sh
# Watch the server log and print the victim key the moment the
# ransomware registers it (which happens BEFORE any file is encrypted).
#
# Usage (start this BEFORE clicking Ransomware Demo):
#
#   Terminal 1: ./scripts/run-server.sh 2>&1 | tee ~/server-key.log
#   Terminal 2: ./scripts/get-key.sh
#   Then run the ransomware — this script prints the key within a second
#   and saves a copy to ~/KEY.txt.
LOG=~/server-key.log

if [ ! -f "$LOG" ]; then
  echo "No $LOG found — start the server with logging first:"
  echo "  ./scripts/run-server.sh 2>&1 | tee ~/server-key.log"
  exit 1
fi

echo "Waiting for the ransomware to register its key..."
while ! grep -q 'Successfully saved key pair' "$LOG"; do
  sleep 1
done

grep 'Successfully saved key pair' "$LOG"
KEY=$(grep -o 'Successfully saved key pair [a-f0-9]* - [a-f0-9]*' "$LOG" | tail -1 | awk '{print $NF}')
echo "$KEY" > ~/KEY.txt
echo
echo "Key captured and saved to ~/KEY.txt:"
cat ~/KEY.txt
