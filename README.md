
# Rabbit

Rabbit is a game that is played in your file system. The goal is to catch and tag as many rabbits as possible, while avoiding killing too many.

## How to Install

If you have Go installed, you can run `go get github.com/mpatraw/rabbit`. Otherwise you will have to download one of the binary packages when they're put up.

## How to Play

`rabbit` offers a couple of commands to find and catch rabbits. The first and most obvious of which is `rabbit check` which checks the current directory for a rabbit. If there is one, you will see "A rabbit is here!!". You only have a few seconds to either catch it or tag it.

Only a few rabbits will exist at any given time (1-15), all of which will never go below your home directory ($HOME). Generally the rabbits will move about an area slowly, only doing 1-2 directory hops every few minutes. The only time they move quickly is when they're spotted, once they leave (a few seconds later) they could be almost anywhere in your home tree.

Killing rabbits is a danger. You can kill a rabbit (accidentally or intentionally) by destroying where they are, either by a `mv` or a `rm -r`. One goal is to avoid killing as many rabbits as possible, but sometimes it's just unavoidable. :(

### Flags & Commands

__Flags__
* -a: Adds ASCII graphics at the end of commands.

__Commands__
* check: Checks the current directory for a rabbit.
* catch: Attempts to catch a rabbit in the current directory.
* tag "string": Tries to tag the rabbit in the current directory with "string".
* stats: Prints the stats of rabbits seen, caught, killed, etc.

### Extras

Obviously typing these out everytime you're in a directory is tiring, so you can add this to your `.bashrc` file.

```bash
cdrabbit() {
	cd $*
	rabbit check
}
lsrabbit() {
	ls $*
	rabbit check
}
alias cd='cdrabbit'
alias ls='lsrabbit'
alias tagr='rabbit tag'
alias catchr='rabbit catch'
```

Or if you use __rc__, add this to `.rcrc`.

```rc
fn cd { builtin cd $*; rabbit check }
fn ls { /bin/ls $*; rabbit check }
fn catchr { rabbit catch }
fn tagr { rabbit tag $* }
```

### Example Session

```
; ls
go  mine  others
; cd mine
; cd ..
; cd others
A rabbit is here!!
(_/  _#
'.'_( )
; catch
You caught the rabbit!
_________
| ()|() |
+---+---+
|(")|(")|
---------
; rabbit stats
Rabbits
...spotted:    10
...caught:     4
...killed:     1
; ls
dwm
; cd ..
; ls
go  mine  others
; cd ..
A rabbit is here!!
(_/  _#
'.'_( )
; catch
The rabbit got away...
  o __(\\
   ) _ --
 //    \\
;
```

## How does it Work?

Don't worry, there aren't __actually__ rabbits in your directories. The program keep a record of where every rabbit is and it's state in `$HOME/.rabbit`, and moves and spawns new ones when necessary.

## Thanks

To `bh`, for some of the rabbit ASCII art.
