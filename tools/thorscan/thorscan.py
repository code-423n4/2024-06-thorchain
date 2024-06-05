import concurrent.futures
import datetime
import json
import logging
import os
import signal
import sys
import threading
import time
from queue import Queue
from typing import Any, Callable, Dict, Optional, Set

import requests
import retry

########################################################################################
# Config
########################################################################################

BLOCK_SECONDS = 6

API_ENDPOINT = os.getenv("API_ENDPOINT", "https://thornode-v1.ninerealms.com")
RPC_ENDPOINT = os.getenv("RPC_ENDPOINT", "https://rpc-v1.ninerealms.com")
PARALLELISM = int(os.getenv("PARALLELISM", 4))

logging.basicConfig(
    format="%(asctime)s - %(filename)s:%(lineno)d - %(message)s",
    level=logging.INFO,
    stream=sys.stderr,
)


########################################################################################
# Internal Helpers
########################################################################################


def _parse_block_time(dt_str):
    dt = datetime.datetime.strptime(dt_str[:26], "%Y-%m-%dT%H:%M:%S.%f")
    timestamp = dt.replace(tzinfo=datetime.timezone.utc).timestamp()
    return timestamp


def _print(ret: Optional[str]):
    """
    Print the provided value if it is not None.
    """
    if ret is None:
        return
    if isinstance(ret, (list, tuple)):
        print(json.dumps(ret))
    elif isinstance(ret, dict):
        print(json.dumps(ret, indent=2))
    else:
        print(ret)


########################################################################################
# Listeners
########################################################################################


def transactions(
    *listeners: Callable[[Dict[str, Any], Dict[str, Any]], Optional[str]],
    failed: bool = False,
):
    """
    Returns a block listener which executes provided transaction listeners.

    Args:
        listeners: listener functions accepting (block, tx)
        failed: listen to failed transactions instead of successful ones
    """

    def _listen(block):
        for tx in block["txs"]:
            if failed and tx["result"]["code"] == 0:
                continue
            elif tx["result"]["code"] != 0:
                continue
            for listener in listeners:
                _print(listener(block, tx))

    return _listen


def events(
    *listeners: Callable[
        [Dict[str, Any], Optional[Dict[str, Any]], Dict[str, Any]], Optional[str]
    ],
    types: Optional[Set[str]] = None,
) -> Callable[[Optional[Dict[str, Any]]], Optional[str]]:
    """
    Returns a block listener which executes provided event listeners.

    Args:
        listeners: listener functions accepting (block, tx, event)
        types: event types to filter
    """

    def _listen(block):
        for event in block["begin_block_events"] + block["end_block_events"]:
            if types is not None and event["type"] not in types:
                continue
            for listener in listeners:
                _print(listener(block, None, event))
        for tx in block["txs"]:
            for event in tx.get("events", []):
                if types is not None and event["type"] not in types:
                    continue
                for listener in listeners:
                    _print(listener(block, tx, event))

    return _listen


def messages(
    *listeners: Callable[
        [Dict[str, Any], Optional[Dict[str, Any]], Dict[str, Any]], Optional[str]
    ],
    types: Optional[Set[str]] = None,
) -> Callable[[Optional[Dict[str, Any]]], Optional[str]]:
    """
    Returns a block listener which executes provided message listeners.

    Args:
        listeners: listener functions accepting (block, tx, message)
        types: event types to filter
    """

    def _listen(block):
        for tx in block["txs"]:
            for msg in tx["tx"]["body"]["messages"]:
                if types is not None and msg["type"] not in types:
                    continue
                for listener in listeners:
                    _print(listener(block, tx, msg))

    return _listen


########################################################################################
# Scan
########################################################################################

stop_scan = False
last_scan_block_time = 0


def signal_handler(signal, frame):
    global stop_scan
    stop_scan = True


signal.signal(signal.SIGINT, signal_handler)

# create session pool
sessions = Queue()
for _ in range(PARALLELISM):
    sessions.put(requests.Session())


