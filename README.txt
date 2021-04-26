NOTE:  This plugin is under active development.  I try to keep this
current. But things change....

This is a Vim plugin that seeks to bring two things
to the Vim golang development experience:

	RedBar/GreenBar/Refactor Test Driven Development style programming

	A better 'go test' experience than I have received with some other
	tools.

The first is via the use of Vim's message line to provide red/green
bar fail/pass indications. There is also a yellow bar message which
is used to convey messages about errors that are not caused by an
actual failing tests.

Right now these yellow bar messages include, [no tests found],
[build failed], [received a panic], and [invalid JSON message],
and the receipt of a message on STDERR.  The build tools sometimes
will issue important messages on STDERR, and I want to know about
them immediately.  These messages are not in JSON format.  They are
issued from the build system tools in their normal output formats
and would end up interspersed with the go test -json output, leading
to goTestParser raising an error due to invalid JSON in its input stream.

There also are messages providing detail information in each red/green
bars.  They report the number of tests run, passed, failed, and skipped,
in addition to the elapsed time for running all the tests as provided
by go test. Test coverage is reported on Green Bars.  Presumably, you
have more to worry about than that on Yellow or Red Bars.

Right now, goTestParser is taking go test -json at its word as to
how many tests are run, passed failed etc, based on the JSON Action
fields.  These are wrong sometimes, since they count a main test that
has subtests right along with the subtests.  It is quite common for
a test to kick off subtests which do the actual testing.  go test -json
counts the main test even though it does no testing itself. So a main
test with 5 subtests count as 6 tests, which is incorrect.  It should
be possible to discern what the actual count should be, but basically
requires going back to parsing the non JSON go test -v output, thus
largely defeating the point of converting to -json in the first place.

For now, I am just taking the go test -json output's word and we need to
realize that the results are approximate.

I have seen many instances where vim-go reports [SUCCESS] when there
actually were not tests run at all, or when tests were skipped with
no notification to the programmer.

My thinking is that if I am a consultant called in to work on a code base,
I do not want my tools delivering overly optimistic reports
to me.  To me, if a package has 100 tests, but 25 are not even being run,
I don't want the tools to report that as [SUCCESS]. If go test issues
messages via STDERR, I want to know that.  The fact that they were issued
on STDERR is info I want to be explicitly told.

goTestParser is designed to work alongside of vim-go, since, really,
vim-go is my most important golang development tool.

It does not interfere with vim-go in anyway that I am aware of.
In my setup, I have replaced vim-go's <Leader>t (<ESC>:GoTest<CR) with
<LocalLeader>t to activate goTestParser. If I desire to use vim-go's
:GoTest command, I call it just like that.

goTestParser provides its own go test parser, written in golang, which
parses the 'go test -v -json' output and in turn, provides a further
processed JSON structure which details for a small Vimscript what
message, and in what color to deliver.  It also provides to Vim a quickfix
list of test failures and/or skipped tests which Vim loads for your use.

In this style of development, the RedBar/GreenBar (and YellowBar)s are the
primary layer of communication with the developer, so goTestParser loads
the quickfix list for you, but it does not open it and take you away from
the file you have open.  I find that in my development work flow,
I often don't need to go to the test at all, but want to peruse  and
fix the function that caused the failure, and I may already be there.

The RedBar/GreenBar/YellowBar message line lingers until any key is
pressed (I typically just hit the space bar).

In my set up, <Leader>q toggles the quickfix window and <C-j> and <C-k>
navigate the quickfix list up and down. <LocalLeader>a takes me to the
Alternate file.

