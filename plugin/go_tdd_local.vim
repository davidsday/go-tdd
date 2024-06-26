
" NOTE: These are the user settable variables and maps
" for go-tdd.
" If you need to change any of these, please copy this file
" to your config (~/.vim or ~/.config/nvim)/after/ftplugin/go_tdd_local.vim
" and make the changes there, that way, your changes
" will not be overwritten every time you refresh this plugin from github.

" Vim has its own true/false, (the v: stands for Vim )
let g:go_tdd_debug=v:false
let g:gocyclo_ignore='vendor|testdata'
let g:go_list_type = 'quickfix'
let g:go_test_timeout = '10s'

nnoremap <LocalLeader>a  call go#alternate#Switch(<bang>0, 'edit')
nnoremap <LocalLeader>e <ESC>:e StdErr.txt<CR>
nnoremap <C-n> :cnext<CR>
nnoremap <C-p> :cprev<CR>

" I try to use <LocalLeader> in ftplugin types of situations
" This one runs GreenBar/RedBar tests
" nmap <silent> <LocalLeader>t <Plug>(RunGoGreenBarTests)
" This one opens a window below and pipes 'go test -v ' output so
" you can see the whole thing
nmap <silent> <LocalLeader>t <Plug>(RunGoGreenBarTests)
nmap <silent> <LocalLeader>v <Plug>(RunGoTestsVerbose)


