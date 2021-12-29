Pixelflut: Multiplayer canvas
===============================
#### Pixelflut Socket server in GoLang.

Pixelflut Protocol
------------------

Pixelflut defines four main commands that are always supported to get you started:

* `HELP`: Returns a short introductional help text.
* `SIZE`: Returns the size of the visible canvas in pixel as `SIZE <w> <h>`.
* `PX <x> <y>` Return the current color of a pixel as `SP <x> <y> <rrggbb>`.
* `PX <x> <y> #<rrggbb(aa)>`: Draw a single pixel at position (x, y) with the specified hex color code.
  If the color code contains an alpha channel value, it is blended with the current color of the pixel.

You can send multiple commands over the same connection by terminating each command with a single newline character (`\n`).

Example:

    $ echo "SIZE" | netcat pixelflut.example.com 1337
    SIZE 60 33
    $ echo "SP 23 42 #ff8000" | netcat pixelflut.example.com 1337
    $ echo "GP 32 42" | netcat pixelflut.example.com 1337
    SP 23 42 #ff8000

## License
#### Available under the MIT License, please see LICENSE.
