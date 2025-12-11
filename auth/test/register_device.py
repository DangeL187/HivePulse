import requests


def register_device(url, token):
    for i in range(500):
        _register_device(url, token, i+1)


def _register_device(url, token, device_id):
    url = f"{url}/devices/register"

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
        pass
        # print('[ERROR]', e)  # raises exception if device exists

    return None
