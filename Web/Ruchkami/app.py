from socket import timeout
from flask import Flask, request, jsonify, render_template, render_template_string
import uuid
from selenium import webdriver
import requests
from time import sleep, time
import os
import platform

app = Flask(__name__)

app.jinja_env.globals['__import__'] = __import__

admin_login = str(uuid.uuid4())
admin_logins = [admin_login]
print(admin_logins)


@app.route('/', methods=['GET'])
def index():
    global admin_logins
    return render_template('index.html')


@app.route('/get_url', methods=['POST'])
def get_url():
    global admin_logins
    data = request.get_json()
    if not data or 'url' not in data:
        return jsonify({'error': 'No URL provided'}), 400

    url = data['url']
    
    chrome_options = webdriver.ChromeOptions()
    chrome_options.add_argument('--headless')
    chrome_options.add_argument('--no-sandbox')
    chrome_options.add_argument('--disable-dev-shm-usage')
    chrome_options.add_argument("start-maximized")
    chrome_options.add_argument("--disable-features=InsecureFormSubmissionMixedContentWarning")
    chrome_options.add_argument("--ignore-certificate-errors")
    chrome_options.add_argument("--allow-insecure-localhost")
    
    chromium_path = os.getenv('CHROMIUM_PATH')
    if chromium_path:
        chrome_options.binary_location = chromium_path
    elif platform.system() == 'Linux' or os.path.exists('/usr/bin/chromium'):
        chrome_options.binary_location = '/usr/bin/chromium'
        
    driver = webdriver.Chrome(options=chrome_options)
    response = requests.get(url)
    status_code = response.status_code
    print(status_code)

    # public_host = os.getenv('PUBLIC_HOST', 'tasks.goidactf.ru')
    public_host = 'tasks.goidactf.ru'
    if ":" in public_host:
        public_domain = public_host.split(':')[0]
    else:
        public_domain = public_host
    print(public_host)

    driver.get(f"http://{public_host}")
    cookie = {
        "name": "admin_login",
        "value": admin_login,
        "domain": public_domain,
        "path": "/",
        "httpOnly": True,
        "secure": False
    }
    driver.add_cookie(cookie)
    print(driver.get_cookies())

    driver.get(url)
    sleep(5)

    driver.quit()
    return jsonify({'status_code': status_code}), 200


@app.route('/create_admin', methods=['POST'])
def create_admin():
    global admin_logins
    if request.method == 'POST':
        login = request.cookies.get('admin_login')
        if login in admin_logins:

            new_admin_raw = request.form.get('new_admin_login')

            new_admin_decoded = requests.utils.unquote(new_admin_raw)

            WAF_BLACKLIST = "[]'()=_"
            for char in WAF_BLACKLIST:
                if char in new_admin_decoded:
                    return jsonify({'error': 'Request blocked by WAF'}), 403

            ESCAPE_BLACKLIST = ['\\x', '\\u', '\\U', '\\0', '\\1', '\\2', '\\3', '\\4', '\\5', '\\6', '\\7']
            for esc in ESCAPE_BLACKLIST:
                if esc in new_admin_decoded:
                    return jsonify({'error': 'Request blocked by WAF'}), 403

            new_admin = new_admin_decoded.encode('utf-8').decode('unicode_escape')
            admin_logins.append(new_admin)

            templ = 'Admin {} created'.replace("{}", new_admin)

            response_message = render_template_string(templ)

            return jsonify({'message': response_message}), 201

        else:
            return jsonify(
                {'message': 'Admin creation failed (missing/incorrect "admin_login" cookie)'}), 401  # Forbidden

    return jsonify({'message': 'This endpoint only accepts POST requests'}), 405


@app.route('/admin', methods=['GET'])
def admin_page():
    global admin_logins
    return render_template('admin.html')


if __name__ == '__main__':
    import os
    cert_dir = os.path.join(os.path.dirname(__file__), 'certs')
    ssl_context = (
        os.path.join(cert_dir, 'server.crt'),
        os.path.join(cert_dir, 'server.key')
    )
    app.run(host='0.0.0.0', port=31337, ssl_context=ssl_context)
