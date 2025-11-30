import json

import requests


def revoke_role_admin_from_user(uid, token):
    return _revoke_role_from_user(uid, token, "admin")


def revoke_role_operator_from_user(uid, token):
    return _revoke_role_from_user(uid, token, "operator")


def _revoke_role_from_user(uid, token, role):
    url = f"http://localhost:8000/users/{uid}/roles/{role}"

    headers = {
        "Content-Type": "application/json",
    }

    cookie = {
        "access_token": token,
    }

    try:
        response = requests.delete(url, headers=headers, cookies=cookie)
        return json.loads(response.content.decode("utf-8"))
    except Exception as e:
        print('[ERROR]', e)

    return None
