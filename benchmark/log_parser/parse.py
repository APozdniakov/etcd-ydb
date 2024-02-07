import re
from typing import List


SUMMARY = ["Total", "Slowest", "Fastest", "Average", "Stddev"]
QUANTILES = [10, 25, 50, 75, 90, 95, 99, 99.9]


class Stats:
    def __init__(self, label: str):
        self.label: str = label
        self.xs: List[float] = []
        self.ys: List[int] = []
        self.latencies: List[float] = []

    def __getitem__(self, key):
        return self.__dict__[key]
    
    def __setitem__(self, key, value):
        self.__dict__[key] = value


def noop_cb(stat: Stats, match: re.Match) -> None:
    pass


def stat_cb(stat: Stats, match: re.Match) -> None:
    stat[match.group(1)] = float(match.group(2))


def rps_cb(stat: Stats, match: re.Match) -> None:
    stat["Rps"] = float(match.group(2))


def bar_cb(stat: Stats, match: re.Match) -> None:
    stat.xs.append(float(match.group(1)))
    stat.ys.append(int(match.group(2)))


def latency_cb(stat: Stats, match: re.Match) -> None:
    stat.latencies.append(float(match.group(2)))


INT = "\d+"
FLOAT = f"{INT}\.{INT}"
FILE_PATTERN = [
    (rf"^\n$", noop_cb),
    (rf"^Summary:\n$", noop_cb),
    *[(rf"^  ({summary}):\t({FLOAT}) secs.\n$", stat_cb) for summary in SUMMARY],
    (rf"^  (Requests/sec):\t({FLOAT})\n$", rps_cb),
    (rf"^\n$", noop_cb),
    (rf"^Response time histogram:\n$", noop_cb),
    *[(rf"^  ({FLOAT}) \[({INT})\]\t\|∎*\n$", bar_cb) for _ in range(11)],
    (rf"^\n$", noop_cb),
    (rf"^Latency distribution:\n$", noop_cb),
    *[(rf"^  ({quantile})% in ({FLOAT}) secs.\n$", latency_cb) for quantile in QUANTILES],
    (rf"^\n$", noop_cb),
]

def parse_label(label: str) -> str:
    match label:
        case "1.txt":
            return "0-4   GB"
        case "2.txt":
            return "4-8   GB"
        case "3.txt":
            return "8-8.4 GB"
        case _:
            raise RuntimeError("Unknown label")


def parse_stat(label: str, lines: List[str]) -> Stats:
    result = Stats(parse_label(label))
    for line, (pattern, callback) in zip(lines, FILE_PATTERN):
        if (match := re.match(pattern, line)) is not None:
            callback(result, match)
        else:
            raise RuntimeError(f"Unknown line: \"{line}\" (expected pattern \"{pattern}\")")
    return result


def parse_stats(label: str, lines: List[str]) -> List[Stats]:
    return [parse_stat(label, lines[i: i + len(FILE_PATTERN)]) for i in range(0, len(lines), len(FILE_PATTERN))]
