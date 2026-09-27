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
  06ca4a6f0e9875f0f7bfd6181b54dd270694bb7128e45d2f42be766e09f01713


V11 - STARTUP INTRO SKIP
------------------------
Status: VALIDATED
Former canonical base before V12A11

File:
  hedge_Win11_V11_INTRO_SKIP_TEST.exe

EXE SHA-256:
  dcbf56965b8072afa84b368a06601db67867f3d4705d327a21078771f5339342

Original test ZIP SHA-256:
  9b6d5e4903c2fb44d57f55177b5339d09de5346f31a59086ddd61f87692e98ea

Implementation:
  Hook _BinkOpen@8 thunk at VA 0x683428
  Trampoline cave 0x4EC600

The hook performs a case-insensitive comparison of the last 8 characters of
movie paths and returns NULL only for:

  lo01.dat
  lo02.dat
  lo03.dat
  lo04.dat

Result:
  Startup intro videos are skipped.
  Story cutscenes remain untouched.

==========================================================================
3. VIDEO / ATTRACT-MODE EXPERIMENTS
==========================================================================

Status: ABANDONED BY DESIGN

A separate series of experiments tried to start Bink/attract-mode video from
UI input.

Important findings:

- The low-level QueueMovie path around 0x5911E0 could start Bink playback but
  the movie remained behind the menu.
- Forcing movieManager +0x08 = 2 caused the game to run accelerated and could
  make mouse input disappear.
- A movie-wrapper path around 0x42E400 was captured and tested.
- Multiple V12/V12C/V12D/V12E video candidates were rejected.

Permanent rule:
  The video/Tab feature is abandoned.
  Do NOT reintroduce it in future builds.

==========================================================================
4. AUTO-PROFILE RESEARCH: REQUIREMENT
==========================================================================

Final desired behavior:

- Inspect all four native profile slots.
- The physical/visible slot number must not matter.
- If EXACTLY ONE slot is in use, automatically activate that profile.
- If zero slots are used, remain fully vanilla/manual.
- If two or more slots are used, remain fully vanilla/manual.
- The auto path must use the game's own high-level native transition.
- No simulated mouse.
- No simulated keyboard.
- No coordinate-based input.
- No forced input-result branch.
- No direct SGLoadSlot shortcut.
- No raw final UI-event injection.

Visible/native mapping:

  visible slot 1 -> native index 0
  visible slot 2 -> native index 1
  visible slot 3 -> native index 2
  visible slot 4 -> native index 3

Important terminology:
  User-facing slot numbers are 1..4.
  Internal native indices are 0..3.

==========================================================================
5. AUTO-PROFILE: RELIABLE SLOT DETECTION
==========================================================================

The profile UI repeatedly presents four real EUIDDSaveSlot widgets.

Class information:

  EUIDDSaveSlot class string around 0x693238
  EUIDDSaveSlot vtable around 0x6930B8
  Forwarder at vtable +0x124 -> function 0x44C5E0
  Receiver around 0x44D140

Reliable SlotInUse observation point:
  0x44D50D

At that point:
  ESI = actual EUIDDSaveSlot widget
  EAX / then EBX = SlotInUse result
  [ESI + 0xD4] = native slot index 0..3

Reliable detection model:

  SEEN_MASK = bitmask of native indices observed
  USED_MASK = bitmask of indices whose SlotInUse != 0

A complete cycle requires:
  SEEN_MASK == 0x0F

Exactly one profile requires USED_MASK to contain exactly one set bit.

This detection method is retained in the final V12A11 implementation.

==========================================================================
6. AUTO-PROFILE: EARLY FAILED APPROACHES
==========================================================================

V12A / EARLY DIRECT LOAD
------------------------
Status: REJECTED

Early attempts scanned profile slots and called a save-slot loading path
around 0x42F260 directly.

Result:
  The save data could be touched, but this did not reproduce the complete UI
  transition to an active profile/main menu state.

Conclusion:
  Direct SGLoadSlot is not the correct high-level profile-selection action.


V12A2 / V12A3 - DEFERRED / UI-EVENT EXPERIMENTS
------------------------------------------------
Status: REJECTED

