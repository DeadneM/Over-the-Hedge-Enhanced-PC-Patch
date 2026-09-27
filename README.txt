OVER THE HEDGE - WINDOWS 11 ENHANCED PATCH
============================================
Canonical build: V12A13 / public release R1
Status: VALIDATED

PACKAGE CONTENTS
----------------
This archive intentionally contains exactly two files:

  hedge.exe
  README.txt

No diagnostic logs, alternate executables, test utilities, video experiments,
or extra files are included in the canonical package.

IMPORTANT BASE RULE
-------------------
All modifications are based on the unpacked/deprotected retail hedge.exe.
The protected retail backup bckhedge.exe is NOT used as a patch base.

Original unpacked hedge.exe:
  Size: 3,124,397 bytes
  SHA-256: 81ce80f1bd5cc74183f621871e3ec4fa079bf0d694b652f7b20002cea82c1f21
  PDB reference: d:\Ports\Oth\bin\hedge.pdb

Protected retail executable supplied for reference:
  bckhedge.exe
  Size: 4,514,737 bytes
  Protected wrapper sections included stxt774 / stxt371.

Current canonical V12A13 / R1 hedge.exe:
  SHA-256: 9f6cf822e1927c4968dc22cc4328ea86cdd658cf486843498208be782c16dcfe

==========================================================================
1. PROJECT GOALS
==========================================================================

The patch modernizes the Windows PC version of Over the Hedge while keeping
native game behavior wherever possible.

Validated final goals include:

- Large Address Aware / 4 GB support.
- Native borderless/windowed presentation suitable for modern Windows.
- Desktop-aware resolution handling.
- Modern resolution choices, including 720p, 1080p, 2K, 4K and Desktop.
- Dynamic aspect-ratio handling with Hor+ behavior.
- Native HUD behavior preserved.
- 16x anisotropic filtering.
- Linear/trilinear mip filtering improvements.
- Conservative stage-0 positive MIP LOD-bias clamp.
- Native anti-aliasing choices preserved rather than forced.
- Windows 11 DPI awareness.
- Reliable Alt+F4 using the game's own quit path.
- Startup intro videos skipped without touching story cutscenes.
- Automatic profile activation only when exactly one profile slot is in use.

The project deliberately favors targeted, surgical binary changes over broad
engine hacks.

==========================================================================
2. VALIDATED GRAPHICS / WINDOWS BUILD HISTORY
==========================================================================

V1 - LARGE ADDRESS AWARE
------------------------
Status: VALIDATED

File:
  hedge_Win11_V1_4GB.exe

Change:
  PE Characteristics byte at file offset 0x136:
    0x010F -> 0x012F

Result:
  Enables Large Address Aware / 4 GB address-space support on 64-bit Windows.

SHA-256:
  e9d10b13e9bdab24c1ba9f01724371289657077b434fb2dec0b0d541b5dc25c8


V2A - NATIVE BORDERLESS
-----------------------
Status: VALIDATED

File:
  hedge_Win11_V2A_BORDERLESS_NATIVE_TEST.exe

Main changes:

- Window-style path around VA 0x57C1B4 keeps WS_POPUP behavior but removes
  WS_EX_TOPMOST.
- D3D present-parameter path around 0x5863D0 forces Windowed = TRUE and
  refresh rate = 0.
- Resolution getters at 0x5881A0 / 0x5881B0 use desktop dimensions.
- Exclusive display-mode switching is avoided.

Result:
  Modern borderless behavior with the game staying in a normal Windows/D3D9
  presentation path.

SHA-256:
  296d928c6ff522104d98350d2c3016cd68f45c151cb21449f7e37a7f90e01cef


V3 - DYNAMIC ASPECT / HOR+
--------------------------
Status: VALIDATED

Main change:
  VA 0x563143 uses the renderer width/height ratio dynamically.

Behavior:
  Vertical FOV remains native while horizontal FOV expands with wider aspect
  ratios, giving Hor+ behavior.

Important conclusion:
  The HUD already anchors correctly to the screen edges. The HUD is therefore
  frozen unless a concrete defect is found. No global HUD rescaling hack is
  part of the patch.

SHA-256:
  25c2ab29b375b15fec23017635c7dcc466971dd39e5935b3813ef789ff94eff9


V4 - ANISOTROPIC FILTERING 16x
------------------------------
Status: VALIDATED

