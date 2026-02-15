#!/usr/bin/env python3

"""Run development stack in separate Docker containers without docker-compose."""

from __future__ import annotations

import argparse
import errno
import shutil
import socket
import subprocess
import sys
import time
from dataclasses import dataclass
from pathlib import Path


ROOT_DIR = Path(__file__).resolve().parents[1]

NETWORK_NAME = "task_tracker_standalone_net"
POSTGRES_CONTAINER = "task_tracker_standalone_postgres"
SERVER_CONTAINER = "task_tracker_standalone_server"
WEB_CONTAINER = "task_tracker_standalone_web"
POSTGRES_VOLUME = "task_tracker_standalone_postgres_data"

SERVER_IMAGE = "task-tracker-dev-server:local"
WEB_IMAGE = "task-tracker-dev-web:local"
POSTGRES_IMAGE = "postgres:16-alpine"


@dataclass
class Ports:
    postgres: int
    server: int
    web: int


def run(cmd: list[str], *, check: bool = True, capture: bool = False) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        cmd,
        cwd=ROOT_DIR,
        check=check,
        text=True,
        capture_output=capture,
    )


def require_docker() -> None:
    if shutil.which("docker") is None:
        print("docker command not found", file=sys.stderr)
        sys.exit(1)


def port_in_use(port: int) -> bool:
    sockets: list[socket.socket] = []
    try:
        # Docker publishes on all interfaces; verify both IPv4 and IPv6 bindings.
        sock4 = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock4.bind(("0.0.0.0", port))
        sockets.append(sock4)

        try:
            sock6 = socket.socket(socket.AF_INET6, socket.SOCK_STREAM)
            sock6.setsockopt(socket.IPPROTO_IPV6, socket.IPV6_V6ONLY, 1)
            sock6.bind(("::", port))
            sockets.append(sock6)
        except OSError as exc:
            # IPv6 might be unavailable; ignore that case only.
            if exc.errno not in {errno.EAFNOSUPPORT, errno.EPROTONOSUPPORT, errno.EADDRNOTAVAIL}:
                raise
    except OSError:
        return True
    finally:
        for sock in sockets:
            sock.close()

    return False


def pick_port(preferred: int, *, allow_auto: bool, upper_bound: int = 65535) -> int:
    if not port_in_use(preferred):
        return preferred

    if not allow_auto:
        raise RuntimeError(f"Port {preferred} is already in use.")

    for candidate in range(preferred + 1, upper_bound + 1):
        if not port_in_use(candidate):
            print(f"Port {preferred} is busy, using {candidate} instead.")
            return candidate

    raise RuntimeError(f"No available port found above {preferred}.")


def docker_container_exists(name: str) -> bool:
    result = run(["docker", "ps", "-a", "--format", "{{.Names}}"], check=False, capture=True)
    return any(line.strip() == name for line in result.stdout.splitlines())


def remove_container_if_exists(name: str) -> None:
    if docker_container_exists(name):
        run(["docker", "rm", "-f", name], check=False, capture=True)


def ensure_network(name: str) -> None:
    result = run(["docker", "network", "inspect", name], check=False, capture=True)
    if result.returncode != 0:
        run(["docker", "network", "create", name])


def ensure_volume(name: str) -> None:
    result = run(["docker", "volume", "inspect", name], check=False, capture=True)
    if result.returncode != 0:
        run(["docker", "volume", "create", name])


def docker_image_exists(name: str) -> bool:
    result = run(["docker", "image", "inspect", name], check=False, capture=True)
    return result.returncode == 0


def build_images() -> None:
    print("Building backend image...")
    run(["docker", "build", "-t", SERVER_IMAGE, "./server"])
    print("Building frontend image...")
    run(["docker", "build", "-t", WEB_IMAGE, "./web"])


def wait_for_postgres(timeout_sec: int = 45) -> None:
    deadline = time.time() + timeout_sec
    while time.time() < deadline:
        result = run(
            ["docker", "exec", POSTGRES_CONTAINER, "pg_isready", "-U", "postgres", "-d", "task_tracker"],
            check=False,
            capture=True,
        )
        if result.returncode == 0:
            return
        time.sleep(1)
    raise RuntimeError("Postgres did not become ready in time.")


