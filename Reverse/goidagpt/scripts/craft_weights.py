import json
import subprocess
import sys

import numpy as np
from scipy.linalg import qr

IO_DIM = 8
HIDDEN_DIM = 16
FLAG = b"goidactf{G01D4_GPT_T0T4L_N3URO_SLOP}"
SEED = 0x601DA


def craft_network():
    np.random.seed(SEED)

    print("[*] Crafting network weights...")

    # https://en.wikipedia.org/wiki/QR_decomposition
    # дает нам ортогональные матрицы
    Q1, _ = qr(np.random.randn(HIDDEN_DIM, HIDDEN_DIM))
    Q2, _ = qr(np.random.randn(HIDDEN_DIM, HIDDEN_DIM))

    W1 = Q1[:, :IO_DIM].astype(np.float32)
    W2 = Q2.astype(np.float32)
    W3 = (Q2 @ Q1[:, :IO_DIM]).T.astype(np.float32)

    prod = W3 @ W2 @ W1
    i_err = np.max(np.abs(prod - np.eye(IO_DIM)))
    print(f"ошибка от I: {i_err:.6f}")
    assert i_err < 0.01, "Not close to I"

    secret_input = np.random.randn(IO_DIM).astype(np.float32) * 5
    print(f"инпут для ключа: {secret_input}")

    # собираем коэфы чтобы ReLU давал положительные значения
    delta = np.ones(HIDDEN_DIM, dtype=np.float32) * 0.1
    b1 = -W1 @ secret_input + delta  # h1 = delta
    b2 = -W2 @ delta + delta  # h2 = delta

    Bdelta = W3 @ W2 @ b1 + W3 @ b2

    b3 = np.zeros(IO_DIM, dtype=np.float32)
    output = secret_input + Bdelta + b3

    #
    # В нашем бинаре мы далаем: (u8)abs(output) XOR (u8)abs(input)
    output_bytes = np.abs(output).astype(np.uint8)
    input_bytes = np.abs(secret_input).astype(np.uint8)
    k = output_bytes ^ input_bytes
    print(f"output: {output}")
    print(f"ключ: {k}")
    encrypted = bytes(f ^ int(k[i % IO_DIM]) for i, f in enumerate(FLAG))
    print(f"Encrypted flag: {encrypted.hex()}")
    # Проверяем что ключ работает, не вырожденный и тд
    decrypted = bytes(e ^ int(k[i % IO_DIM]) for i, e in enumerate(encrypted))
    assert decrypted == FLAG, "can't decrypt"

    return {
        "weights": {
            "W1": W1.tolist(),
            "W2": W2.tolist(),
            "W3": W3.tolist(),
            "b1": b1.tolist(),
            "b2": b2.tolist(),
            "b3": b3.tolist(),
        },
        "secret_input": secret_input.tolist(),
        "encrypted_flag": list(encrypted),
        "flag_length": len(FLAG),
    }


def main():
    data = craft_network()

    with open("weights.json", "w") as f:
        json.dump(data, f, indent=2)

    # так как мы хотим comptime веса в бинаре то нельзя парсить json
    subprocess.run([sys.executable, "scripts/generate_zig_weights.py"], check=True)

    print("Done")


if __name__ == "__main__":
    main()
