import numpy as np
from matplotlib import axes
from matplotlib import figure
from matplotlib import pyplot as plt
from typing import List
from parse import PERCENTILES, Stats


class Plot:
    def __init__(self, figure: figure.Figure):
        self.figure = figure

    def settitle(self, title: str):
        self.figure.suptitle(title)

    def save(self, filename) -> None:
        self.figure.savefig(filename)
    
    def show(self) -> None:
        plt.show()


def draw(stats: List[Stats]) -> Plot:
    stats = sorted(stats, key=lambda stat: stat.label)
    fig, ax = plt.subplots(figsize=(len(stats), 8))
    width = 0.05
    xs = np.arange(len(PERCENTILES))
    colors = reversed(plt.colormaps["RdYlGn"](np.linspace(0, 1, len(stats))))
    for (index, (stat, color)) in enumerate(zip(stats, colors)):
        assert stat.percentiles == PERCENTILES
        ax.bar(xs + index * width, stat.latencies, width=width, label=stat.label, color=color)
    ax.set_xlabel("percentiles (%)")
    ax.set_xticks(xs + (len(stats) - 1) / 2 * width, PERCENTILES)
    ax.set_ylabel("response time (ns)")
    ax.legend()
    ax.set_title("Latency distributions")
    return Plot(fig)
