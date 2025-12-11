import json

import requests


def get_roles_admin(url, token):
    return _get_roles(url, 1, token)


def get_roles_user(url, token):
    return _get_roles(url, 2, token)


def _get_roles(url, uid, token):
    url = f"{url}/users/{uid}/roles"

    headers = {
        "Content-Type": "application/json",
    }

    cookie = {
        "access_token": token,
    }

    try:
        response = requests.get(url, headers=headers, cookies=cookie)
        return json.loads(response.content.decode("utf-8"))
    except Exception as e:
        print('[ERROR]', e)

    return None