Various deferred/run/WndProc and UI-event approaches were explored based on
an early assumption about the selection event.

Critical correction:
  An earlier SLOT_ACTIVATE pointer was NOT a real EUIDDSaveSlot.
  Event 0x762D502E belonged to the wrong UI class and must not be used for
  auto-profile logic.

==========================================================================
7. V11F - PROFILE FLOW DIAGNOSTIC
==========================================================================

Status: DIAGNOSTIC

File:
  hedge_Win11_V11F_PROFILE_FLOW_DIAGNOSTIC.exe

Original diagnostic EXE SHA-256:
  82b4af1899617d0ce3d41549a26644664efdc237bfd7108f51e2a7a02916ed35

Original diagnostic ZIP SHA-256:
  6ef10725e57839c962838bf90c1bf43c42037648d10fc2898f5cfac98be486ac

The diagnostic logged:

  UNIQUE_SLOT
  SAVESLOT_FORWARD
  SAVESLOT_RECEIVE
  SGLOADSLOT

Important result:
  During manual confirmation of the real EUIDDSaveSlot, the game emitted:

    event = 0x8EE7556F

At the retail code around 0x44D016:

    push 0
    push esi              ; same EUIDDSaveSlot
    push 0x8EE7556F
    call [vtable+0x124]

This established 0x8EE7556F as the real confirmation event for the correct
EUIDDSaveSlot class.

The same diagnostic confirmed the user's active test profile was native
index 0 in that specific run. This was only an observation, never a rule.

==========================================================================
8. READY STATE / V11G DIAGNOSTIC
==========================================================================

The real EUIDDSaveSlot receives event:

  0x1579DAF8

Its native receiver sets:

  [EUIDDSaveSlot + 0x191] = 1

This is a genuine ready/enable state.

V11G profile-gate diagnostic:

File:
  hedge_Win11_V11G_PROFILE_GATE_DIAGNOSTIC.exe

EXE SHA-256:
  57a9e5ff65f9d9fd7ead95d0a537d435a33264568127d5cdefd014f40d2456ac

ZIP SHA-256:
  778b5eda8ee2187519cd5830eabf8ac77785de3b0509b56bbda1ab374b1477db

Logged:
  CYCLE
  READY_EVENT
  GATE_STATE
  CONFIRM_FORWARD

Result:

- All four native indices were repeatedly seen:
    CYCLE A=0x0F
- Exactly one used slot was present:
    B=0x01 in that test run
- READY_EVENT was observed.
- Packed GATE_STATE ended in 0101, showing both:
    [widget+0x191] == 1
    READY_EVENT observed
- Manual click still emitted the real confirmation path.

Conclusion:
  The ready gate was valid, but readiness alone did not make a direct
  confirmation-event injection equivalent to the real native flow.

==========================================================================
9. V12A5 - TRUE UI EVENT ATTEMPT
==========================================================================

Status: REJECTED

File:
  hedge_Win11_V12A5_AUTO_SINGLE_PROFILE_TRUE_UI_EVENT_TEST.exe

EXE SHA-256:
  db73112b1261e7284386f0ff59a3f571f28a53e9b592bd8e62b25b88ebcdae8a

ZIP SHA-256:
  bbbbc8266ab62047a6bd18384583654f67f8bd86abbc71d2a619a9b0740a2c9f

Method:

- Detect exactly one used profile.
- Wait for a valid parent [widget+0x8C].
- Wait for [widget+0x191] == 1.
- Call the EUIDDSaveSlot vtable+0x124 with:
    event  = 0x8EE7556F
    sender = widget
    arg3   = 0

Result:
  Nothing happened.

Conclusion:
  Injecting the final EUIDDSaveSlot confirmation event is not sufficient.
  The native path has important state/control flow before that event.

==========================================================================
10. NATIVE INPUT GATE DISCOVERY
==========================================================================

Static audit around the real manual confirm path near 0x44D00C found an
important input gate before 0x44D016.

Retail sequence:

- Read control pointer [ESI+0x168].
- Pass token [ESI+0x170].
- Call control virtual method +0x68.
- Test returned AL.
- If AL == 0, do not confirm.
- Only if AL != 0 does execution reach 0x44D016 and emit 0x8EE7556F.

