#!/usr/bin/env python3
"""Control an Ikemen GO X11 window and capture named screenshots.

This helper is intentionally dependency-light for local visual-probe passes:
- window discovery/focus uses wmctrl
- screenshots use ImageMagick import
- key events use libX11 + libXtst via ctypes, so xdotool is not required
"""

from __future__ import annotations

import argparse
import ctypes
import os
import re
import shlex
import signal
import subprocess
import sys
import time
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable


DEFAULT_STEPS = [
    "wait:2",
    "snap:boot",
    "key:Return",
    "wait:1",
    "snap:after-enter",
]

KEY_ALIASES = {
    "enter": "Return",
    "return": "Return",
    "esc": "Escape",
    "escape": "Escape",
    "space": "space",
    "up": "Up",
    "down": "Down",
    "left": "Left",
    "right": "Right",
    "tab": "Tab",
    "backspace": "BackSpace",
    "ctrl": "Control_L",
    "control": "Control_L",
    "alt": "Alt_L",
    "shift": "Shift_L",
}


@dataclass
class Window:
    wid: str
    desktop: str
    pid: str
    x: int
    y: int
    w: int
    h: int
    host: str
    title: str


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Launch/attach to Ikemen GO, send X11 key input, and save window screenshots.",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument("--bin", default="./Ikemen_GO_Linux", help="Ikemen executable to launch")
    parser.add_argument("--workdir", default=".", help="Runtime working directory for launch mode")
    parser.add_argument("--no-launch", action="store_true", help="Attach to an existing Ikemen window instead of launching")
    parser.add_argument("--window-title", default=r"(?i)ikemen", help="Regex used to find the target window title")
    parser.add_argument("--output-dir", default="artifacts/visual-probes", help="Directory for screenshots")
    parser.add_argument("--probe-mode", default=os.environ.get("IKEMEN_RENDER_PROBES", "edge"), help="IKEMEN_RENDER_PROBES value for launch mode; empty disables")
    parser.add_argument("--timeout", type=float, default=15.0, help="Seconds to wait for a matching window")
    parser.add_argument("--settle", type=float, default=0.15, help="Delay after focus/key actions")
    parser.add_argument(
        "--step",
        action="append",
        default=[],
        help=(
            "Workflow step. Repeatable. Supported forms: wait:SECONDS, key:KEY, "
            "hold:KEY:SECONDS, keys:KEY+KEY, text:TEXT, snap:LABEL, focus. "
            "If omitted, captures boot and after-enter."
        ),
    )
    parser.add_argument("--dry-run", action="store_true", help="Print planned actions without launching, keying, or capturing")
    parser.add_argument("--keep-running", action="store_true", help="Do not terminate the launched process at the end")
    return parser.parse_args()


