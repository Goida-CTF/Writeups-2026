from sage.all import *
import os
import random
import secrets
from Crypto.Cipher import AES
from Crypto.Util.Padding import pad


FLAG = os.getenvb(b"FLAG", b"flag{test_flag}")

H = QuaternionAlgebra(QQ, -1, -1)
i, j, k = H.gens()

def random_quaternion(bound = 1000000):
    return sum(random.randint(-bound, bound) * x for x in (1, i, j, k))

def byte_quaternion(bts):
    return sum(x * y for x, y in zip(bts, (1, i, j, k)))

key = secrets.token_bytes(16)
blocks = [byte_quaternion(key[i:i+4]) for i in range(0, len(key), 4)]
quater = random_quaternion()

print([x * quater for x in blocks])

cipher = AES.new(key, AES.MODE_ECB)
print(cipher.encrypt(pad(FLAG, cipher.block_size)).hex())

