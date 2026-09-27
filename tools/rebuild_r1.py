#!/usr/bin/env python3
from __future__ import annotations
import base64, hashlib, lzma, struct, zlib
from pathlib import Path

ORIGINAL_SIZE = 3124397
ORIGINAL_SHA256 = "81ce80f1bd5cc74183f621871e3ec4fa079bf0d694b652f7b20002cea82c1f21"
TARGET_SIZE = 3124397
TARGET_SHA256 = "9f6cf822e1927c4968dc22cc4328ea86cdd658cf486843498208be782c16dcfe"
TARGET_XZ_SHA256 = "4ea7321c93fe313199663085192eb66c70ad5bc03aef543a602872cded4b547e"
MAGIC = b"RTDP1\0"

def sha(b: bytes) -> str:
    return hashlib.sha256(b).hexdigest()

def main() -> int:
    root = Path(__file__).resolve().parents[1]
    parts = sorted((root / "payload").glob("R1_target.xz.b64.part*"))
    if not parts:
        raise SystemExit("R1 payload chunks are missing")
    encoded = "".join(p.read_text(encoding="ascii").strip() for p in parts)
    packed = base64.b64decode(encoded, validate=True)
    if sha(packed) != TARGET_XZ_SHA256:
        raise SystemExit("compressed target hash mismatch")
    target = lzma.decompress(packed)
    if len(target) != TARGET_SIZE or sha(target) != TARGET_SHA256:
        raise SystemExit("R1 target reconstruction failed")
    (root / "hedge.exe").write_bytes(target)

    raw = bytearray(MAGIC)
    raw += struct.pack("<QQ", ORIGINAL_SIZE, TARGET_SIZE)
    raw += bytes.fromhex(ORIGINAL_SHA256)
    raw += bytes.fromhex(TARGET_SHA256)
    raw += struct.pack("<I", 1)
    raw += b"\x01" + struct.pack("<I", len(target)) + target
    patch = zlib.compress(bytes(raw), 9)
    out = root / "patcher" / "over_the_hedge_r1.rtdp1.zlib"
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_bytes(patch)

    print("R1 hedge.exe SHA-256:", sha(target))
    print("RTDP1 payload SHA-256:", sha(patch))
    return 0

if __name__ == "__main__":
    raise SystemExit(main())
