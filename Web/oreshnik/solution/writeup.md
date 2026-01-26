# Oreshnik Writeup

https://medium.com/@mark_huber/decoding-the-jwt-anomaly-when-changing-a-tokens-last-character-doesn-t-break-verification-d6ab68627afb

Используя данную аномалию при декодировании base64 в JWT, можно обойти проверку на revoked токен и получить доступ к админ панели с флагом.

Заходим в /revoked, ищем токен админа, перебираем замену последнего символа на все возможные значения, пока не найдем рабочий токен.

Делаем запрос на /admin с полученным токеном.

Подробный сплойт [solve.py](./solve.py)