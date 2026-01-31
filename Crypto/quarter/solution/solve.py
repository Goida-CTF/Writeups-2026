from sage.all import *
from Crypto.Cipher import AES


H = QuaternionAlgebra(QQ, -1, -1)
i, j, k = H.gens()

C = "[-3645850 + 11125544*i + 160913233*j + 139894445*k, -202319279 + 136538151*i + 263558462*j + 15650808*k, -104983557 + 179922981*i + 243034582*j + 219381904*k, -31002009 - 27834531*i + 245343263*j + 93937383*k]"
flag_enc = bytes.fromhex("7fa1c8d3db311d50c5f5a26e023df5543d580ee0e4ad2c9c29105b79a78c82a8a15e6164d4d783a241950d081df2bbc6")

def xor(a, b):
    return bytes(x ^ y for x, y in zip(a, b))

C = sage_eval(C, locals={"i": i, "j": j, "k": k})
D = [c / C[0] for c in C[1:3]]

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
    if not all(0 <= x <= 256 and x >= 0 for x in b):
        continue
    P0 = H(list(b))
    quarter = P0 ** -1 * C[0]
    C = [c / quarter for c in C]
    key = bytes(x for c in C for x in list(c))
    
    cipher = AES.new(key, AES.MODE_ECB)
    print(cipher.decrypt(flag_enc))
