import os
import json
import jwt
from datetime import datetime
from dotenv import load_dotenv

load_dotenv(dotenv_path="../../.env")

PORT = os.getenv("PORT", "8080")
JWT_SECRET = os.getenv("JWT_SECRET")

if not JWT_SECRET:
    raise ValueError("No JWT_SECRET in .env")

def generate_token(user_id, username, role, secret):
    exp = int(datetime(2050, 1, 1).timestamp())
    payload = {
        "lid": user_id,
        "usr": username,
        "rol": role,
        "exp": exp
    }
    return jwt.encode(payload, secret, algorithm="HS256")

# Valid tokens
user1_token = generate_token("1", "user1", "user", JWT_SECRET)
user2_token = generate_token("2", "user2", "user", JWT_SECRET)
admin_token = generate_token("3", "admin", "admin", JWT_SECRET)


# Invalid tokens:

# Username missing
incomplete_token = (
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."
    "eyJsaWQiOiIxIiwicm9sIjoidXNlciIsImV4cCI6MjUyNDYwNDQwMH0."
    "y3_zzKHCcwmP70-xrp7-rMyoffF7My8FFVW_FVVMtgs."
)

# Setting the role to admin, without changing the signature
fake_token = (
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."
    "eyJsaWQiOiIxIiwidXNyIjoidXNlcjEiLCJyb2wiOiJhZG1pbiIsImV4cCI6MjUyNDYwNDQwMH0."
    "Z7IxtCnrmO6oYqMKipSzbkq3CVvfYiQTipoHk3qAcDA"
)

# Postman enviroment variables:
values = [
    {"key": "base_url", "value": f"http://localhost:{PORT}", "type": "default", "enabled": True},
    {"key": "user1_valid_token", "value": user1_token, "type": "default", "enabled": True},
    {"key": "user2_valid_token", "value": user2_token, "type": "default", "enabled": True},
    {"key": "admin_valid_token", "value": admin_token, "type": "default", "enabled": True},
    {"key": "incomplete_token", "value": incomplete_token, "type": "default", "enabled": True},
    {"key": "fake_token", "value": fake_token, "type": "default", "enabled": True}
]

# Postman FULL environment JSON
env = {
    "id": "592af9d5-b95a-407b-ae6f-15aa35b4ffeb",
    "name": "tests",
    "values": values,
    "_postman_variable_scope": "environment",
    "_postman_exported_at": datetime.now().isoformat() + "Z",
    "_postman_exported_using": "Postman/11.49.1"
}

# Saving into the postman env file
with open("./postman-files/tests.postman_enviroment.json", "w") as file:
    json.dump(env, file, indent=4)