import numpy as np
from matplotlib import axes
from matplotlib import figure
from matplotlib import pyplot as plt
from typing import List
from parse import QUANTILES, Stats


class Plot:
    def __init__(self, figure: figure.Figure):
        self.figure = figure

    def settitle(self, title: str):
        self.figure.suptitle(title)

    def save(self, filename) -> None:
        self.figure.savefig(filename)
    
    def show(self) -> None:
        plt.show()


def transpose(stats_list: List[List[Stats]]) -> List[List[Stats]]:
    length = max(map(lambda stats: len(stats), stats_list))
    return [[stats[i] for stats in stats_list if i < len(stats)] for i in range(length)]


def draw_response_time_histograms(ax: axes.Axes, stats: List[Stats]) -> None:
    for stat in stats:
        xs = stat.xs
        ys = stat.ys / np.sum(stat.ys) * 100
        ax.plot(xs, ys, label=stat.label)
    ax.set_xlabel("response time (sec)")
    ax.set_ylabel("operations (%)")
    ax.legend()
    ax.set_title("Response time plots")


def draw_latency_distribution_bars(ax: axes.Axes, stats: List[Stats]) -> None:
    width = 0.25
    xs = np.arange(len(QUANTILES))
    colors = reversed(plt.colormaps["RdYlGn"](np.linspace(0, 1, len(stats))))
    for (index, (stat, color)) in enumerate(zip(stats, colors)):
        ax.bar(xs + index * width, stat.latencies, width=width, label=stat.label, color=color)
    ax.set_xlabel("percentiles (%)")
    ax.set_xticks(xs + (len(stats) - 1) / 2 * width, QUANTILES)
    ax.set_ylabel("response time (sec)")
    ax.legend()
    ax.set_title("Latency distributions")


def draw_stats(subfig: figure.Figure, stats: List[Stats]) -> None:
    ax0, ax1 = subfig.subplots(nrows=1, ncols=2)
    draw_response_time_histograms(ax0, stats)
    draw_latency_distribution_bars(ax1, stats)


def draw(stats_list: List[List[Stats]]) -> Plot:
    stats_list = transpose(stats_list)
    fig = plt.figure(figsize=(10, len(stats_list) * 4.8))
    subfigs = fig.subfigures(nrows=len(stats_list), ncols=1)
    subfigs = np.reshape(np.array([subfigs]), (len(stats_list)))
    for subfig, stats in zip(subfigs, stats_list):
        draw_stats(subfig, stats)
    return Plot(fig)