Main changes:

Stage-0 material minimum filtering changed to anisotropic around:
  0x45A968
  0x45AF1E
  0x45EA28

Renderer initialization/reset applies anisotropic MAG/MIN filtering and
MAXANISOTROPY = 16 for texture stages 0..3.

Relevant caves/hooks:
  Cave 0x4ED095
  Cave 0x4EC095
  Hook 0x583628
  Hook 0x5AC9E9

SHA-256:
  3452e25551eba33ce2563a38e09be10eb9d96352ea0ee0d55c91c3ee61781342


V5 - FORCED MAXIMUM MSAA
------------------------
Status: REJECTED

The attempt to force maximum MSAA globally was rejected.

Permanent rule:
  Do NOT reintroduce forced anti-aliasing.
  The game's native Off / 2x / 4x choices must remain selectable.


V6A - MODERN RESOLUTION LIST
----------------------------
Status: VALIDATED as resolution-menu base

Visible order is intentionally:

  1. 720p    = 1280x720
  2. 1080p   = 1920x1080
  3. 2K      = 2560x1440
  4. 4K      = 3840x2160
  5. Desktop = current desktop resolution

Project terminology rule:
  "2K" means 2560x1440 in this project.

The Desktop entry is the internal default and stays last in the visible list.

Relevant locations:
  Labels around 0x6C500C..0x6C5038
  Desktop setter cave 0x4EBD00
  Desktop getters 0x4EBD40 / 0x4EBD60
  Default index patch around 0x58B5ED -> index 4

SHA-256:
  3338ea12a1383968bcd9cdcc6271855fd677531117377a0550fb1c9877fabe08


V6B - LOW-RES BORDERLESS ATTEMPT
--------------------------------
Status: REJECTED / INSUFFICIENT

The first attempt to reconcile a lower game render resolution with the
borderless desktop-sized window was insufficient.


V6C - DESKTOP-SIZED HWND + INDEPENDENT BACKBUFFER
-------------------------------------------------
Status: VALIDATED

Goal:
  The borderless window itself stays desktop-sized even when the selected
  internal render/backbuffer resolution is lower.

Changes:
  Initial SetWindowPos path 0x5865C5 -> cave 0x4EBE00
  Video-apply path 0x58A3F2..0x58A413 -> cave 0x4EBE40

Result:
  Desktop-sized HWND is decoupled from the render backbuffer resolution.

SHA-256:
  ec61ad02e6e7d4b9597833f989073bed4fbefe691e2d8d2b3114de6532b493f7


V7 - LINEAR / TRILINEAR MIP FILTERING
-------------------------------------
Status: VALIDATED

Patch:
  VA 0x58CDB5
  74 2D -> 90 90

Result:
  Enables the desired mip-filtering path instead of the original conditional
  branch behavior.

SHA-256:
  702ce1ff6ac2c1946cd02a6bd47fd6a341b6ad76d1ffad9fd24cd03328f7bea4


V8 - WINDOWS DPI AWARENESS
--------------------------
Status: VALIDATED

Entry path:
  0x65DD70 -> trampoline cave 0x4EC200

Behavior:
  Requests modern per-monitor DPI awareness where available, with fallback to
  SetProcessDPIAware.

SHA-256:
  dcbfe6ecb678462bebe7be9580f202faec1f8ab121615c190f72f186e523e808


V9 - CONSERVATIVE POSITIVE MIP LOD-BIAS CLAMP
---------------------------------------------
Status: VALIDATED

Central MIPMAPLODBIAS path:
  0x58CE36 -> cave 0x4EC400

Rule:
  Positive stage-0 MIP bias is clamped to 0.
  Zero and negative values are left unchanged.
  Stages 1+ are left unchanged.

This is intentionally conservative and not a global texture/LOD rewrite.

SHA-256:
  199247aa268121933311b3004762b1ac8c6de56b79b03c96f6a78f06bc882924


V10 - NATIVE ALT+F4
-------------------
Status: VALIDATED

WndProc:
  0x41C7E0

Logic:
  Detect WM_SYSKEYDOWN + VK_F4 + Alt context, then jump to the game's native
  WM_CLOSE handling at 0x41C81D.

The native close path sets:
  quit flag 0x734540 = 1
and calls PostQuitMessage.

Important rule:
  Preserve this native shutdown behavior in all future builds.

SHA-256:
