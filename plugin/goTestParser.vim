" Dave's GoLang Testing Stuff

" RunTest(), ProcessStandardOutput(), ColorBar(),
" and GoTestParser, a Golang program that runs go test -v -json
" and parses its output for use here,
" and these few mappings, work together to allow the running of
" RedBar/GreenBar Golang tests from within Vim.
"
" I prefer RedBar/GreenBar Test Driven Development, and missed my
" comfortable and pleasant development cycle when I moved to Golang.
"
" <LocalLeader>t runs 'go test ' and provides RedBar/GreenBar indications
" of failure or success.  It also provides useful summary info:
"   Tests Run
"   Tests Passed
"   Tests Skipped
"   Test Failed
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
    let g:packageDir = shellescape(expand('%:p:h'))
    "-count=1 ensures uncached results.  This is optional.
    " let s:cmdLine = 'go test -v -count=1 ' . p . ' | goTestParser '
    "
    " This script, goTestParser.vim, lives in goTestParser/plugin
    " The go source code lives in goTestParser/go.  The install.shell
    " script builds and moves the goTestParser binary into
    " goTestParser/go/bin.  So, from here ../go/bin/goTestParser should be a
    " reliable path to the gotTestParser binary.
    "
    "
    " TODO: I'd like to use ../go/bin/goTestParser but don't know
    " how to do it here....
    let g:goTestParserBinary="${HOME}/.config/nvim/plugged/goTestParser/bin/goTestParser"
    let l:oneSpace=" "
    let s:cmdLine=g:goTestParserBinary . oneSpace . g:packageDir
    if a:toScreen == v:true
      echon system(s:cmdLine)
    else
      let out = system(s:cmdLine)
     call s:ProcessStdOutput(out)
    endif
    let ch = getchar()
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

  if l:json_object.counts.fails > 0 || l:json_object.counts.skips > 0
    call setqflist(l:json_object.qflist,'r')
  endif

  call go#color_bar#DoColorBar(l:json_object.barmessage.color,
        \ l:json_object.barmessage.message)

  " Encode l:pgmdata into JSON and write it out for our inspection
  let l:tmp_json_object = json_encode(l:json_object)
  let l:logPath = expand('%:p:h') . '/gotestlog.json'
  call writefile(split(l:tmp_json_object,'\n'), l:logPath, 'b')
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
