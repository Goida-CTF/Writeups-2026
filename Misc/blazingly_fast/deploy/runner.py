import subprocess
import os
import secrets


print("Input your code to insert into main.rs.")
print("End with 3 blank lines!")

lines = []
while len(lines) < 3 or any(x.strip() != "" for x in lines[-3:]):
    lines.append(input())
code = "\n".join(lines)

nonce = secrets.token_hex(8)
os.mkdir(f"/chroot/opt/{nonce}")
with open(f"/chroot/opt/{nonce}/flag.txt", "w") as f:
    f.write(os.getenv("FLAG", "flag{test_flag}"))

with open(f"./src/src/main.rs", "w") as f:
    f.write(code)

p = subprocess.Popen(["nsjail", "--config", "nsjail.cfg", "--", "/home/user/.cargo/bin/cargo", "build"])
p.wait()
