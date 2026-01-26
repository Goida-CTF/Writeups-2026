# Решение
### Обзор решения
1 минута тупого фаззинга afl++ по гайду от кодекса.

### Детали решения
Билдим afl++: https://github.com/AFLplusplus/AFLplusplus/blob/stable/docs/INSTALL.md
Билдим melonDS-headless с поддержкой afl++:
```
CC=afl-cc CXX=afl-c++ cmake -B build-afl -S . -DCMAKE_BUILD_TYPE=Release -DBUILD_QT_SDL=OFF -DBUILD_HEADLESS=ON -DENABLE_OGLRENDERER=OFF
cmake --build build-afl --target melonDS-headless
```
Фаззим 1 минуту на 6 exec/sec:
```
afl-fuzz -i seeds -o out -- ./melonDS-headless/build-afl/melonDS-headless --frames 1 @@
```
Если все прошло успешно, в out/crashes будет решение.
Приложен `solution.nds`