from typing import List, Tuple
from sage.all import *
import random
import gmpy2

FLAG = b'goidactf{b4efd05cdc57cd5656192303b555b193}'

def generate_coeffs() -> Tuple[List[int], List[int]]:
    while True:
        coeffs = [[0] * 4 for _ in range(4)]
        m = matrix([[0] * 4 for _ in range(4)])
        for i in range(4):
            for j in range(4):
                coeffs[i][j] = random.randrange(0, 17)
                m[i, j] = gmpy2.next_prime(coeffs[i][j])
        if m.is_invertible():
            break
    return coeffs, m

def marshal_coeffs(coeffs: List[int]) -> str:
    return f"[{','.join(str(c) for c in coeffs)}]"

def marshal_coeffs_and_target(coeffs: List[int], target: int) -> str:
    return f"({marshal_coeffs(coeffs)}, \"{target * '1'}\")"

def marshal_quad(coeffs: List[List[int]], targets: List[int]) -> str:
    res = []
    for c, t in zip(coeffs, targets):
        res.append(marshal_coeffs_and_target(c, t))

    return f"[{','.join(res)}]"

def main():
    flag = FLAG[:]
    flag_padded = flag + b'\0' * ((4 - len(flag)) % 4)
    res = []
    for i in range(0, len(flag_padded), 4):
        coeffs, m = generate_coeffs()
        v = vector(ZZ, list(flag_padded[i:i+4]))
        res.append((coeffs, m * v))

    with open("coeffs.txt", "w") as f:
        print(res, file=f)
    with open("src/coeffs.rs", "w") as f:
        print(f"[{','.join(marshal_quad(coeffs, targets) for coeffs, targets in res)}]", file=f)
        
if __name__ == "__main__":
    main()