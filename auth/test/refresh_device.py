import requests


def refresh_device(url, token):
    url = f"{url}/devices/refresh"

    headers = {
        "Content-Type": "application/json",
    }

    cookie = {
        "refresh_token": token,
    }

    try:
        response = requests.post(url, headers=headers, cookies=cookie)
        return response.cookies['access_token']
    except Exception as e:
        print('[ERROR]', e)

    return None
