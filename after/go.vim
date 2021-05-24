

let g:go_tdd_debug=0
let g:gocyclo_ignore="'vendor|testdata'"
let g:go_list_type = 'quickfix'

nnoremap <LocalLeader>a  call go#alternate#Switch(<bang>0, 'edit')
nnoremap <C-j> :cnext<CR>
nnoremap <C-k> :cprev<CR>

nmap <silent> <LocalLeader>t <Plug>(RunGoGreenBarTests)
nmap <silent> <LocalLeader>v <Plug>(RunGoTestsVerbose)


