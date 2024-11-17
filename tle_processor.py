import sys
import json
from sgp4.api import Satrec, jday
from datetime import datetime

def calculate_positions_from_file(file_path):
    with open(file_path, 'r') as f:
        tles = json.load(f)

    results = []
    for tle in tles:
        tle1, tle2 = tle["TLE_LINE1"], tle["TLE_LINE2"]
        satellite = Satrec.twoline2rv(tle1, tle2)

        # Using current UTC time to propagate satellite's position
        now = datetime.utcnow()
        jd, fr = jday(now.year, now.month, now.day, now.hour, now.minute, now.second + now.microsecond / 1e6)

        error, position, velocity = satellite.sgp4(jd, fr)
        if error == 0:
            results.append({"position": position, "velocity": velocity})
        else:
            results.append({"error": "SGP4 error"})
    return results

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 tle_processor.py <path_to_tle_file>")
        sys.exit(1)

    tle_file_path = sys.argv[1]
    positions = calculate_positions_from_file(tle_file_path)
    print(json.dumps(positions))
