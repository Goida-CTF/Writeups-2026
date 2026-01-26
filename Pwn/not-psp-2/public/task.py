from random import randbytes
from hashlib import sha256
from bottle import Bottle, request, run, HTTPError, response, static_file, redirect
from pathlib import Path
import subprocess
import os

# pow related

POW_DIFFICULTY = 0 # might be adjusted on the remote
POW_ACTIVE_SALTS = []

def pow_get_solver():
    POW_ACTIVE_SALTS.append(randbytes(8).hex())
    salt = POW_ACTIVE_SALTS[-1]
    return f"from hashlib import sha256; print(next(w for w in ('{salt}' + str(i) for i in range(10 ** 9)) if sha256(w.encode()).hexdigest().startswith('{'0' * POW_DIFFICULTY}')))"

def pow_submit_work(work: str) -> bool:
    if POW_DIFFICULTY == 0:
        return True
    
    salt = next(salt for salt in POW_ACTIVE_SALTS if work.startswith(salt))
    assert sha256(work.encode()).hexdigest().startswith("0" * POW_DIFFICULTY)
    POW_ACTIVE_SALTS.remove(salt)

    return True

# the logic

MAX_UPLOAD = 200 * 1024 * 1024  # 200 MB

app = Bottle()

@app.get("/")
def index_handler():
    if POW_DIFFICULTY > 0:
        msg = f"""
Для начала, нужно запустить этот PoW солвер:
<pre><code>{pow_get_solver()}</code></pre>
Чтобы запустить ром и получить скриншот, отправьте форму:
<form method="post" action="/" enctype="multipart/form-data">
    <input name="work" type="text" placeholder="Результат PoW" required />
    <input name="file" type="file" required />
    <button type="submit">Запустить ром</button>
</form>
"""
    else:
        msg = """
Чтобы запустить ром и получить скриншот, отправьте форму:
<form method="post" action="/" enctype="multipart/form-data">
    <input name="file" type="file" required />
    <button type="submit">Запустить ром</button>
</form>
"""
    if request.query.get('screenshot') == "yes":
        msg += """\n<img src="/screenshot.png" style="border: 2px solid grey;"/>"""
    return msg

@app.post("/")
def run_handler():
    try:
        print(request.forms.get("work"))
        assert pow_submit_work(request.forms.get("work")) == True
    except:
        return HTTPError(400, "PoW is invalid (already used?)")

    if request.content_length is None or request.content_length > MAX_UPLOAD:
        raise HTTPError(413, "File too large")
    
    upload = request.files.get("file")
    if not upload:
        raise HTTPError(400, "No file provided")
    
    upload.save("rom.nds", overwrite=True)

    flag = os.getenv("FLAG", "flag{default_fake_flag}")
    with open("flag.txt", "w") as f:
        f.write(flag)

    env = os.environ.copy()
    if "FLAG" in env:
        del env["FLAG"]
    p = subprocess.Popen(
        ["melonDS-headless", "--jit", "rom.nds", "screenshot.png"],
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        env=env
    )
    return_code = p.wait()

    if return_code != 0:
        return f"Return code: {return_code}. Ошибка."
    else:
        return redirect("/?screenshot=yes")

@app.get("/screenshot.png")
def screenshot_handler():
    response.set_header('Content-type', 'image/png')
    response.set_header('Cache-Control', 'no-store, no-cache, must-revalidate')
    response.set_header('Pragma', 'no-cache')
    return static_file("screenshot.png", root=".")

run(app, host="0.0.0.0", port=8000)
