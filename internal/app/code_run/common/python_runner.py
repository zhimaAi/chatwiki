"""Trusted bootstrap embedded in code_run. User source arrives only through stdin."""
import builtins
import datetime
import errno
import json
import math
import os
import platform
import random
import re
import sys
import typing
import _strptime  # datetime.strptime imports this lazily; preload trusted code.

SAFE_BUILTINS = ["abs","all","any","bin","bool","chr","complex","dict","divmod","enumerate","filter","float","format","frozenset","hex","int","isinstance","issubclass","iter","len","list","map","max","min","next","oct","ord","pow","print","range","repr","reversed","round","set","slice","sorted","str","sum","tuple","type","zip","Exception","PermissionError","TypeError","ValueError","KeyError","AttributeError","IndexError","RuntimeError"]

def import_runtime_helper(name, globals=None, locals=None, fromlist=(), level=0):
    # CPython datetime calls __import__ even when _strptime is preloaded.
    # Return only the already-loaded trusted helper; never load arbitrary code.
    if name == "_strptime" and level == 0:
        return _strptime
    raise ImportError("module import is disabled")

def install_guard():
    if sys.platform != "linux" or platform.machine() != "x86_64":
        raise RuntimeError("unsupported sandbox platform")
    import pyseccomp as seccomp
    # All library discovery/import activity finishes before restricting syscalls.
    for entry in os.listdir("/proc/self/fd"):
        fd = int(entry)
        if fd > 2:
            try:
                os.close(fd)
            except OSError as exc:
                if exc.errno != errno.EBADF:
                    raise
    policy = seccomp.SyscallFilter(defaction=seccomp.ALLOW)
    policy.set_attr(seccomp.Attr.CTL_NNP, 1)
    policy.set_attr(seccomp.Attr.CTL_TSYNC, 1)
    blocked = (
        "fork", "vfork", "clone", "clone3", "execve", "execveat",
        "socket", "socketpair", "connect", "bind", "listen", "accept", "accept4",
        "sendto", "recvfrom", "sendmsg", "recvmsg", "sendmmsg", "recvmmsg",
        "getsockopt", "setsockopt", "getsockname", "getpeername", "shutdown",
        "io_uring_setup", "io_uring_enter", "io_uring_register",
        "ptrace", "process_vm_readv", "process_vm_writev", "pidfd_getfd",
    )
    for name in blocked:
        policy.add_rule(seccomp.ERRNO(errno.EPERM), name)
    policy.load()

def run():
    # Never handle initialization failure as an optional feature.
    try:
        install_guard()
    except Exception:
        print(json.dumps({"error": "python_sandbox_init_failed"}))
        return 1
    try:
        request = json.loads(sys.stdin.buffer.read().decode("utf-8"))
        source, params = request["main_func"], request["params"]
        if not isinstance(source, str) or not isinstance(params, dict):
            raise ValueError("invalid python request")
        scope = {
            "__builtins__": {
                **{name: getattr(builtins, name) for name in SAFE_BUILTINS},
                "__import__": import_runtime_helper,
            },
            "__name__": "user_code",
            "json": json, "math": math, "random": random,
            "datetime": datetime, "re": re,
        }
        for name in ("Any", "Dict", "List", "Tuple", "Optional"):
            scope[name] = getattr(typing, name)
        exec(compile(source, "<user_code>", "exec"), scope, scope)
        result = scope["main"](**params)
        print(json.dumps(result))
        return 0
    except BaseException as exc:
        print(json.dumps({"error": str(exc)}))
        return 1

if __name__ == "__main__":
    sys.exit(run())
