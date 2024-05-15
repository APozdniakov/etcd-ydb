import json
from typing import Any, Dict, List


PERCENTILES = [10, 25, 50, 75, 90, 95, 99, 99.9]


class Stats:
    def __init__(self, label: str, json: Dict[str, Any]):
        self.label = label
        self.total_time = json["TotalTime"]
        self.total = json["Total"]
        self.fastest = json["Fastest"]
        self.slowest = json["Slowest"]
        self.average = json["Average"]
        self.rps = json["RPS"]
        self.percentile = list(map(lambda p: p["Percentile"], json["Percentiles"]))
        self.latency = list(map(lambda p: p["Latency"], json["Percentiles"]))
        self.percentiles = {p["Percentile"]: p["Latency"] for p in json["Percentiles"]}

def get_label(label: str) -> str:
    return f"{int(label):02d} ГБ"


def parse_stats(label: str, lines: List[str]) -> Stats:
    return Stats(get_label(label), json.loads("".join(lines)))
