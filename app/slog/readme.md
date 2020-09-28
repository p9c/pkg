# slog

Flexible, extensible logging framework designed to be extended arbitrarily to
enable any kind of log output or handling.

## Rationale

Concurrent programming is quite unsuited to the use of debuggers as they are
generally only designed to follow a single path of execution. Much more useful is
to have the ability to watch the activity of multiple threads of execution.

The standard library for logging is very simple and does not lend itself to easy
extension, unlike the error library, which can be extended as it defines its
function set in an interface. 

The debugging method that the author of this library uses involves making use
of the hypertext linking functions and opening of source code at code locations
relevant to an error. The error library has been extended and plans exist to
augment it such that errors can now be even further deeply embedded in multiple
layers of execution. This is good and well for errors that can be handled, but it
just further obstructs the ease of eliminating unexpected errors caused by 
incorrect algorithms.

As such, this library comes with a minimal implementation that prints out
code locations of the site of the error so that the developer can immediately 
jump to the site of the problem and trace it more easily back to its origin.
