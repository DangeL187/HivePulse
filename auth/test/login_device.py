import requests


def login_device(url):
    url = f"{url}/devices/login"

    headers = {
        "Content-Type": "application/json",
    }

    data = {
        "device_id": "dev-1",
        "password": "secret"
    }

    try:
        response = requests.post(url, headers=headers, json=data)
        return response.cookies['access_token'], response.cookies['refresh_token']
    except Exception as e:
        print('[ERROR]', e)

    return None, None
