" Dave's GoLang Testing Stuff

" RunTest(), ProcessStandardOutput(), ColorBar(),
" and GoTestParser, a Golang program that runs go test -v -json -color
" and parses its output for use here, and these few mappings, work
" together to allow the running of RedBar/GreenBar Golang tests from within
" Vim.
"
" I prefer RedBar/GreenBar Test Driven Development, and missed my
" comfortable and pleasant development cycle when I moved to Golang.
"
" <LocalLeader>t runs 'go test -v -json -cover ' and provides RedBar/GreenBar
" indications of failure or success.
"
" It also provides useful summary info:
"   Tests Run
"   Tests Passed
"   Tests Skipped
"   Tests Failed
"   Test  Coverage
"   Cyclomatic Complexity
"   And the file and line number of the first failure or the cause of
"   the first skip.
"
" There is also a yellow bar indication for [No Tests Found], [Build Failed]
" [Invalid JSON] and [Received a Panic] sorts of problem indications
" for non test failures.
"
" <LocalLeader>v runs 'go test -v ' verbosely to the screen, allowing you
" to see the entire output, and inspect more closely.
"
if exists('g:did_gotst_ftplugin')
  finish
endif
let g:did_gotst_ftplugin = 1

" NOTE: Important to understand, this will not work
" inside a function.  Must be done like this, external
" to any function, prior to them being called.  Once
" you've got the result in a script var, you can use
" it in the functions as normal, I spent days.....
" Here the script is go-tdd/plugin/go-tdd.vim
" the expand removes plugin/go-tdd.vim, leaving us
" with pathtoplugin, to which I can add '/bin/go-tdd'
" and have the path to our binary, where ever the plugin
" manager might have put it.
let s:plugin_dir = expand('<sfile>:p:h:h')

" toScreen needs to either be v:true or v:false
" If toScreen is v:true, stdout goes to the terminal
" If not, stdout is captured in a variable 'out'
" and passed to ProcessStdOutput(out) to be processed
" with the results displayed on the message line.
function! s:RunTest(toScreen)
    " Turn off Vim-Go's automatic type display until we are done here
    let g:go_echo_go_info = 0
    echon 'Testing...'
    "shellescape(expand('%:p:h')) gives path to this docs directory
    "(package dir)
    let l:packageDir = shellescape(expand('%:p:h'))
    " Ensure Vim's working directory is the same as the file we are editing
    " Without this, sometimes, when opening opening a file found by FZF
    " (<Leader>f), Vim's working directory stays at the directory we just
    " left.  So don't delete it.
    chdir %:p:h
    " let l:goTestParserBinary="${HOME}/.config/nvim/plugged/go-tdd/bin/go-tdd"
    let l:goTestParserBinary=s:plugin_dir . '/bin/go-tdd'
    let l:oneSpace=' '
    let l:screencolumns=string(&columns - 1)

    let l:cmdLine=l:goTestParserBinary . oneSpace . l:packageDir . oneSpace . l:screencolumns
    if a:toScreen == v:true
      echon system(l:cmdLine)
    else
      let l:out = system(l:cmdLine)
     call s:ProcessStdOutput(l:out)
    endif
    let l:ch = getchar()
    echon "\r\n\n"
    redraw
    " Turn Vim-Go's automatic type display back on
    let g:go_echo_go_info = 1
endfunction


" The stdout arg is a pointer to the stdout stream, indicating that
" this function operates on that stream.
" This function operates on the captured stdout stream from the tests
" and makes decisions about what to show in the message bar, green or red
" or yellow background, and what messages to show

function! s:ProcessStdOutput(stdout) abort
  let l:packageDir = shellescape(expand('%:p:h'))
  let l:json_object = json_decode(a:stdout)
  if l:json_object.quickfixlist != []
    call setqflist(l:json_object.quickfixlist,'r')
  endif
  call go#color_bar#DoColorBar(l:json_object.color,
        \ l:json_object.message)
endfunction

noremap <unique> <Plug>(RunGoTestsVerbose) :call <SID>RunTest('True')<CR>
noremap <unique> <Plug>(RunGoGreenBarTests) :call <SID>RunTest('False')<CR>

" I try to use <LocalLeader> in ftplugin types of situations
" This one runs GreenBar/RedBar tests
nmap <silent> <LocalLeader>t <Plug>(RunGoGreenBarTests)
" This one opens a window below and pipes 'go test -v ' output so
" you can see the whole thing
nmap <silent> <LocalLeader>v <Plug>(RunGoTestsVerbose)
""nnoremap <silent> <LocalLeader>v :!clear;./runTests<cr>

""=============================================================================
"" GoTest
""=============================================================================
" go test this project
" nnoremap <silent><LocalLeader>p :GoTest ./... <CR>
" go test this package
" nnoremap <silent><LocalLeader>m :GoTest<CR>
" go test this file
" nnoremap <silent><LocalLeader>f :GoTest expand('%:p')<CR>
" go test this function
" nnoremap <silent><LocalLeader>u :GoTestFunc file<CR>
