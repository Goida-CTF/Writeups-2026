import json
import numpy as np
from scipy.optimize import linprog

Q = 64.0

def quantize(x):
    return np.clip(np.floor(np.abs(x) * Q), 0, 255).astype(np.uint8)

def load():
    with open("weights.json") as f:
        d = json.load(f)
    w = d["weights"]
    return (
        np.array(w["W1"], np.float32),
        np.array(w["W2"], np.float32),
        np.array(w["W3"], np.float32),
        np.array(w["b1"], np.float32),
        np.array(w["b2"], np.float32),
        np.array(w["b3"], np.float32),
        bytes(d["encrypted_flag"]),
    )


DELTA = 1.0 / (16 * 64)

def solve_lp(W1, W2, b1, b2):
    A1 = -W1
    b1c = b1 - DELTA

    A2 = -(W2 @ W1)
    b2c = (W2 @ b1 + b2) - DELTA

    A = np.vstack([A1, A2])
    b = np.hstack([b1c, b2c])

    res = linprog(
        c=np.zeros(W1.shape[1]),
        A_ub=A,
        b_ub=b,
        bounds=(None, None),
        method="highs",
    )
    assert res.success
    return res.x.astype(np.float32)


def forward(x, W1, W2, W3, b1, b2, b3):
    h1 = np.maximum(0, W1 @ x + b1)
    h2 = np.maximum(0, W2 @ h1 + b2)
    return W3 @ h2 + b3

def main():
    W1, W2, W3, b1, b2, b3, enc = load()
    x = solve_lp(W1, W2, b1, b2)
    out = forward(x, W1, W2, W3, b1, b2, b3)

    key = quantize(out) ^ quantize(x)
    flag = bytes(enc[i] ^ key[i % len(key)] for i in range(len(enc)))

    print("[+] input:", x)
    print("[+] output:", out)
    print("[+] key:", key)
    print("[+] FLAG:", flag.decode())

if __name__ == "__main__":
    main()
