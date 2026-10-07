# Why io.Writer
Setup chooses the destination once. Both os.Stdout and *os.File implement io.Writer. Tests can use an in-memory bytes.Buffer.

# Why MultiWriter has the console first
MultiWriter writes to destinations in order and stops at the first error. With the console first, a log line still reaches the terminal if the file write fails. With the file first, that error would stop the line before it reached the terminal.

# Why Close uses a stored file
In ModeBoth, destination is a MultiWriter, not an *os.File, so a file type check failed. In ModeStdout, a path is ignored and destination is os.Stdout, so the old check would have closed the terminal.

# What the mutex protects
It keeps each ModeBoth line's console and file writes together. It also stops Close from closing the file during a log write. The concurrent test passes without the mutex because it only writes to a file and counts lines; it does not race Close against logging or check ordering across destinations.

# What's out of scope on purpose
Colors, JSON output, log rotation, and asynchronous logging are out of scope. Changing the level while logging is also out of scope; it would need the mutex too.

# Close decision
Close is idempotent. A repeat call returns nil.