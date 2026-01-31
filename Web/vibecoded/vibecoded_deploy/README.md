# Навайбкодили

Таск деплоится из Docker Compose. Нужно при билде фронта проставить аргумент `HCAPTCHA_SITEKEY`, а для бэкенда проставить environment vars `HCAPTCHA_SECRET` и `FLAG` (брать из [`../solution/flag.txt`](../solution/flag.txt)). Все остальные environment vars проставлены.
