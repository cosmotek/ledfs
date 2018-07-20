# LedFS
> A daemon for driving WS281x LEDs (Neopixels) using a FUSE filesystem
> on the Raspberry Pi

#### Build setup

In order to compile the ledfs binary, you will need to setup a few things on
Raspberry Pi.

+ Go:
	This should do the trick...
	`sudo apt install golang`
	`echo "export GOPATH=$HOME/go" > .profile"`
	`. .profile`

+ WS281x C and Go libs
	Copy the libraries locationed in the `deps` directory to the path they are in minus the
	containing folder. For example files in, `deps/usr/local/lib`
	should be placed in `/usr/local/lib`.

+ Root Priviledges
	We'll leave that up to you and the internet.

As long as you have all of the above, compiling your go source like any other
Go code should work fine. Just make sure to compile and then run the resulting
binary as root, because the leds require root access to run, and funky errors
appear if you aren't root.
