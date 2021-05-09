NOTE:  This plugin is under active development.  I try to keep this
current. But things change....

This is a Vim plugin that seeks to bring two things
to the Vim golang development experience:

	RedBar/GreenBar/Refactor Test Driven Development style programming

	A better 'go test' experience

The first is via the use of Vim's message line to provide red/green
bar fail/pass indications. I have also added a yellow bar message which
I use to convey messages about errors that are not caused by an actual
failing test.

Right now these yellow bar messages include
	[no tests found],
	[no tests to execute] - basically an empty test file
	[build failed],
	[received a panic],
	[invalid JSON message],
	and the receipt of any message on STDERR.

The build tools sometimes will issue important messages on STDERR,
and I want to know about them immediately.

These messages are not in JSON format.  Since I am parsing
"go test -v -json -cover" output, goTestParser expects valid JSON.
I capture stderr and process it as best I can, by showing a snippet of
the message in a yellow bar.  If the STDERR message is longer than can
be shown in a one line yellow bar, I capture the entire message in
stdERR.txt in the package directory.

If goTestParser encounters non JSON lines on stdout, it issues a yellow
bar message and quits. That rarely, if ever happens.

There also are supplemental messages providing detail information in each
red/green bars.  They report the number of tests run, passed, failed,
and skipped, in addition to the elapsed time for running all the tests
as provided by go test. Test coverage is reported on Green Bars.

I have also added Average Cyclomatic Complexity to the Green Bars.
This has little to do with testing but a lot to do with design and
it is a metric I want to be aware of. I hear that several well known IDEs
start warning about Cyclomatic Complexity at 10.  I like to keep mine
below 2.5. Uncle Bob Martin says his teams achieve about 1.3-1.7 routinely.
I have heard him say he uses Java, Ruby, and Smalltalk.
This project is at 1.77 as I write this.  Many experienced developers
find that test driven development, along with low cyclomatic complexities
help to achieve robust applications more quickly than might otherwise
be achieved.  I have also found that to be the case, so much so that
I built this tool to supplement vim-go for my own use.

By, the way, a hat tip here to the well written github.com/fzipp/gocyclo,
which provides the code for determining cyclomatic complexity in Golang
code.  It took one import and exactly two extra lines of code to incorporate
this feature into goTestParser. Previous to that I had been running gocyclo
via a system call and grepping and awking the output before storing the
result away for later display. And, of course, the user would have to
install gocyclo on their systems themselves, adding to the barrier to
usefulness.

For now, I am just taking the go test -json output's word and we need to
realize that the results are approximate. I have not yet found even one
tool that gets this right.  Just don't be surprised if you think you
have written 33 tests and 35 get reported.

I have also seen instances where vim-go reports [SUCCESS] when there
actually were not tests run at all, or when tests were skipped with
no notification to the programmer.

My thinking is that if I am a consultant called in to work on a code base,
I do not want my tools delivering overly optimistic reports
to me.  To me, if a package has 100 tests, but 25 are not even being run,
I don't want the tools to report that as [SUCCESS]. If go test issues
messages via STDERR, I want to know that.  The fact that they were issued
on STDERR is info I want to be explicitly told.  Right now, vim-go reports
[SUCCESS] when there are skipped tests test files with no tests and
directories with no test files at all.  It does catch and relay the
message indicating there is no go.mod file in the directory though.

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
the quickfix list for you, but it does not force you to go to the failed
test every time. I find that often in my development work flow,
I don't need to go to the test at all, but want to peruse  and
fix the function that caused the failure, and I may already be there.

The RedBar/GreenBar/YellowBar message line lingers until any key is
pressed (I typically just hit the space bar).

In my personal set up, I have told vim-go to use the quickfix window
exclusively. I have <Leader>q toggle the quickfix window and <C-j> and <C-k>
navigate the quickfix list up and down. <LocalLeader>a takes me to the
Alternate file.

Life is good.....