def _get(*args, **kwargs) -> requests.Response:
    # set x-client-id header
    kwargs["headers"] = kwargs.get("headers", {})
    kwargs["headers"]["x-client-id"] = "thorscan"

    # 5 second timeout
    kwargs["timeout"] = kwargs.get("timeout", 5)

    # get session from pool for request
    session = sessions.get()
    res = session.get(*args, **kwargs)
    sessions.put(session)

    return res


@retry.retry(delay=1, backoff=2, max_delay=10)
def _current_height() -> int:
    res = _get(f"{RPC_ENDPOINT}/status").json()
    return int(res["result"]["sync_info"]["latest_block_height"])


@retry.retry(delay=BLOCK_SECONDS / 2, backoff=2, max_delay=BLOCK_SECONDS)
def _fetch_block(height: int) -> Optional[Dict[str, Any]]:
    while True:
        if stop_scan:
            return None
        res = _get(f"{API_ENDPOINT}/thorchain/block", params={"height": height})
        if res.status_code == 404:
            time.sleep(BLOCK_SECONDS / 2)
            continue  # expected near tip
        res.raise_for_status()
        return res.json()


def _fetch_thread(queue: Queue, start: int, stop: int):
    height = start
    with concurrent.futures.ThreadPoolExecutor(PARALLELISM) as pool:
        while not stop_scan:
            # wait when last scan block is too recent (avoid 404s)
            if last_scan_block_time > time.time() - BLOCK_SECONDS:
                time.sleep(BLOCK_SECONDS)
                continue

            # wait if queue is full or we are caught up
            bootstrapping = last_scan_block_time == 0
            caught_up = last_scan_block_time > time.time() - BLOCK_SECONDS * PARALLELISM
            if queue.full() or (queue.qsize() > 0 and (caught_up or bootstrapping)):
                time.sleep(0.1)
                continue

            fut = pool.submit(_fetch_block, height)
            queue.put(fut)
            height += 1
            if stop != 0 and height > stop:
                queue.put(None)
                return


def scan(
    *listeners: Callable[[Optional[Dict[str, Any]]], Optional[str]],
    start: int = 0,
    stop: int = 0,
):
    """
    Scan blocks and call all listeners on every block in the range.

    Args:
        listeners: listener functions accepting the block
        start: height to start scan, current height if zero, reverse index if negative
        stop: height to stop scan, infinite if zero, reverse index if negative
    """
    if len(listeners) == 0:
        raise Exception("no listeners were provided to scan")

    # reverse index from current height for negative start or stop
    if start < 0 or stop < 0:
        current_height = _current_height()
        if start < 0:
            start = current_height + start
        if stop < 0:
            stop = current_height + stop

    # print to stderr
    print(
        f"scanning blocks {start} to {stop or 'infinite'} with parallelism {PARALLELISM}",
        file=sys.stderr,
    )

    # start queue fetching in separate thread
    queue = Queue(maxsize=int(PARALLELISM))
    fetch_thread = threading.Thread(target=_fetch_thread, args=(queue, start, stop))
    fetch_thread.start()

    while not stop_scan:
        try:
            fut = queue.get(timeout=0.1)
            queue.task_done()
        except Exception:
            continue

        # check if we are done
        if fut is None:
            break

        # wait for block result
        while True:
            if fut.done():
                break
            time.sleep(0.1)
        block = fut.result()
        if stop_scan or block is None:
            return

        # if block is less than 1 minute old, stop parallel fetching
        global last_scan_block_time
        last_scan_block_time = _parse_block_time(block["header"]["time"])

        # modify block transaction type field for convenience
        for tx in block["txs"]:
            for msg in tx["tx"]["body"]["messages"]:
                msg["type"] = msg["@type"].lstrip("/types.")
                del msg["@type"]

        # call all block listeners
        for listener in listeners:
            _print(listener(block))


########################################################################################
# CLI
########################################################################################


def _cli():
    args = " ".join(sys.argv[1:])
    eval(f"scan({args})")
