import requests


def register_device(token):
    for i in range(1000):
        _register_device(token, i+1)

def _register_device(token, device_id):
    url = "http://localhost:8000/devices/register"

    headers = {
        "Content-Type": "application/json",
    }

    data = {
        "device_id": "dev-" + str(device_id),
        "password": "secret"
    }

    cookie = {
        "access_token": token
    }

    try:
        response = requests.post(url, headers=headers, json=data, cookies=cookie)
        return response.cookies['access_token'], response.cookies['refresh_token']
    except Exception as e:
        print('[ERROR]', e)

    return None