This explained why calling the final event directly did not reproduce the
native behavior.

==========================================================================
11. V12A6 / V12A7 - FORCED INPUT RESULT
==========================================================================

V12A6
------
Status: REJECTED / CRASH

File:
  hedge_Win11_V12A6_AUTO_SINGLE_PROFILE_NATIVE_INPUT_TEST.exe

EXE SHA-256:
  fd00f2c30a6fb68683b08be66a6829d7e353c8be31a0a10e55f739d9e3dee284

ZIP SHA-256:
  c8e8fb8e0f363111feee3d4e234af7416b6b7e9471c20f656cd02e8dbf670f40

Method:
  Hook the native input test around 0x44D012 and force the true branch when
  the unique slot was ready.

Result:
  Crash.


V12A7
------
Status: REJECTED / CRASH

File:
  hedge_Win11_V12A7_AUTO_PROFILE_DELAYED_NATIVE_INPUT_TEST.exe

EXE SHA-256:
  265b5a2121b2e156c7b631219d3b7ab4414d1af5c7d6625c358362aeb30cd68e

ZIP SHA-256:
  0a6ca89ce9e73b69afdaffac60044a968ecf434b81eea97ca07067b04c307624

Method:
  Same native input-result override, but delayed by at least 1000 ms after the
  genuine ready event using GetTickCount.

Result:
  Crash.

Permanent rule:
  Never force or modify the AL result/branch at 0x44D012 again.

==========================================================================
12. V11H - INPUT DIAGNOSTIC
==========================================================================

Status: DIAGNOSTIC / CRASHED, BUT LOG WAS USEFUL

File:
  hedge_Win11_V11H_PROFILE_INPUT_DIAGNOSTIC.exe

EXE SHA-256:
  ff8c49877639f21b628efe0edc888a46c499899a4799ff6ec5640f86a522504d

ZIP SHA-256:
  c005cf1f5df11e569378938f2c3ed4ece1ff6aee8ac2d5d95b1792251db0b5c1

Log:
  hedge_profile_input.log

Useful lines captured before the diagnostic crashed:

  WINMSG A=00000201 B=00000001 C=02F307CE
  WINMSG A=00000202 B=00000000 C=02F307CE
  INPUT_TRUE A=029CC068 B=00000003 C=061EA160

Interpretation:

- The user's real manual confirmation was a normal left mouse click.
- WM_LBUTTONDOWN = 0x0201
- WM_LBUTTONUP   = 0x0202
- The game's native input control then saw token 3 as true.

The click in that 3840x2160 test occurred at x=1998, y=755.

Important design decision:
  Although this identified the physical input source, the project explicitly
  rejects mouse/keyboard emulation as an auto-profile solution.

==========================================================================
13. V12A8 - SYNTHETIC WIN32 CLICK
==========================================================================

Status: REJECTED BY DESIGN

File:
  hedge_Win11_V12A8_AUTO_PROFILE_POSTMESSAGE_CLICK_TEST.exe

EXE SHA-256:
  68779908829a302df84fbb1a96815288c42b1654154d6daeb18ebd52490e48b8

ZIP SHA-256:
  3160ee2e1bd3c7ffd3743c939e6bd18456bcdd8ef286e049d8df31d44007b0ae

Method:
  Attempted to reproduce the real manual click through Win32 messages after
  profile detection/readiness.

This approach was explicitly rejected.

Permanent rule:
  Do not emulate mouse clicks, keyboard presses, window messages or screen
  coordinates for auto-profile.

==========================================================================
14. V12A9 / V12A10 - WRONG NATIVE LOAD WRAPPER
==========================================================================

V12A9
------
Status: REJECTED / STABLE BUT NO EFFECT

File:
  hedge_Win11_V12A9_AUTO_SINGLE_PROFILE_NATIVE_LOAD_WRAPPER_TEST.exe

EXE SHA-256:
  9ecd910f83e37c88bd18440c7badfee038e8b74a2d3ec06306d6a8bc22cc02f7

