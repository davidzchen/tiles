# tiles

Easy way to manage tmux sessions

## Installation

To install, simply put `tiles` in your `PATH.`

## Usage

```sh
tiles <command> [<flags>]
```

Simple wrapper script for creating and managing tmux (terminal multiplexer)
sessions.

Possible commands are:

* `start`  - Start a new tmux session.
* `attach` - Attach to an existing tmux session.
* `ls`     - List currently running tmux sessions.
* `help`   - Prints the help text.

Tmux sessions are defined with a `.tiles` file, which must be in the user's home
directory. The syntax of the `.tiles` file is as follows:

```python
tmux_session(
    name = "session-name,
    windows = [],
)
```

The `name` parameter is used to reference the session when invoking this script.
The `windows` parameter is a list of tuples of `[window_name,
working_directory]` For example, the following configuration:

```python
tmux_session(
    name = 'work',
    windows = [
        ['blog', '~/Projects/blog'],
        ['tensorflow', '~/Projects/tensorflow']
    ],
)
```

defines a tmux session named 'work' with two windows.

## Roadmap

* Check existence of directories before running tmux commands
* Add tests
* Long-term: Support GNU Screen
* Long-term: Support configuring panes within windows

## License

`tiles` is licensed under the Apache 2.0 license:

```
Copyright 2015 David Z. Chen

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
