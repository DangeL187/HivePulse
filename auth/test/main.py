from get_roles import get_roles_admin, get_roles_user
from grant_role import grant_role_admin_to_admin, grant_role_admin_to_user, grant_role_operator_to_user
from login import login_admin, login_user
from login_device import login_device
from refresh_device import refresh_device
from register import register_user
from register_device import register_device
from revoke_role import revoke_role_admin_from_user, revoke_role_operator_from_user

# URL = "http://localhost:8000"  # local
URL = "http://localhost:30080"  # k8s


def check(condition: bool, label: str = ""):
    status = "[PASSED]" if condition else "[FAILED]"
    color = "\033[92m" if condition else "\033[91m"
    reset = "\033[0m"
    if label != "":
        print(f"{label} ", end='')
    print(f"{color}{status}{reset}")


def permissions_check_user_not_admin(admin_token, user_token):
    print('\n[*] Permissions check (user is not admin)...')

    check('admin' in get_roles_admin(URL, admin_token)["roles"], 'get roles admin')

    res = get_roles_user(URL, user_token)
    check("roles" in res and not res["roles"], 'get roles user')

    res = grant_role_admin_to_admin(URL, user_token)
    check("error" in res and res["error"] == 'forbidden', 'grant role "admin" to admin by user')

    res = grant_role_admin_to_user(URL, user_token)
    check("error" in res and res["error"] == 'forbidden', 'grant role "admin" to user by user')

    res = revoke_role_admin_from_user(URL, 2, user_token)
    check("error" in res and res["error"] == 'forbidden', 'revoke role "admin" from user by user')

    res = revoke_role_admin_from_user(URL, 2, admin_token)
    print(res)
    check("message" in res and res["message"] == 'role has been revoked', 'revoke role "admin" from user by admin')


def permissions_check_user_admin(admin_token, user_token):
    print('\n[*] Permissions check (user is admin)...')

    check('admin' in get_roles_admin(URL, admin_token)["roles"], 'get roles admin')
    check('admin' in get_roles_user(URL, user_token)["roles"], 'get roles user')

    res = grant_role_admin_to_admin(URL, user_token)
    check("error" in res and res["error"] == 'user already has this role', 'grant role "admin" to admin by user')

    res = grant_role_admin_to_user(URL, user_token)
    check("error" in res and res["error"] == 'user already has this role', 'grant role "admin" to user by user')


def main():
    print('\n[*] Register and login...')

    register_user(URL)

    admin_token = login_admin(URL)
    check(admin_token is not None, 'login admin')

    user_token = login_user(URL)
    check(user_token is not None, 'login user')

    permissions_check_user_not_admin(admin_token, user_token)

    print('\n[*] Permissions granting...')

    res = grant_role_admin_to_admin(URL, admin_token)
    check("error" in res and res["error"] == 'user already has this role', 'grant role "admin" to admin by admin')

    res = grant_role_admin_to_user(URL, admin_token)
    check("message" in res and res["message"] == 'role has been granted', 'grant role "admin" to user by admin')

    res = grant_role_operator_to_user(URL, admin_token)
    check("message" in res and res["message"] == 'role has been granted', 'grant role "operator" to user by admin')

    permissions_check_user_admin(admin_token, user_token)

    print('\n[*] Device register and login... (it might take some time)')

    register_device(URL, user_token)

    device_access_token, device_refresh_token = login_device(URL)
    check(device_access_token is not None and device_refresh_token is not None, 'login device')

    device_access_token = refresh_device(URL, device_access_token)
    check(device_access_token is None, 'failed to refresh device with access token')

    device_access_token = refresh_device(URL, device_refresh_token)
    check(device_access_token is not None, 'refresh device with refresh token')

    print('\n[*] Permissions revoking...')

    res = revoke_role_admin_from_user(URL, 2, user_token)
    check("message" in res and res["message"] == 'role has been revoked', 'revoke role "admin" from user by user')

    res = revoke_role_admin_from_user(URL, 2, admin_token)
    check("message" in res and res["message"] == 'role has been revoked', 'revoke role "admin" from user by admin')

    res = revoke_role_operator_from_user(URL, 2, admin_token)
    check("message" in res and res["message"] == 'role has been revoked', 'revoke role "operator" from user by admin')

    permissions_check_user_not_admin(admin_token, user_token)


if __name__ == '__main__':
    main()
