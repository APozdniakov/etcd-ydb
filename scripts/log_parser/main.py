import argparse
from fs import Dir
from parse import parse_stats
from plot import draw


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--save", action="store_true", help="save plot to file")
    parser.add_argument("dirs", metavar="dir", type=str, nargs="+", help="dir to process")
    return parser.parse_args()


def main(args: argparse.Namespace) -> None:
    for dir in list(map(Dir, args.dirs)):
        stats = [parse_stats(file.native.stem, file.lines) for file in dir.files]
        plot = draw(stats)
        plot.settitle(f"Benchmark results for different fill rate for {str(dir.native).replace('/', ' ')}")
        if args.save:
            plot.save(dir.native / f"{str(dir.native).replace('/', '_')}_plot.png")
        else:
            plot.show()


if __name__ == "__main__":
    main(parse_args())
