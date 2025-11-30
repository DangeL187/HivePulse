import json

import requests


def get_roles_admin(token):
    return _get_roles(1, token)


def get_roles_user(token):
    return _get_roles(2, token)


def _get_roles(uid, token):
    url = f"http://localhost:8000/users/{uid}/roles"

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
