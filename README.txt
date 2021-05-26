This is a Vim plugin that I wrote in Golang for my own use to smooth my Golang
learning and development experience.

The plugin is 92+% in Golang and about 7% in Vimscript and the
tiniest bit of shell script thrown in.

It seeks to add two things to your Golang TDD development:

	RedBar/GreenBar/Refactor Test Driven Development style programming

	A marginally better 'vim-go go test' experience

The first is via the use of Vim's message line to provide red/green
bar fail/pass indications. I have also added a yellow bar message which
I use for messages about errors that are not caused by an actual
failing test.

Right now these yellow bar messages include
	[no tests found],
	[no tests to execute] - basically empty test file(s)
	[build failed],
	[received a panic],
	[invalid JSON message],
	and the receipt of any message on STDERR.

I capture stderr separately and process it by showing a snippet of the
message in a yellow bar.  If the STDERR message is longer than can be
shown in a one line yellow bar, I capture the entire message in StdErr.txt
in the package directory.

If go-tdd encounters non JSON lines on stdout, it issues a yellow
bar message and quits. I don't remember that happening in actual use.

I have added supplemental messages providing detail information in each
red or green bar.  They report the number of tests run, passed, failed,
and skipped, in addition to the elapsed time for running all the tests
as provided by go test. Test coverage is reported on Green Bars.

I have also added Average Cyclomatic Complexity to the Green Bars.
It is a metric I want to be aware of. I hear that several well known IDEs
start warning about Cyclomatic Complexity at 10.  I like to keep mine below
2.5.  This project is at 1.77 as I write this.

Many experienced developers find that test driven development, along with
low cyclomatic complexities help to achieve robust applications more quickly
than might otherwise be achieved.  I certainly have, so much so that I built
this tool to supplement vim-go for my own use. This project has ~90% test
coverage as I write this (basically everything but main() and a few
logFatal() type calls).

And, as an aside, I have had occasion to go back and retrofit (rewrite)
a few personal projects in this style recently, and I have been
surprised at the bugs that writing tests brings to the surface
even on projects which seem to be working fine.  I suspect that these
personal projects have only been exercised with my personal use patterns,
and I had of course corrected the bugs that come to light in my own use
scenario, so the apps seemed to be bug free to me.  Writing tests to
substantially cover 85-95% of their code further flushes out bugs I had
not found.  Humbling....  It has further confirmed to me that TDD is
a worthwhile style of development and especially when combined with
the use of cyclomatic complexity metrics to make sure your functions
stay small, simple, debuggable and readable.  I have found that
having the cyclomatic complexity metric more or less ever present in my
development efforts helps me preemptively keep things simple.  It is
easier to write a simple test for a simple function.

By, the way, a hat tip here to the well written github.com/fzipp/gocyclo,
which provides the code for determining cyclomatic complexity in Golang
code.  It is compiled directly into this plugin, so the user need not
do a thing, except work on the complexity of his/her code.

I am taking the go test -json output's word as to the number
of tests run, passed, failed, etc, but we need to realize that the results
are approximate. When you write a test with subtests, go test seems to count
the mother/father test in addition to all the subtests.  Thing is, the parent
test often doesn't actually do any testing itself and perhaps
shouldn't be counted. Just don't be surprised at counts that differ
slightly from your counts, if you use subtests.

go-tdd is designed to work alongside of vim-go, since, really,
vim-go is my most important golang development tool.

It does not interfere with vim-go in anyway that I am aware of.
In my setup, I have replaced vim-go's <Leader>t (<ESC>:GoTest<CR) with
<LocalLeader>t to activate go-tdd. If I desire to use vim-go's
:GoTest command, I call it just like that.

The second benefit listed above was a "marginally better go test
experience".  vim-go reports [SUCCESS] in directories with no test files
at all or where there are test files but they are empty, or where one, or
many tests are skipped. I am not an old Golang hand, but this does not
strike me as "[SUCCESS]".  Especially if I am looking at a code base that
is new to me, I don't want my tools reporting these situations as "[SUCCESS]".
So in go-tdd I have incorporated a "Yellow Bar", message for situations
which are not directly due to a failing test but which the developer
should be aware of, thus providing that "marginally better go test
experience" I mentioned above.

To accomplish this, go-tdd provides its own go test parser, written
in golang, somewhat simpler than vim-go's and synchronous instead of
asynchronous, which parses the 'go test -v -json' output and in turn,
provides a further processed JSON structure which details for a small
Vimscript what message, and in what color to deliver.  It also provides to
Vim a quickfix list of test failures and/or skipped tests which Vim loads
for your use.

go-tdd's synchronous invocation of 'go test -v -json -cover' has not
really been noticeable for me. I rarely see go test take more than a few
hundredths of a second to complete even hundreds of tests. Most reported
times are in the thousandths of seconds and go-tdd, written almost
entirely in Golang itself, does its job and returns control to Vim quickly.

In this style of development, the RedBar/GreenBar (and YellowBar)s are the
primary layer of communication with the developer, so go-tdd loads
the quickfix list for you, but it does not force you to go to the failed
test. I find that often in my development work flow, I sometimes don't need
to go to the test at all, but instead want to peruse  and
fix the function that caused the failure, and I may already be there.

The RedBar/GreenBar/YellowBar message line lingers until any key is
pressed (I typically just hit the space bar).

<Space> (or any other key, for that matter) dismisses the Green/Red/Yellow
bars.

Configuration:

If you are using vim-plug:
	Plug 'davidsday/go-tdd', {'for': [ 'go' ], 'branch': 'master', 'do': './install.sh'}
	Plug 'Valloric/ListToggle'
		If you use this, <Leader>q toggles the quickfix window open and closed


go-tdd comes with a small number of user alterable configuration
settings and maps.

I have gathered them all into go-tdd/plugin/go_tdd_local.vim and it looks
like this:
	" the actual file does not have these comments
	let g:go_tdd_debug=0 " set to 1 if you want to log JSON sent to vim
	let g:gocyclo_ignore="'vendor|testdata'"  "regex for gocyclo to ignore
	let g:go_list_type = 'quickfix'  " my pref is to only use quickfix

	" toggle between mygofile.go and mygofile_test.go
	nnoremap <LocalLeader>a  call go#alternate#Switch(<bang>0, 'edit')
	" navigate quickfix list down and up
	nnoremap <C-j> :cnext<CR>
	nnoremap <C-k> :cprev<CR>

	" run go-tdd on package in current dir
	nmap <silent> <LocalLeader>t <Plug>(RunGoGreenBarTests)
	" run verbosely - just shows go test -v -json -cover verbose output
	nmap <silent> <LocalLeader>v <Plug>(RunGoTestsVerbose)

If you need to make changes, please copy this file to
your config dir

(~/.vim/ or ~/.config/nvim/)after/ftplugin/go_tdd_local.vim

and make any changes there.  That way your changes will not be overwritten
every time you refresh this plugin from github.

Life is good.....

Dave
