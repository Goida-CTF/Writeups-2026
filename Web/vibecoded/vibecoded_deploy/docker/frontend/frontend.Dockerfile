FROM nginx:1.29.3-alpine3.22
ARG HCAPTCHA_SITEKEY
COPY frontend/dist /tmp/dist

RUN test -n "$HCAPTCHA_SITEKEY" || (echo "HCAPTCHA_SITEKEY is required" && exit 1)

RUN ESCAPED_KEY="$(printf '%s' "$HCAPTCHA_SITEKEY" | sed 's/[\/&]/\\&/g')" \
  && sed "s/__HCAPTCHA_SITEKEY__/${ESCAPED_KEY}/g" /tmp/dist/index.html > /usr/share/nginx/html/index.html \
  && mv /tmp/dist/styles.css /usr/share/nginx/html/.

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
