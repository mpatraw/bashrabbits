
# Rabbit

Rabbit is a game that is played in your file system. The goal is to catch and tag as many rabbits as possible.

## How to Install

If you have Go installed, run `go get github.com/mpatraw/rabbit`. Otherwise you will have to download one of the binary packages when they're put up.

## How to Play

`rabbit` offers a couple of commands to find and catch rabbits. The first and most obvious of which is `rabbit check` which checks the current directory for a rabbit. If there is one, you will see "A rabbit is here!!". You only have a few seconds to either catch it with the `rabbit catch` command or tag it with `rabbit tag <string>` command. One last command `rabbit stats` tells you how many rabbits you've seen, caught, and... killed. Yes, you can kill rabbits, mostly by accident.

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

Or if you use __rc__.

```rc
fn cd { builtin cd $*; rabbit check }
fn ls { /bin/ls $*; rabbit check }
fn catchr { rabbit catch }
fn tagr { rabbit tag }
```

So, you must be wondering, "how do you kill rabbits?" Well, when you remove or move directories around, rabbit may have lived there and now dead.

## How does it Work?

Don't worry, there aren't __actually__ rabbits in your directories. The program keep a record of where every rabbit is and it's state in `$HOME/.rabbit`, and moves and spawns new ones when necessary.
