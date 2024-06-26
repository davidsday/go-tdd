
" RunTest(), ProcessStandardOutput(), ColorBar(),
" and go-tdd, a Golang program that runs go test -v -json -cover
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
if exists('g:did_gotdd_ftplugin')
  finish
endif
let g:did_gotdd_ftplugin = 1

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
let s:plugin_dir = shellescape(fnameescape(expand('<sfile>:p:h:h')))


" toScreen needs to either be v:true or v:false
" If toScreen is v:true, stdout goes to the terminal
" If not, stdout is captured in a variable 'out'
" and passed to ProcessStdOutput(out) to be processed
" with the results displayed on the message line.
function! s:RunTest(toScreen)
    " Turn off Vim-Go's automatic type display until we are done here
    let g:go_echo_go_info = 0
    echon 'Testing...'
    let l:go_tdd_executable_filename='go-tdd'
    if has('win32')
        let l:go_tdd_executable_filename='go-tdd.exe'
    endif
    "shellescape(expand('%:p:h')) gives path to this docs directory
    "(package dir)
    let l:package_dir=shellescape(fnameescape(expand('%:p:h')))
    " Ensure Vim's working directory is the same as the file we are editing
    " Without this, sometimes, when opening opening a file found by FZF
    " (<Leader>f), Vim's working directory stays at the directory we just
    " left.  So don't delete it.
    chdir %:p:h
    " let l:go_tdd_binary=s:plugin_dir . '/bin/go-tdd'
    let l:go_tdd_binary=go#util#Join(s:plugin_dir, 'bin',l:go_tdd_executable_filename)
    let l:oneSpace=' '
    let l:screencolumns=string(&columns - 1)
    if !exists('g:go_test_timeout')
      let g:go_test_timeout='10s'
    endif
    if !exists('g:go_tdd_debug')
      let g:go_tdd_debug=v:false
    endif
    if !exists('g:gocyclo_ignore') || g:gocyclo_ignore ==# ''
      let g:gocyclo_ignore='vendor|testdata'
    endif

    let l:arg_dict={}
    if has('win32')
        let l:arg_dict['package_dir']=trim(l:package_dir,"'")
        let l:arg_dict['plugin_dir']=trim(s:plugin_dir,"'")
    else
        let l:arg_dict['plugin_dir']=trim(s:plugin_dir,"'", 0)
        let l:arg_dict['package_dir']=trim(l:package_dir,"'", 0)
    endif
    let l:arg_dict['screen_columns']=l:screencolumns
    let l:arg_dict['gocyclo_ignore']=g:gocyclo_ignore
    let l:arg_dict['go_tdd_debug']=g:go_tdd_debug
    let l:arg_dict['timeout']=g:go_test_timeout

    let l:json_args=shellescape(json_encode(l:arg_dict))
    let l:cmdLine=l:go_tdd_binary
    let l:cmdLine.= oneSpace . l:json_args


    if a:toScreen == v:true
      echon system(l:cmdLine)
    else
      let l:out = system(l:cmdLine)
     call s:ProcessStdOutput(l:out)
    endif
    let l:ch = getchar()
    " DEBUG: do I need to OS neutralize this?  I'm thinking probably
    " But it seemed to be working OK the last time I tried on Windows
    echon "\r\n\n"
    redraw
    " Turn Vim-Go's automatic type display back on
    let g:go_echo_go_info = 1
endfunction


" This function operates on the JSON object passed back to Vim from
" the go-tdd binary, creating the quickfixlist if appropriate and
" calling go#color_bar#DoColorBar() to display a Red/Green/Yellow Bar
" message as provided by go-tdd.
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
