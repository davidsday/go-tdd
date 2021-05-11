" Dave's GoLang Testing Stuff

function! go#color_bar#DoColorBar(color, msg)
  let l:length = strlen(a:msg)
  redraw

  if a:color ==# 'yellow'
    hi YellowBar ctermfg=black ctermbg=yellow guibg=yellow guifg=black
    echohl YellowBar
  elseif a:color ==# 'red'
    hi RedBar ctermfg=white ctermbg=red guibg=red guifg=white
    echohl RedBar
  else
    hi GreenBar ctermfg=white ctermbg=green guibg=#719e07 guifg=black
    echohl GreenBar
  endif
  if l:length >= &columns
    let l:length = &columns - 1
  endif
  echon a:msg . repeat(' ',&columns - (l:length) )
  echohl None
endfunction
