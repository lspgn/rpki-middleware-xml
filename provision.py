#!/usr/bin/env python3

import requests
import json

checkssl = False

krill1 = "https://localhost:3001"
krill2 = "https://localhost:3002"

apipath = "/api/v1"

token1 = "test"
token2 = "test"

ca1 = "ta"

rpki_notify = "https://krill2:3000/rrdp/notification.xml"
base_uri = "rsync://krill2:3000/repo/{}/".format(ca1)

service_uri = "https://middleware:3000/rfc8181/{}".format(ca1)
# service_uri = "https://krill2:3000/rfc8181/ta"

print("Fetching Repo Request")
data = requests.get(
    "{}{}/cas/{}/repo/request.json".format(krill1, apipath, ca1),
    headers={"Authorization": "Bearer {}".format(token1)},
    verify=checkssl,
)
repo_request = data.json()
pub_handle = repo_request.get("publisher_handle")

print("Fetching Publisher Response")

data = requests.post(
    "{}{}/publishers".format(krill2, apipath),
    headers={
        "Authorization": "Bearer {}".format(token2),
        "Content-Type": "application/json",
    },
    verify=checkssl,
    data=json.dumps(
        {"publisher_handle": pub_handle, "id_cert": repo_request.get("id_cert"),}
    ),
)
print(data)
publisher_response = data.json()
if data.status_code != 200 and publisher_response.get("code") == 2102:
    print("Handle already exists")
    data = requests.get(
        "{}{}/publishers/{}".format(krill2, apipath, pub_handle),
        headers={"Authorization": "Bearer {}".format(token2),},
        verify=checkssl,
    )
    publisher_response = data.json()

print(publisher_response)

print("Adding publisher to repo")

data = requests.post(
    "{}{}/cas/{}/repo".format(krill1, apipath, ca1),
    headers={
        "Authorization": "Bearer {}".format(token1),
        "Content-Type": "application/json",
    },
    verify=checkssl,
    data=json.dumps(
        {
            "rfc8181": {
                "publisher_handle": "my-pub",
                "id_cert": publisher_response.get("id_cert"),
                "repo_info": {"rpki_notify": rpki_notify, "base_uri": base_uri,},
                "service_uri": service_uri,
            }
        }
    ),
)
if data.status_code != 200:
    print(data.json())
