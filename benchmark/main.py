import os
import sys
import re
from matplotlib import pyplot as plt
from typing import List, Tuple


SUMMARY_STATS = ["Total", "Slowest", "Fastest", "Average", "Stddev"]
QUANTILES = [10, 25, 50, 75, 90, 95, 99, 99.9]
INT = "\d+"
FLOAT = f"{INT}\.{INT}"
STAT_PATTERN = rf"^  ({'|'.join(SUMMARY_STATS)}):\t({FLOAT}) secs.\n$"
RPS_PATTERN = rf"^  (Requests/sec):\t({FLOAT})\n$"
BAR_PATTERN = rf"^  ({FLOAT}) \[({INT})\]\t\|∎*\n$"
LATENCY_PATTERN = rf"^  ({'|'.join(map(str, QUANTILES))})% in ({FLOAT}) secs.\n$"
FILE_PATTERN = [
    rf"^\n$",
    rf"^Summary:\n$",
    *[rf"^  {stat}:\t({FLOAT}) secs.\n$" for stat in SUMMARY_STATS],
    RPS_PATTERN,
    rf"^\n$",
    rf"^Response time histogram:\n$",
    *[BAR_PATTERN for _ in range(11)],
    rf"^\n$",
    rf"^Latency distribution:\n$",
    *[rf"^  {quantile}% in ({FLOAT}) secs.\n$" for quantile in QUANTILES],
    rf"^\n$",
]

class Stats:
    def __init__(self, label: str):
        self.label: str = label
        self.xy: List[Tuple[float, int]] = []
        self.latencies: List[Tuple[float, float]] = []

    def __getitem__(self, key):
        return self.__dict__[key]
    
    def __setitem__(self, key, value):
        self.__dict__[key] = value


def read_dir(dir: str) -> List[str]:
    return list(os.scandir(dir))


def read_file(filename: str) -> List[str]:
    with open(filename) as file:
        return file.readlines()


def parse_lines(label: str, lines: List[str]) -> Stats:
    result = Stats(label)
    for line, pattern in zip(lines, FILE_PATTERN):
        if re.match(pattern, line) is not None:
            if (match := re.match(STAT_PATTERN, line)) is not None:
                result[match.group(1)] = float(match.group(2))
            elif (match := re.match(RPS_PATTERN, line)) is not None:
                result["Rps"] = float(match.group(2))
            elif (match := re.match(BAR_PATTERN, line)) is not None:
                result.xy.append((float(match.group(1)), int(match.group(2))))
            elif (match := re.match(LATENCY_PATTERN, line)) is not None:
                result.latencies.append((match.group(1), float(match.group(2))))
        else:
            raise RuntimeError()
    return result


def unzip(l: List[Tuple]) -> Tuple[List, List]:
    return list(map(lambda t: t[0], l)), list(map(lambda t: t[1], l))


def draw_plot(stats: List[Stats]) -> None:
    _, ax = plt.subplots(layout='constrained')
    for stat in stats:
        xs, ys = unzip(stat.xy)
        quantiles, latencies = unzip(stat.latencies)
        # plt.plot(xs, ys, label=stat.label)
        # plt.bar(xs, ys, label=stat.label, alpha=0.5)
        plt.bar(quantiles, latencies, label=stat.label, alpha=0.5)
    # plt.yscale("log")
    ax.legend()
    plt.show()


def main():
    for dir in sys.argv[1:]:
        stats = [parse_lines(file.path, read_file(file)) for file in read_dir(dir)]
        draw_plot(stats)


if __name__ == '__main__':
    main()
