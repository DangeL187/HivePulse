import requests


def login_admin():
    return _login("admin", "admin")


def login_user():
    return _login("test@example.com", "secret")


def _login(email, password):
    url = "http://localhost:8000/users/login"

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