def run(cmd: list[str], *, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(cmd, check=check, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)


def require_tool(name: str) -> None:
    if subprocess.run(["/usr/bin/env", "bash", "-lc", f"command -v {shlex.quote(name)} >/dev/null"], check=False).returncode != 0:
        raise RuntimeError(f"Required tool not found on PATH: {name}")


def require_icon_assets(workdir: Path) -> None:
    missing = [
        rel
        for rel in (
            Path("external/icons/IkemenCylia_256.png"),
            Path("external/icons/IkemenCylia_96.png"),
            Path("external/icons/IkemenCylia_48.png"),
        )
        if not (workdir / rel).is_file()
    ]
    if missing:
        rendered = "\n".join(f"  - {path}" for path in missing)
        raise RuntimeError(
            f"Missing required Ikemen icon assets in {workdir}:\n{rendered}\n"
            "The default config expects the Cylia icon set under external/icons/."
        )


def list_windows() -> list[Window]:
    proc = run(["wmctrl", "-lpG"], check=False)
    windows: list[Window] = []
    for line in proc.stdout.splitlines():
        parts = line.split(None, 8)
        if len(parts) < 9:
            continue
        wid, desktop, pid, xs, ys, ws, hs, host, title = parts
        try:
            windows.append(Window(wid, desktop, pid, int(xs), int(ys), int(ws), int(hs), host, title))
        except ValueError:
            continue
    return windows


def find_window(title_regex: str, *, pid: int | None, timeout: float) -> Window:
    pattern = re.compile(title_regex)
    deadline = time.time() + timeout
    last: list[Window] = []
    while time.time() <= deadline:
        windows = list_windows()
        last = windows
        matches = [w for w in windows if pattern.search(w.title)]
        if pid is not None:
            pid_matches = [w for w in matches if w.pid == str(pid)]
            if pid_matches:
                return pid_matches[-1]
        if matches:
            return matches[-1]
        time.sleep(0.25)
    titles = "; ".join(f"{w.wid} pid={w.pid} title={w.title!r}" for w in last[-10:]) or "no windows listed"
    raise RuntimeError(f"No window matched {title_regex!r} within {timeout}s. Last windows: {titles}")


def focus_window(window: Window, settle: float) -> None:
    run(["wmctrl", "-ia", window.wid])
    time.sleep(settle)


class XKeySender:
    def __init__(self) -> None:
        self.x11 = ctypes.cdll.LoadLibrary("libX11.so.6")
        self.xtst = ctypes.cdll.LoadLibrary("libXtst.so.6")
        self.x11.XOpenDisplay.argtypes = [ctypes.c_char_p]
        self.x11.XOpenDisplay.restype = ctypes.c_void_p
        self.x11.XStringToKeysym.argtypes = [ctypes.c_char_p]
        self.x11.XStringToKeysym.restype = ctypes.c_ulong
        self.x11.XKeysymToKeycode.argtypes = [ctypes.c_void_p, ctypes.c_ulong]
        self.x11.XKeysymToKeycode.restype = ctypes.c_uint
        self.x11.XFlush.argtypes = [ctypes.c_void_p]
        self.x11.XCloseDisplay.argtypes = [ctypes.c_void_p]
        self.xtst.XTestFakeKeyEvent.argtypes = [ctypes.c_void_p, ctypes.c_uint, ctypes.c_int, ctypes.c_ulong]
        self.display = self.x11.XOpenDisplay(os.environ.get("DISPLAY", ":0").encode())
        if not self.display:
            raise RuntimeError("Could not open X display for key injection")

    def close(self) -> None:
        if self.display:
            self.x11.XCloseDisplay(self.display)
            self.display = None

    def keycode(self, key_name: str) -> int:
        normalized = KEY_ALIASES.get(key_name.lower(), key_name)
        keysym = self.x11.XStringToKeysym(normalized.encode())
        if keysym == 0 and len(normalized) == 1:
            keysym = self.x11.XStringToKeysym(normalized.lower().encode())
        if keysym == 0:
            raise RuntimeError(f"Unknown X11 key name: {key_name!r}")
        code = self.x11.XKeysymToKeycode(self.display, keysym)
        if code == 0:
            raise RuntimeError(f"No keycode for X11 key name: {key_name!r}")
        return int(code)

    def tap(self, key_name: str, settle: float) -> None:
        self.hold(key_name, 0.05, settle)

    def hold(self, key_name: str, duration: float, settle: float) -> None:
        code = self.keycode(key_name)
        self.xtst.XTestFakeKeyEvent(self.display, code, 1, 0)
        self.x11.XFlush(self.display)
        time.sleep(duration)
        self.xtst.XTestFakeKeyEvent(self.display, code, 0, 0)
        self.x11.XFlush(self.display)
        time.sleep(settle)

    def combo(self, key_names: Iterable[str], settle: float) -> None:
        codes = [self.keycode(k) for k in key_names]
        for code in codes:
            self.xtst.XTestFakeKeyEvent(self.display, code, 1, 0)
        for code in reversed(codes):
            self.xtst.XTestFakeKeyEvent(self.display, code, 0, 0)
        self.x11.XFlush(self.display)
        time.sleep(settle)

    def text(self, value: str, settle: float) -> None:
        for ch in value:
            if ch == " ":
                self.tap("space", settle)
            elif ch == "\n":
                self.tap("Return", settle)
            else:
                self.tap(ch, settle)


def capture(window: Window, output_dir: Path, label: str) -> Path:
    safe = re.sub(r"[^A-Za-z0-9_.-]+", "-", label).strip("-") or "snapshot"
    stamp = time.strftime("%Y%m%d-%H%M%S")
    path = output_dir / f"{stamp}-{safe}.png"
    output_dir.mkdir(parents=True, exist_ok=True)
    # ImageMagick import accepts the wmctrl hex window id and captures just the game window.
    run(["import", "-window", window.wid, str(path)])
    return path


def launch(args: argparse.Namespace) -> subprocess.Popen[str] | None:
    if args.no_launch:
        return None
    env = os.environ.copy()
    if args.probe_mode:
        env["IKEMEN_RENDER_PROBES"] = args.probe_mode
    workdir = Path(args.workdir).resolve()
    binary = Path(args.bin)
    if not binary.is_absolute():
        binary = (workdir / binary).resolve()
    if not binary.exists():
        raise RuntimeError(f"Ikemen binary not found: {binary}")
    if not os.access(binary, os.X_OK):
        raise RuntimeError(f"Ikemen binary is not executable: {binary}")
    require_icon_assets(workdir)
    return subprocess.Popen([str(binary)], cwd=str(workdir), env=env, text=True)


def main() -> int:
    args = parse_args()
    steps = args.step or DEFAULT_STEPS
    repo_root = Path(args.workdir).resolve()
    output_dir = Path(args.output_dir)
    if not output_dir.is_absolute():
        output_dir = repo_root / output_dir

    print("Ikemen control/snapshot workflow")
    print(f"  workdir: {repo_root}")
    print(f"  launch: {'no, attach existing window' if args.no_launch else args.bin}")
    print(f"  probe mode: {args.probe_mode!r}")
    print(f"  window title regex: {args.window_title}")
    print(f"  output dir: {output_dir}")
    print("  steps:")
    for step in steps:
        print(f"    - {step}")
    if args.dry_run:
        return 0

    require_tool("wmctrl")
    require_tool("import")

    proc: subprocess.Popen[str] | None = None
    sender: XKeySender | None = None
    screenshots: list[Path] = []
    try:
        proc = launch(args)
        window = find_window(args.window_title, pid=proc.pid if proc else None, timeout=args.timeout)
        print(f"Matched window: {window.wid} pid={window.pid} title={window.title!r} geom={window.w}x{window.h}+{window.x}+{window.y}")
        focus_window(window, args.settle)
        sender = XKeySender()
        for raw in steps:
            if ":" in raw:
                command, value = raw.split(":", 1)
            else:
                command, value = raw, ""
            command = command.strip().lower()
            value = value.strip()
            if command == "wait":
                time.sleep(float(value))
            elif command == "focus":
                focus_window(window, args.settle)
            elif command == "key":
                focus_window(window, args.settle)
                sender.tap(value, args.settle)
            elif command == "hold":
                focus_window(window, args.settle)
                if ":" not in value:
                    raise RuntimeError(f"hold step requires hold:KEY:SECONDS, got {raw!r}")
                key_name, duration = value.rsplit(":", 1)
                sender.hold(key_name.strip(), float(duration), args.settle)
            elif command == "keys":
                focus_window(window, args.settle)
                sender.combo([part.strip() for part in value.split("+") if part.strip()], args.settle)
            elif command == "text":
                focus_window(window, args.settle)
                sender.text(value, args.settle)
            elif command == "snap":
                path = capture(window, output_dir, value)
                screenshots.append(path)
                print(f"snapshot: {path}")
            else:
                raise RuntimeError(f"Unknown step {raw!r}")
        if screenshots:
            print("Captured screenshots:")
            for path in screenshots:
                print(f"  {path}")
        return 0
    finally:
        if sender is not None:
            sender.close()
        if proc is not None and not args.keep_running:
            proc.terminate()
            try:
                proc.wait(timeout=3)
            except subprocess.TimeoutExpired:
                proc.send_signal(signal.SIGKILL)


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except KeyboardInterrupt:
        raise SystemExit(130)
    except Exception as exc:
        print(f"ERROR: {exc}", file=sys.stderr)
        raise SystemExit(1)