ZIP SHA-256:
  802ee910cb68158daa921839d81f8c69bac37df421570262dfdd164eb08f417d

Method:
  Use wrapper 0x42F690(index), which internally loads:

    ECX = [0x6FA7E8]
    push index
    call 0x42F260

V12A9 did not crash, but did not perform the profile transition.


V12A10
-------
Status: REJECTED / NO EFFECT

File:
  hedge_Win11_V12A10_AUTO_SINGLE_PROFILE_NATIVE_FLOW_TEST.exe

EXE SHA-256:
  b4110283faa08c6c8e64fc6a6f62ddfd07a6aa02a107a81418e03dcb206811dc

ZIP SHA-256:
  10a61525646a84616fa0bdff033d8b5fa310702e192dd9c706d63a99027b27ab

Change from A9:
  Removed WndProc from auto-profile and tried the same wrapper from the game's
  own profile UI flow after additional stable slot cycles.

Result:
  Still no visible profile activation.

Conclusion:
  The problem was not timing. 0x42F690 was the wrong semantic operation.

==========================================================================
15. V11I - TRUE NATIVE LOAD PATH DIAGNOSTIC
==========================================================================

Status: DIAGNOSTIC

File:
  hedge_Win11_V11I_NATIVE_LOAD_PATH_DIAGNOSTIC.exe

EXE SHA-256:
  efda5235aad98d553faf349015835331fcc41cfe21b6d0608e095ed72e9c7012

ZIP SHA-256:
  0f3c631bf78d5b9d115f5c3aa3c412e97e3e7c32dae06ea54d194d6d0352ff79

Log:
  hedge_native_load.log

Hooks:

0x42F690 entry:
  LOAD_WRAPPER
  A = requested native slot index
  B = caller return address
  C = [0x6FA7E8]

0x42F260 entry:
  SGLOAD_ENTRY
  A = requested native slot index
  B = caller return address
  C = actual incoming save-manager ECX

0x42F6A0 return transform:
  WRAPPER_RETURN
  A = FFFFFFFF success / 00000000 failure

Manual-selection result:

The game called the wrapper for native indices:

  0, 1, 2, 3, then 0 again

Typical caller:
  0x47D3B8

The SGLOAD_ENTRY calls came from the 0x42F690 wrapper, with every call
reporting success.

Critical conclusion:
  0x42F690 / 0x42F260 enumerate/load slot data for UI/script state.
  They are NOT the high-level action that activates the selected profile.

This invalidated the semantic basis of V12A9 and V12A10.

==========================================================================
16. V11J - PROFILE PARENT ROUTE DIAGNOSTIC
==========================================================================

Status: DIAGNOSTIC

File:
  hedge_Win11_V11J_PROFILE_PARENT_ROUTE_DIAGNOSTIC.exe

EXE SHA-256:
  1b26ad39b5ac6c6167a6f4152fab6f4d36d736ba1c609f502936651b116bd7c8

ZIP SHA-256:
  62e1c2f3f4426b99daa7920af8aa6f17420a4b8b52144049319907847156e557

Log:
  hedge_profile_parent.log

Manual click produced:

  PARENT_ROUTE A=06164160 B=0616D430 C=0044E080
  PARENT_META  A=0616D430 B=006932E0 C=06164160

Meaning:

  child EUIDDSaveSlot = 0x06164160 in that run
  parent UI object    = 0x0616D430
  parent vtable       = 0x006932E0
  parent handler      = 0x0044E080

Static audit of 0x44E080 showed that it is itself only another forwarder.

Behavior of 0x44E080:

- If [this+0x8C] exists, forward the same event to:
    parent = [this+0x8C]
    handler = [parent.vtable+0x124]
- Otherwise, forward to the global fallback:
    object = [0x702828]
    handler = [global.vtable+0x08]

Conclusion:
  V11J identified an intermediate UI layer, not the final profile action.

==========================================================================
17. V11K - PROFILE PARENT CHAIN DIAGNOSTIC
==========================================================================

Status: DIAGNOSTIC

File:
  hedge_Win11_V11K_PROFILE_PARENT_CHAIN_DIAGNOSTIC.exe

EXE SHA-256:
