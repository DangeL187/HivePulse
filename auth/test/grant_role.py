import json

import requests


def grant_role_admin_to_admin(url, token):
    return _grant_role(url, 1, token, 'admin')


def grant_role_admin_to_user(url, token):
    return _grant_role(url, 2, token, 'admin')


def grant_role_operator_to_user(url, token):
    return _grant_role(url, 2, token, 'operator')


def _grant_role(url, uid, token, role):
    url = f"{url}/users/{uid}/roles"

    headers = {
        "Content-Type": "application/json",
    }

    cookie = {
        "access_token": token,
    }

    data = {
        "role": role,
    }

    try:
        response = requests.post(url, headers=headers, json=data, cookies=cookie)
        return json.loads(response.content.decode("utf-8"))
    except Exception as e:
        print('[ERROR]', e)

    return None
