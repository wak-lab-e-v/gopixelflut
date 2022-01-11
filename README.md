Pixelflut: Multiplayer canvas
===============================
#### Pixelflut Socket server in GoLang.

Pixelflut Protocol
------------------

Pixelflut defines four main commands that are always supported to get you started:

* `HELP`: Returns a short introductional help text.
* `INFO`: Returns size and some other statistics. 
* `SIZE`: Returns the size of the visible canvas in pixel as `SIZE <w> <h>`.
* `GP <x> <y>` Return the current color of a pixel as `PX <x> <y> #<rrggbb>`. 
* `PX <x> <y> #<rrggbb>`: Draw a single pixel at position (x, y) with the specified hex color code.
* `PX <x> <y> 255 255 255`: Draw a single pixel at position (x, y) with the specified color values (R G B).

You can send multiple commands over the same connection by terminating each command with a single newline character (`\n`).

Example:

    $ echo "SIZE" | netcat pixelflut.example.com 1337
    SIZE 800 600
    $ echo "PX 23 42 #ff8000" | netcat pixelflut.example.com 1337
    $ echo "GP 32 42" | netcat pixelflut.example.com 1337
    PX 23 42 #ff8000


## License
#### Available under the MIT License, please see LICENSE.
