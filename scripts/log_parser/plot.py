import numpy as np
from matplotlib import figure
from matplotlib import pyplot as plt
import pandas as pd
from typing import Dict, List
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


def draw(dirs: Dict[str, List[Stats]]) -> Plot:
    df = pd.DataFrame([{"label": label, "name": stat.label, **stat.percentiles} for label, stats in dirs.items() for stat in stats], columns=["label", "name", *PERCENTILES])
    fig, ax = plt.subplots(figsize=(2 * len(PERCENTILES),6))
    colors = reversed(plt.colormaps["RdYlBu"](np.linspace(0, 1, len(dirs))))
    legend = {}
    for label, color in zip(dirs.keys(), colors):
        for p_index, percentile in enumerate(PERCENTILES):
            percentiles = df[df["label"] == label].sort_values("name")[percentile].to_numpy()
            plot, = ax.plot(p_index + (np.arange(len(percentiles)) - (len(percentiles) - 1) / 2) * 0.05, percentiles, color=color)
        legend[label] = plot
    ax.set_xlabel("Перцентили (%)")
    ax.set_xticks(np.arange(len(PERCENTILES)), PERCENTILES)
    ax.set_ylabel("время ответа (нс)")
    ax.set_yscale("log")
    ax.legend(legend.values(), legend.keys())
    ax.grid()
    ax.set_title("Перцентили задержек")
    return Plot(fig)
