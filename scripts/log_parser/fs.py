from pathlib import Path
from typing import List


def read_dir(dir: Path) -> List[Path]:
    return [File(file) for file in dir.iterdir() if file.is_file() and file.suffix == ".json"]


def read_file(file: Path) -> List[str]:
    with open(file) as file:
        return file.readlines()


class File:
    def __init__(self, file: Path):
        self.native: Path = file
        self.lines: List[str] = read_file(self.native)


class Dir:
    def __init__(self, name: str):
        self.native: Path = Path(name)
        self.files: List[File] = read_dir(self.native)
