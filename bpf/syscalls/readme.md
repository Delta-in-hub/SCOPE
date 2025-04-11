Trace syscall entries using BPF.

Filters events based on PID and/or command name and prints details.

USAGE: ./syscalls [-p PID] [-c COMM] [-v]

  -c, --comm=COMMAND         Filter by command name (exact match)
  -p, --pid=PID              Filter by process ID (TGID)
  -v, --verbose              Verbose debug output
  -?, --help                 Give this help list
      --usage                Give a short usage message
  -V, --version              Print program version

