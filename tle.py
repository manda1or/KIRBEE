import requests
import json 
from dotenv import load_dotenv

load_dotenv()

# import credentials from .env
username = os.getenv('SPACE_USER')
password = os.getenv('SPACE_PASSWORD')

def fetch_tles():
    url = "https://www.space-track.org/ajaxauth/login"
    query = "https://www.space-track.org/basicspacedata/query/class/tle_latest/ORDINAL/1/DECAYED/false/PERIOD/<128/format/json"

    with requests.Session() as session:
        session.post(url, data={"identity": username, "password": password})
        response = session.get(query)
        response.raise_for_status()
        return response.json()

if __name__ == "__main__":
    tles = fetch_tles()
    print(json.dumps(tles))

