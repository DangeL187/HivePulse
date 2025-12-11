import requests


def login_admin(url):
    return _login(url, "admin", "admin")


def login_user(url):
    return _login(url, "test@example.com", "secret")


def _login(url, email, password):
    url = f"{url}/users/login"

    headers = {
        "Content-Type": "application/json",
    }

    data = {
        "email": email,
        "password": password
    }

    try:
        response = requests.post(url, headers=headers, json=data)
        return response.cookies['access_token']
    except Exception as e:
        print('[ERROR]', e)

    return None
