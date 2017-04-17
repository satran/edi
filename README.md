**This project is abandoned as I currently do not have time to work on it. Please feel free to fork it.**

# EDI

EDI is an IDE. Rather than an Integrated Development Environment it strives to be an Integrating Development Environment in the spirit of Acme. EDI does not try to do it all. It rather leans on the wonderful tools that already exist on your machine. It acts as a bridge between those commands and the editor.


### REQUIREMENTS
You just need [Go](http://golang.org/doc/install) and a recent browser that supports websockets. It has been tested on the latest Firefox and Chrome versions.
Tested both on GNU/Linux and Mac. Not tested on Windows(I'm sure there are parts that will not work.)


### INSTALLATION
To install make sure you have set `$GOPATH` and `go` is in your `$PATH`. Execute

    sh INSTALL.sh
    
This should build a static binary for your OS as `$GOPATH/bin/edi`.


### USAGE
EDI runs a webserver on port 8312. Open your browser on that port. 
The command view accepts these commands:
- `E file_name` to open file (relative paths also allowed)
- `W` saves the current buffer
- `W file_name` saves the current buffer as `file_name`
- `cd` to change directory
- `L` to display the open buffers
- `Ctrl+l` clears the command view

Right clicking on a file name in the command view performs what is called a context. Currently context opens a valid file in the editor view.

`Ctrl+space` toggles focus between the command view and the editor view.

You can execute many of the commands your OS provides as long as it does not require standard input or ncurses.


### NOTES
- This is an alpha version. Things will break. I do not take responsibility of your lost code.
- There is no customization at the moment. This will change in the next release. But for now edit the code to customize.
- Keybindings for the editor is currently vim. This can be changed by editing `app/s/j/edi.js`. emacs, sublime-text key bindings are available.
- Its probably not secure so don't try running it on a server that is facing the internet.
