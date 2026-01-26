from fastecdsa.curve import P256
from fastecdsa.point import Point
from generator import TOKEN, Q
from random import seed, randbytes
from crypto import _generate_prime, RSA 

with open("./output.txt", "r") as f:
    first_rand = list(map(int, (f.readline()[5:-2]).split(",")))
    n = int(f.readline()[4:])
    c = int(f.readline()[4:])

e = 65537

def get_y_on_p256(x) -> int | None:
    y_2 = ((x**3) - (3 * x) + P256.b) % P256.p
    y = pow(y_2, (P256.p + 1) // 4, P256.p)

    if y_2 == y**2 % P256.p:
        return y
    else:
        return None

def test_A_candidate(A: Point, expected: int) -> [bool, Point | None, int | None]:
    s_test = (A * TOKEN).x
    r_test = (Q * s_test).x
    if r_test & (2 ** (8 * 30) - 1) == expected:
        return True, A, s_test
    
    return False, None, None

for bits in range(2**16+1):
    x = (bits << (8 * 30)) | (first_rand[0])
    y = get_y_on_p256(x)
    if y is None:
        continue
    A = Point(x, y, P256)
    status, A, s = test_A_candidate(A, first_rand[1])
    if status == True:
        break
    
predict_s = (s*Q*TOKEN).x
predict_r = (predict_s * Q).x & (2 ** (8 * 30) - 1)

p = _generate_prime(s=predict_r)
q = n // p

cipher = RSA(p=p, q=q)
m = cipher.decrypt(c)
print(m)