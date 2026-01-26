from fastecdsa.curve import P256
from fastecdsa.point import Point

from Crypto.Random.random import *

TOKEN = 1337
Q = Point(
    0xC97445F45CDEF9F0D3E05E1E585FC297235B82B5BE8FF3EFCA67C59852018192,
    0xB28EF557BA31DFCBDD21AC46E2A91E3C304F44CB87058ADA2CB815151E610046,
    P256,
)

class DGen:
    def __init__(self):
        self.Q = Q

        self.P = self.Q * TOKEN
        self.s = getrandbits(256)

        first_rands = []
        for _ in range(2):
            first_rands.append(self.get_random_seed())
        print(f"r = {first_rands}")

    def get_random_seed(self) -> int:
        s = (self.s * self.P).x
        rand = (s * self.Q).x
        self.s = s

        return rand & (2 ** (8 * 30) - 1)
