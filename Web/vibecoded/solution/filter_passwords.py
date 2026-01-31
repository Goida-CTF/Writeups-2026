import string


allowed_chars = (
    string.digits
    + string.ascii_lowercase
    + string.ascii_uppercase
    + string.punctuation
)


def is_accepted_password(password: str) -> bool:
    return (
        8 <= len(password) <= 32
        and any(char in string.digits for char in password)
        and any(char in string.ascii_lowercase for char in password)
        and any(char in string.ascii_uppercase for char in password)
        and any(char in string.punctuation for char in password)
        and all(char in allowed_chars for char in password)
    )


def main() -> None:
    with (
        open("/usr/share/wordlists/rockyou.txt", "rb") as fr,
        open("./filtered.txt", "a") as fw,
    ):
        for i, line in enumerate(fr):
            password: str
            try:
                password = line.decode()[:-1]
            except UnicodeDecodeError:
                continue

            if not is_accepted_password(password):
                continue

            fw.write(f"{password}\n")
            print(f"{i}: {password}")


if __name__ == "__main__":
    main()
