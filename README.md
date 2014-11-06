
# Rabbit

Rabbit is a game that is being played all of the time in your file system. The goal is to catch and tag as many rabbits as possible.

## How to Install

If you have Go installed, running `go get github.com/mpatraw/Rabbit`. Otherwise you will have to download one of the binary packages when they're put up.

## How to Play

`rabbit` offers a couple of commands to find and catch rabbits. The first and most obvious of which is `rabbit check` which checks the current directory for a rabbit. If there is one, which you will know by a message telling you so, you only have a few seconds to either catch it with the `rabbit catch` command or tag it with `rabbit tag blue` command. One last command `rabbit stats` tells you how many rabbits you've seen, caught, and... killed. Yes, you can kill rabbits, mostly by accident.

Obviously typing these out everytime you're in a directory is tiring, so you can add this to your `.bashrc` file.

```bash
rcd() {
	cd $*
	rabbit check
}
rls() {
	ls $*
	rabbit check
}
alias cd='rcd'
alias ls='rls'
```