def up(args: argparse.Namespace) -> None:
    allow_auto = not args.strict_ports

    ports = Ports(
        postgres=pick_port(args.postgres_port, allow_auto=allow_auto),
        server=pick_port(args.server_port, allow_auto=allow_auto),
        web=pick_port(args.web_port, allow_auto=allow_auto),
    )

    # Avoid accidental duplicate port assignments when auto mode adjusts values.
    if len({ports.postgres, ports.server, ports.web}) < 3:
        raise RuntimeError(
            f"Resolved duplicate ports: postgres={ports.postgres}, server={ports.server}, web={ports.web}"
        )

    remove_container_if_exists(WEB_CONTAINER)
    remove_container_if_exists(SERVER_CONTAINER)
    remove_container_if_exists(POSTGRES_CONTAINER)

    ensure_network(NETWORK_NAME)
    ensure_volume(POSTGRES_VOLUME)

    if args.skip_build:
        if not docker_image_exists(SERVER_IMAGE):
            raise RuntimeError(f"Image {SERVER_IMAGE} not found. Re-run without --skip-build.")
        if not docker_image_exists(WEB_IMAGE):
            raise RuntimeError(f"Image {WEB_IMAGE} not found. Re-run without --skip-build.")
    else:
        build_images()

    try:
        print("Starting postgres container...")
        run(
            [
                "docker",
                "run",
                "-d",
                "--name",
                POSTGRES_CONTAINER,
                "--network",
                NETWORK_NAME,
                "-e",
                "POSTGRES_USER=postgres",
                "-e",
                "POSTGRES_PASSWORD=postgres",
                "-e",
                "POSTGRES_DB=task_tracker",
                "-p",
                f"{ports.postgres}:5432",
                "-v",
                f"{POSTGRES_VOLUME}:/var/lib/postgresql/data",
                POSTGRES_IMAGE,
            ]
        )

        print("Waiting for postgres readiness...")
        wait_for_postgres()

        print("Starting backend container...")
        run(
            [
                "docker",
                "run",
                "-d",
                "--name",
                SERVER_CONTAINER,
                "--network",
                NETWORK_NAME,
                "-e",
                f"PUBLIC_HOST=http://localhost:{ports.server}",
                "-e",
                "PORT=:8000",
                "-e",
                "DB_USER=postgres",
                "-e",
                "DB_PASSWORD=postgres",
                "-e",
                f"DB_HOST={POSTGRES_CONTAINER}",
                "-e",
                "DB_PORT=5432",
                "-e",
                "DB_NAME=task_tracker",
                "-e",
                "DB_SSLMODE=disable",
                "-e",
                "JWT_EXP=604800",
                "-e",
                "JWT_SECRET=dev-secret-change-me",
                "-p",
                f"{ports.server}:8000",
                SERVER_IMAGE,
            ]
        )

        print("Starting frontend container...")
        run(
            [
                "docker",
                "run",
                "-d",
                "--name",
                WEB_CONTAINER,
                "--network",
                NETWORK_NAME,
                "-e",
                f"BACKEND_URL=http://{SERVER_CONTAINER}:8000",
                "-e",
                f"NUXT_BACKEND_URL=http://{SERVER_CONTAINER}:8000",
                "-e",
                "NODE_ENV=production",
                "-e",
                "NITRO_HOST=0.0.0.0",
                "-e",
                "NITRO_PORT=3000",
                "-p",
                f"{ports.web}:3000",
                WEB_IMAGE,
            ]
        )
    except subprocess.CalledProcessError as exc:
        remove_container_if_exists(WEB_CONTAINER)
        remove_container_if_exists(SERVER_CONTAINER)
        remove_container_if_exists(POSTGRES_CONTAINER)
        raise RuntimeError(f"Failed to start stack: {' '.join(exc.cmd)}") from exc

    print("\nStandalone dev stack is running:")
    print(f"- Frontend: http://localhost:{ports.web}")
    print(f"- Backend:  http://localhost:{ports.server}")
    print(f"- Postgres: localhost:{ports.postgres} (postgres/postgres, db=task_tracker)")
    print("\nUseful commands:")
    print(f"- docker logs -f {SERVER_CONTAINER}")
    print(f"- docker logs -f {WEB_CONTAINER}")
    print(f"- python3 scripts/dev_docker_stack.py down")


def down(args: argparse.Namespace) -> None:
    remove_container_if_exists(WEB_CONTAINER)
    remove_container_if_exists(SERVER_CONTAINER)
    remove_container_if_exists(POSTGRES_CONTAINER)

    if args.remove_network:
        run(["docker", "network", "rm", NETWORK_NAME], check=False, capture=True)

    if args.remove_volume:
        run(["docker", "volume", "rm", POSTGRES_VOLUME], check=False, capture=True)

    print("Standalone dev containers stopped.")


def status(_: argparse.Namespace) -> None:
    result = run(
        [
            "docker",
            "ps",
            "-a",
            "--format",
            "table {{.Names}}\t{{.Status}}\t{{.Ports}}",
        ],
        capture=True,
    )
    lines = result.stdout.splitlines()
    if not lines:
        print("No containers found.")
        return

    tracked = {POSTGRES_CONTAINER, SERVER_CONTAINER, WEB_CONTAINER}
    filtered = [lines[0]] + [line for line in lines[1:] if line.split()[0] in tracked]
    print("\n".join(filtered))


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run Task Tracker dev stack in separate Docker containers.")
    subparsers = parser.add_subparsers(dest="command", required=True)

    up_parser = subparsers.add_parser("up", help="Build and start standalone dev containers.")
    up_parser.add_argument("--postgres-port", type=int, default=5543, help="Host port for postgres (default: 5543)")
    up_parser.add_argument("--server-port", type=int, default=8001, help="Host port for backend (default: 8001)")
    up_parser.add_argument("--web-port", type=int, default=3001, help="Host port for frontend (default: 3001)")
    up_parser.add_argument(
        "--strict-ports",
        action="store_true",
        help="Fail if requested ports are busy instead of auto-selecting next free ports.",
    )
    up_parser.add_argument("--skip-build", action="store_true", help="Skip docker image builds and reuse existing tags.")
    up_parser.set_defaults(func=up)

    down_parser = subparsers.add_parser("down", help="Stop and remove standalone dev containers.")
    down_parser.add_argument(
        "--remove-network",
        action="store_true",
        help=f"Also remove Docker network '{NETWORK_NAME}'.",
    )
    down_parser.add_argument(
        "--remove-volume",
        action="store_true",
        help=f"Also remove Postgres volume '{POSTGRES_VOLUME}'.",
    )
    down_parser.set_defaults(func=down)

    status_parser = subparsers.add_parser("status", help="Show container status table.")
    status_parser.set_defaults(func=status)

    return parser.parse_args()


def main() -> int:
    require_docker()
    args = parse_args()
    try:
        args.func(args)
    except RuntimeError as exc:
        print(str(exc), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
