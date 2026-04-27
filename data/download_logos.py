import requests
import os

teams = [
    "atl", "bos", "bkn", "cha", "chi", "cle", "dal", "den", "det", "gsw",
    "hou", "ind", "lac", "lal", "mem", "mia", "mil", "min", "no", "ny",
    "okc", "orl", "phi", "phx", "por", "sac", "sa", "tor", "utah", "wsh"
]

output_dir = "web/static/images/logos"
os.makedirs(output_dir, exist_ok=True)

for team in teams:
    url = f"https://a.espncdn.com/i/teamlogos/nba/500/{team}.png"
    response = requests.get(url)
    if response.status_code == 200:
        with open(f"{output_dir}/{team}.png", "wb") as f:
            f.write(response.content)
        print(f"Downloaded: {team}")
    else:
        print(f"Failed: {team} (status {response.status_code})")