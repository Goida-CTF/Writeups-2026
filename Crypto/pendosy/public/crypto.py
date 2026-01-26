from generator import DGen
from Crypto.Util.number import getPrime, long_to_bytes
from random import seed, randbytes

class RSA:
    def __init__(self, p: int = 0, q: int = 0):
        if p == 0 or q == 0:
            generator = DGen()
            p = _generate_prime(s=generator.get_random_seed())
            q = _generate_prime(s=generator.get_random_seed())

        self.n = p * q
        print(f"n = {self.n}")

        self.e = 65537
        self.d = pow(self.e, -1, (p - 1) * (q - 1))
    
    def encrypt(self, m: bytes) -> int:
        m = int.from_bytes(m, byteorder='big')
        c = pow(m, self.e, self.n)

        return c

    def decrypt(self, c: int) -> bytes:
        m = pow(c, self.d, self.n)

        return long_to_bytes(m)

def _generate_prime(s: int) -> int:
    seed(s)
    return getPrime(512, randfunc=randbytes)