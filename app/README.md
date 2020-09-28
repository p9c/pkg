# app

App is a library wrapped around 
github.com/urfave/cli that makes neater 
declarations, dynamically generates a 
concurrent safe, tag-grouped configuration 
database with the change hooks that sanitise 
values, can trigger running responses such
as restarting a subsystem or as needed, then
efficiently modifies a cached, marshalled form 
of the configuration (JSON) and syncs to disk,
in the configured data directory.

See [example folder](example/) for an example. The
example can easily be used to start a new app.

In addition here in this repo are platform 
specific data directory code, various generic 
filesystem helpers, and an interrupt handling
library to enable clean shutdown and restarts.

## Rationale

When writing applications, there is normally
a constant need to change, add and remove
config and subcommands before the app is
complete. With complex, pluggable 
applications, subcommands allow a simpler, 
monolithic runtime that can launch child
processes, for use with pipe IPC. 

App centralises the specification, and 
automatically generates a config system
with caching and change hooks to allow hot
dynamic reconfiguration, and simple context
access to subsystems to read and potentially
write configuration data concurrently.
