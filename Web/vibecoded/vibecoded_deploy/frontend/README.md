# frontend

build prod frontend with
```bash
tailwindcss -i ./src/input.css -o ./dist/styles.css --minify
html-minifier --collapse-whitespace --remove-comments --minify-css true --minify-js true -o ./dist/index.html ./src/index_stage.html
```
