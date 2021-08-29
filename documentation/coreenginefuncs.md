# Core Engine Functions

## Functions Overview

Function | Description
-------- | -----------
[shutdown()](#shutdown) | Shuts down the engine (exits to desktop)
[log(level: string, message: string)](#log) | Writes a log entry message based on the log level
[fmt(format: string, values...: any)](#fmt) | Returns a formatted string from the input format and values
[setBootText(text: string)](#setBootText) | Sets the lower text in the boot splash screen

## Function Details

### shutdown
```lua
shutddown()
```
Shuts down the engine (exits to desktop)

---

### log
```lua
log(level, message)
```
Writes a log entry message based on the log.

The first parameter is a string that defines the log level based on the below list.

The second parameter is a single string that contains the message to log. Note that
if you'd like to have a formatted string, you should use the [fmt](#fmt) function, which
can be used in-line as the second parameter if necessary.

Log Levels:
* info
* error
* fatal
* warn
* debug
* trace

---

### fmt
```lua
fmt(string, value, ...)
```

Returns a formatted string from the input format and values

The format string is based on the [golang format verbs](https://pkg.go.dev/fmt).


---

### setBootText
```lua
setBootText()
```

Sets the lower text in the boot splash screen.

This function is only useful during the boot phase of the mod. After the boot phase
is over, this function does nothing.


