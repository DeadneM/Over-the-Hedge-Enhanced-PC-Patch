# Over the Hedge - Enhanced PC Patch

Modernization patch for the original PC release of **Over the Hedge**.

## R1

R1 is the first public cumulative release, based on the validated V12A13 executable.

### Main improvements

- Large Address Aware / 4 GB executable support
- Native borderless rendering on modern Windows
- Modern resolution list: 720p, 1080p, 2K, 4K, Desktop
- Desktop-size borderless window with independent render resolution
- Dynamic aspect ratio and Hor+ FOV
- 16x anisotropic filtering and trilinear mip filtering
- Conservative stage-0 positive MIP LOD-bias clamp
- Native FSAA Off / 2x / 4x preserved
- Windows 11 DPI awareness
- Native orderly Alt+F4
- Startup intro skip while preserving story cinematics
- Automatic selection of the only existing profile through the game's native high-level profile flow
- Windows key restored by removing only DirectInput `DISCL_NOWINKEY`

The full technical notebook is in `README.txt`.

## Installation

### Patcher

1. Use the unpacked/deprotected original PC `hedge.exe` supported by this project.
2. Put `OverTheHedge_Patcher_R1.exe` next to `hedge.exe`.
3. Run the patcher.
4. A verified backup is created before installation.

The patcher is fail-closed. Unknown executables are refused and never modified.

Supported original SHA-256:
`81ce80f1bd5cc74183f621871e3ec4fa079bf0d694b652f7b20002cea82c1f21`

R1 patched SHA-256:
`9f6cf822e1927c4968dc22cc4328ea86cdd658cf486843498208be782c16dcfe`

## Notes

- The protected/encrypted `bckhedge.exe` is **not** a supported patch source.
- Native 60 Hz timing is deliberately preserved after runtime audit.
- Video/TAB experiments were abandoned and are not part of R1.
