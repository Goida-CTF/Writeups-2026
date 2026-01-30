# EFN
Сложность: `Сложно`

Описание:

```
[Сообщение от Пэма Бонди, Министерство Юстиции США]:
Bro, sorry to bother you, but our secret files related to the Epstein case were stolen! There was a PDF document and a track that played frequently on the island.
Fucking Chinese routers, I told them not to install them!! I swear the fucking Chinese stole files from our file storage through that fucking router!!
The admins dumped wireless traffic from the DOJ office, check this out.
Also to build a full picture of an incident we provide you original sample of stolen archive, but it's encrypted AND DON'T EVEN TRY TO DECRYPT THIS!
Bro, I know you're the best Russian hacker, please help us!! America will be proud of you!!
_______

Необходимо найти имя хакера - это одна из частей флага.
Формат флага: goidactf{NameOfHacker_...}
```


## Теги
`#forensics` `#network` `#files`

## Выдать участникам
[officewifidump.cap](https://github.com/k10nex/goidactf2026quals-tasks/blob/main/Forensics/EFN/officewifidump.cap)

[UnpublishedEpsteinFiles.zip](https://github.com/k10nex/goidactf2026quals-tasks/blob/main/Forensics/EFN/UnpublishedEpsteinFiles.zip)

## Решение
В задании выдается дамп беспроводного трафика прямиком из офиса Мин.Юстиции и архив с украденными файлами. Начнем решение с поиска имени хакера.

Откроем дамп трафика. Нашему взгяду предстут фреймы протокола 802.11, которые в текущем виде никакой полезной информации нам не дают (кроме SSID точки доступа, но даже его мы позже найдем другим способом). Для решения этой задачи нужно дешифровать этот трафик.

Перейдя в иерархию протоколов заметим, что в трафике присутствуют пакеты аутентификации. Перейдем к этим пакетикам и увидим, что было совершено четырехстороннее рукопожатие, что означает, что мы сможем дешифровать трафик если найдем пароль для этой WiFi-сети.

Вспоминаем светлые школьные денёчки, когда каждый хотел взломать школьный вай-фай чтобы не платить по 150 рублей за гигабайт интернета - а вместе с ними вспоминаем про такую штуку как `aircrack-ng`. В чистом виде он нам не понадобится, но знать, как работает аутентификация в WiFi-сетях для решения данной задачки было бы полезно.

Прочитав [некоторые статейки](https://www.opennet.ru/tips/3025_wpa_wpa2_wifi_aircrackng_hashcat.shtml) и подтянув матчасть возвращаемся к задачке - устанавливаем `hcxpcapngtool`, достаем брутабельный хэш и отправляем его брутиться в `hashcat`:

<img width="1920" height="1080" alt="cap-processed" src="https://github.com/user-attachments/assets/1a6199dc-00fc-4245-a084-01f7526d93dd" />

<img width="1115" height="628" alt="handshake-cracked" src="https://github.com/user-attachments/assets/84b16e65-5861-4245-ac90-3b14c46d2add" />

Получаем пароль от WiFi-сети. Смотрим, как нам дешифровать трафик, дешифруем и теперь видим вполне понятные всем протоколы и пакеты:

<img width="1920" height="1032" alt="dec-wifi-filter" src="https://github.com/user-attachments/assets/b13a9c04-d3ba-41a8-ae25-560f053236bc" />

А среди них видим и работу MDNS на айпаде злоумышленника (который спокойно рассылает всем свое имя по мультикасту даже с включенным частным MAC - безопасненько, Крейг, молодец). Устройства Apple всегда берут себе название формата `iDevice (AppleID-FirstName)`, соответственно, то, что мы ищем нам прямым текстом показал MDNS.

> *Имя злоумышленника: `SunXuiVcay`*

Идем дальше. Пэм в своем сообщении говорил нам не пытаться расшифровать архив, но мы его, конечно же, не послушаем.

Как многие заметили - архив не получается забрутить дефолтными словарями (спойлер - никогда никакими и не получилось бы). Здесь нужно вспомнить про одну [очень интересную атаку на ZipCrypto](https://link.springer.com/chapter/10.1007/3-540-60590-8_12) - `Known Plaintext Attack`. Эта атака применима только к архивам, зашифрованным с использованием ZipCrypto и не использующим сжатие файлов.

Наш архив как раз такой, это можно заметить по размеру файлов в нем (файлы становятся даже больше из-за добавления доп. информации к ним):

<img width="782" height="565" alt="image" src="https://github.com/user-attachments/assets/23c18e0c-41ef-4c18-9bed-08d9770d5a48" />

Ищем [и находим](https://github.com/kimci86/bkcrack) инструмент, реализующий данную атаку. Прочитав инструкцию нам становится понятно, что для того, чтобы атака была возможной, нам нужно знать плейнтекст какого-нибудь файла.

В архиве есть цифровая подпись в формате XML. Хедер у XML по стандарту всегда одинаковый - это `<?xml version="1.0" encoding="UTF-8"?>`. Читаем доку к bkcrack, сохраняем плейнтекст в файл, пишем команду, брутим и на выходе получаем расшифрованный архив:

<img width="1088" height="272" alt="archive-bkcracked" src="https://github.com/user-attachments/assets/fb3c24fc-aec9-4e7b-b536-d17314f4f493" />

В расшифрованном архиве сначала посмотрим XML-файл, так как это единственный файл текстового формата. В некоторых полях электронной подписи найдем base64-данные, а в самом сертификате обнаружим интересную строчку:

<img width="1346" height="527" alt="image" src="https://github.com/user-attachments/assets/72779d86-cc1e-45b3-8d2e-2d8f73574340" />


> *Что-то, что нам нужно взять: `J8ZJoFsJMn2CP62wuYvZM49DbWEXNaTx`*


Продолжим. PDF-файл оказался зашифрованным, но хорошо, что мы умеем пользоваться [онлайн-инструментами](https://hashes.com/ru/johntheripper/pdf2john) и брутить по rockyou - быстренько восстанавливаем пароль и открываем файлы Эпштейна.

<img width="1115" height="628" alt="image" src="https://github.com/user-attachments/assets/02c40e34-e593-425c-9ad4-a9ecc3bcae3c" />

> *Пароль от архива: `03542332238`*

Воспользовавшись подсказкой из названия файла переходим на страницу № sixseven. Полное копирование этой страницы не приводит нас к успеху:

<img width="530" height="607" alt="image" src="https://github.com/user-attachments/assets/3ded3adf-18ba-48de-a3ee-864be98a27c3" />

Почитаем контекст и заметим, что один из блоков закрывает не приватную информацию о ком-то, а информацию о каком-то событии, которая вообще не должна была закрываться цензорами. Пробуем выделить именно этот кусочек:

<img width="740" height="487" alt="image" src="https://github.com/user-attachments/assets/faeb0860-f1f7-4cb8-b5e2-267ac3bda978" />

Получаем вторую часть флага.

> *Вторая часть флага: `C0MPR0M153D_3p5731n`*

*В целом, можно было полностью вытащить текст из pdf-файла и погрепать флаг. Но это больше для тех, кто не шарит на сикссевен или не умеет видеть банальные подсказки.*

Из непроверенных нами остался только аудиофайл. Его музыкальное содержание безусловно прекрасно, но нас больше интересует флаг, поэтому попробуем разобраться.

После того, как анализ вейвформа и спектрограммы ни к чему не привели, вспоминаем, какая программа по стандарту возвращает стегоконтейнеры в формате WAV. Это DeepSound, довольно популярная программка на цтфках. Пихаем туда нашу прекрасную песню и видим запрос пароля.

Вот и пригодилось нам то, что мы получили из ЭЦП. Вставляем строку `J8ZJoFsJMn2CP62wuYvZM49DbWEXNaTx` в DeepSound как пароль, достаем файлы и видим финальную часть флага:

<img width="329" height="171" alt="image" src="https://github.com/user-attachments/assets/07f4230b-ab15-46df-92ad-15077176826f" />

> *Третья часть флага: `_f1le5_nowaaaaay!!}`*

Собираем и сдаем флаг.

## Флаг
    goidactf{SunXuiVcay_C0MPR0M153D_3p5731n_f1le5_nowaaaaay!!}
