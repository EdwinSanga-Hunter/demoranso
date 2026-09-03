# Prebuilt Windows demo binaries

These were built from this repository (`make dist`) with the server address
baked in as `http://localhost:8080` — run everything inside one Windows VM.

- `ransomware.exe` — the malware. Double-click, accept the UAC prompt, and a
  console window shows the whole attack. Without `RANSOMWARE_DIR` it walks all
  drives; set `RANSOMWARE_DIR=C:\demo` to scope it to one folder.
- `unlocker.exe` — the unlocker. Paste the 32-char key, press Enter, answer `Y`.
- `server.exe` — the demo server. Run it **first** (double-click). It keeps
  `database.db` next to it and prints `Successfully saved key pair <ID> - <KEY>`
  when the malware registers.
- `private.pem` — the RSA private key matching the public key embedded in
  `ransomware.exe`. It must stay in the same folder as `server.exe`.

Full step-by-step guides: see the main [README](../README.md).
