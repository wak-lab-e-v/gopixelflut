Pixelflut: Multiplayer canvas for LED Boards
===============================
#### Pixelflut Socket server in GoLang.

Pixelflut Protocol
------------------

Pixelflut defines four main commands that are always supported to get you started:

* `HELP`: Returns a short introductional help text.
* `INFO`: Returns some informations and statistics. 
* `SIZE`: Returns the size of the visible canvas in pixel as `SIZE <w> <h>`.
* `GP <x> <y>`: Return the current color of a pixel as `PX <x> <y> <r> <g> <b>`. 
* `GM`: Return the the full matrix as binary stream. Row oriented. Start with Y1.
* `PX <x> <y> #<rrggbb>`: Draw a single pixel at position (x, y) with the specified hex color code.
* `PX <x> <y> 255 255 255`: Draw a single pixel at position (x, y) with the specified color values (R G B).
* `EXIT`: Close your session by server.

You can send multiple commands over the same connection by terminating each command with a single newline character (`\n`).

Example:

    $ echo "SIZE" | netcat pixelflut.example.com 1337
    SIZE 800 600
    $ echo "PX 10 10 #ff8000" | netcat pixelflut.example.com 1337
    $ echo "PX 20 20 255 50 0 | netcat pixelflut.example.com 1337
	$ echo "GP 20 20" | netcat pixelflut.example.com 1337
    PX 20 20 255 50 0


## License
#### Available under the MIT License, please see LICENSE.
