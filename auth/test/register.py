import requests


def register_user(url):
    url = f"{url}/users/register"

    headers = {
        "Content-Type": "application/json",
    }

    data = {
        "email": "test@example.com",
        "full_name": "John Doe",
        "password": "secret"
    }

    try:
        response = requests.post(url, headers=headers, json=data)
        return response.content
    except Exception as e:
        print('[ERROR]', e)

    return None
