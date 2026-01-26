import os
from crypto import RSA

FLAG = os.getenv("FLAG")

def main():
    cipher = RSA()
    
    c = cipher.encrypt(bytes(FLAG.encode()))
    print(f"c = {c}")
    
if __name__ == "__main__":
    main()