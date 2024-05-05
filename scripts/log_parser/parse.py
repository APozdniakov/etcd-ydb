import json
from typing import Any, Dict, List


PERCENTILES = [10, 25, 50, 75, 90, 95, 99, 99.9]


class Stats:
    def __init__(self, label: str, json: Dict[str, Any]):
        self.label: str = label
        self.total_time = json["TotalTime"]
        self.total = json["Total"]
        self.fastest = json["Fastest"]
        self.slowest = json["Slowest"]
        self.average = json["Average"]
        self.rps = json["RPS"]
        self.percentiles = list(map(lambda p: p["Percentile"], json["Percentiles"]))
        self.latencies = list(map(lambda p: p["Latency"], json["Percentiles"]))


def parse_stats(label: str, lines: List[str]) -> Stats:
    return Stats(label.replace("_", " "), json.loads("".join(lines)))
