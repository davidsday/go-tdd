NOTE:  This plugin is under active development.  I try to keep the
documentation current. But things change not always in real time.

This is a Vim plugin that seeks to bring two things
to the Vim golang development experience:

	RedBar/GreenBar/Refactor Test Driven Development style programming

	A better 'go test' experience than I have received with some other
	tools.

The first is via the use of Vim's message line to provide red/green
bar fail/pass indications. There is also a yellow bar message which
is used to convey messages about errors that are not caused by an
actual failing test.

Right now these yellow bar messages include, [no tests found],
[build failed], [received a panic], and [invalid JSON message].

There also are messages providing detail information in each red/green
bars.  They report the number of tests run, passed, failed, and skipped,
in addition to the elapsed time for running all the tests as provided
by go test.

Much, if not all of this information is absent from the rather sparce
reporting done by vim-go, which is by far the leading vim/neovim golang
development plugin, one that I use everyday and cannot imagine
programming in golang without.  But I have seen many instances where
vim-go reports [SUCCESS] when there actually were not tests run at all,
or when tests were skipped with no notification to the programmer.

My thinking is that if I am a consultant called in to work on a code base,
I do not want my tools delivering erroneous, overly optimistic reports
to me.  To me, if a package has 100 tests, but 25 are not even being run,
I want to know that right up front.

goTestParser is designed to work alongside of vim-go.
It does not interfere with vim-go in anyway.

goTestParser provides its own go test parser, written in golang, which
parses the 'go test -v -json' output and in turn, provides a much further
processed JSON structure which details for a small Vimscript what
message, and in what color to deliver.  It also provides to Vim a quickfix
list of test failures and/or skipped tests which Vim loads for your use.



