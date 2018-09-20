# LedFS
> A daemon for driving WS281x LEDs (Neopixels) using a FUSE filesystem
> on the Raspberry Pi

### Usage
Follow the build instructions below to install the LedFS Daemon on the Raspberry Pi of your
choice. 

##### Config
Config is provided within `<mountpoint>/options.json` and automatically resets the program
when changed. When reset, all leds are zeroed out, and the connection to the Neopixels
is toggled off and on again.

Defaults to:
```json
{
    "numLeds": 24,
    "gpioPin": 18,
    "brightness": 220
}
```

##### Colors
Colors are set using `<mountpoint>/colors.json` and automatically updated with the file is
modified. Colors should be formated as an JSON array of hex color code strings like so:
```json
{ "values": ["#FAFAFA", "#000000", "#1a0000"] }
```

Note that array values map directly to the Neopixels. If you have 24 leds, but only provide
an array of three colors, only the first three leds will be colored.


### Build setup

In order to compile the ledfs binary, you will need to setup a few things on
Raspberry Pi. Once installed, using LedFS is as simple as running the executable
like so: `ledfsd /directory/to/mount`. Please note that you may need to run the 
executable as root (since it needs access to the GPIO). While running, LedFS enables
control over Neopixels using the configuration as found in `<mountpoint>/options.json` 
and the array of hex colors found in `<mountpoint>/colors.json`.

+ Go: 
	This should do the trick...

	```sh
	sudo apt install golang
	mkdir -p $HOME/go
	echo "export GOPATH=$HOME/go" > .profile
	. .profile
	```


+ WS281x C and Go libs: 

	Copy the libraries locationed in the `deps` directory to the path they are in minus the
	containing folder. For example files in, `deps/usr/local/lib`
	should be placed in `/usr/local/lib`.

+ Root Priviledges: 
	We'll leave that up to you and the internet.


+ Fuse Libraries: 

    This will install libfuse2 and the binaries required for mounting. 
    
    `sudo apt install fuse`


As long as you have all of the above, compiling your go source like any other
Go code should work fine. Just make sure to compile and then run the resulting
binary as root, because the leds require root access to run, and funky errors
appear if you aren't root.
