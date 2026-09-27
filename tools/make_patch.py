#!/usr/bin/env python3
"""Generate the compact source-locked RTDP1 delta used by the R1 patcher.

Usage:
    python tools/make_patch.py ORIGINAL_UNPACKED.exe R1.exe OUTPUT.b64

Both files must have the same size. The generated payload contains COPY
instructions for unchanged regions and LITERAL instructions only for changed
regions. Source and target SHA-256 hashes are embedded in the patch header.
"""
from __future__ import annotations

import base64
import hashlib
import struct
import sys
import zlib
from pathlib import Path

MAGIC = b"RTDP1\0"
MAX_LITERAL_GAP = 16


def sha(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def main() -> int:
    if len(sys.argv) != 4:
        print("Usage: make_patch.py ORIGINAL_UNPACKED.exe R1.exe OUTPUT.b64")
        return 2

    source_path, target_path, output_path = map(Path, sys.argv[1:])
    source = source_path.read_bytes()
    target = target_path.read_bytes()

    if len(source) != len(target):
        raise SystemExit("R1 generator expects equal-size source and target")

    runs: list[list[int]] = []
    i = 0
    while i < len(source):
        if source[i] == target[i]:
            i += 1
            continue
        j = i + 1
        while j < len(source) and source[j] != target[j]:
            j += 1
        runs.append([i, j])
        i = j

    groups: list[list[int]] = []
    for start, end in runs:
        if not groups or start - groups[-1][1] > MAX_LITERAL_GAP:
            groups.append([start, end])
        else:
            groups[-1][1] = end

    instructions: list[tuple[int, int, bytes]] = []
    pos = 0
    for start, end in groups:
        if start > pos:
            instructions.append((0, pos, source[pos:start]))
        instructions.append((1, 0, target[start:end]))
        pos = end
    if pos < len(source):
        instructions.append((0, pos, source[pos:]))

    raw = bytearray(MAGIC)
    raw += struct.pack("<QQ", len(source), len(target))
    raw += hashlib.sha256(source).digest()
    raw += hashlib.sha256(target).digest()
    raw += struct.pack("<I", len(instructions))

    for opcode, offset, data in instructions:
        raw += bytes([opcode])
        if opcode == 0:
            raw += struct.pack("<QI", offset, len(data))
        else:
            raw += struct.pack("<I", len(data))
            raw += data

    packed = zlib.compress(bytes(raw), 9)
    output_path.write_text(base64.b64encode(packed).decode("ascii") + "\n", encoding="ascii")

    print("source sha256:", sha(source))
    print("target sha256:", sha(target))
    print("changed runs:", len(runs))
    print("literal groups:", len(groups))
    print("literal bytes:", sum(end - start for start, end in groups))
    print("compressed patch bytes:", len(packed))
    print("patch sha256:", sha(packed))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
