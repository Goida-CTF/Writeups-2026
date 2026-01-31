from sage.all import *
from Crypto.Cipher import AES


H = QuaternionAlgebra(QQ, -1, -1)
i, j, k = H.gens()

C = "[-15197553861 + 71114182955*i - 41801438776*j - 10172845760*k, -22328311008 + 82523670788*i - 43598195417*j - 27140296104*k]"
flag_enc = bytes.fromhex("ce07f0b8e3c83fab445025fd77435b2636e2ce8db9cba2481ef6bc62917f739ab0e1b31a6e9655dffa44a34a3bfcbeee")

def xor(a, b):
    return bytes(x ^ y for x, y in zip(a, b))

C = sage_eval(C, locals={"i": i, "j": j, "k": k})
D = [c / C[0] for c in C[1:]]

p = next_prime(2 ** 64)
Hr = QuaternionAlgebra(GF(p), -1, -1)
i, j, k = Hr.gens()

D = [Hr(list(x)) for x in D]

def q2mat(q):
    return matrix([list(Hr(q * m)) for m in (1, i, j, k)])

B = matrix.block(ZZ,
    [[q2mat(q) for q in D] + [q2mat(1)]]
).stack(
    (matrix.identity(len(D) * 4) * p).augment(matrix.zero(len(D) * 4, 4))
)

B = B.LLL()

for row in B:
    b = row[-4:]
    if b[0] < 0:
        b = -b
    if not all(0 <= x <= 256 ** 2 and x >= 0 for x in b):
        continue
    P0 = H(list(b))
    quarter = P0 ** -1 * C[0]
    C = [c / quarter for c in C]
    key = bytes(b for c in C for x in list(c) for b in [int(x) >> 8, int(x) & 0xFF])
    
    cipher = AES.new(key, AES.MODE_ECB)
    print(cipher.decrypt(flag_enc))
