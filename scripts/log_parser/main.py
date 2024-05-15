import argparse
from fs import Dir
from parse import parse_stats
from plot import draw


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--save", type=str, help="file in which save the plot")
    parser.add_argument("dirs", metavar="dir", type=str, nargs="+", help="dir to process")
    return parser.parse_args()


def main(args: argparse.Namespace) -> None:
    stats = {str(dir.native).replace("/", " "): [parse_stats(file.native.stem, file.lines) for file in dir.files] for dir in list(map(Dir, args.dirs))}
    plot = draw(stats)
    # plot.settitle("Результаты бенчмарков")
    if args.save is not None:
        plot.save(args.save)
    else:
        plot.show()


if __name__ == "__main__":
    main(parse_args())
