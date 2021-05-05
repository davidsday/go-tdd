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

function! go#color_bar#BuildBarMessage(
      \ runCount,skipCount,failCount,passCount,
      \ runTime,firstFailTestFname,
      \ firstFailTestLineNo
      \ ) abort
  " As we build this message (barMessage), each line is responsible
  " for providing any blank spaces preceding it, and must NOT leave
  " any blank spaces after itself. Life is simpler that way.

  let l:oneSpace=' '
  let l:commaSpace=', '


  if a:runCount > 1
    " plural
    let l:txtmsg = 'tests run'
  else
    " singular
    let l:txtmsg = 'test run'
  endif
  let l:barMessage = a:runCount . l:oneSpace . l:txtmsg
  let l:barMessage .= l:commaSpace . a:passCount . l:oneSpace . 'passed'
  if a:skipCount > 0
    let l:barMessage .= commaSpace . a:skipCount . oneSpace . 'skipped'
  endif
  " Our failCount will potentially include skipped files
  " So the true error count is failCount - skipCount
  if a:failCount - a:skipCount > 0
    if a:failCount == 1
      " Singular
      let l:fails = 'failure'
      let l:barMessage .=  l:commaSpace . a:failCount . l:oneSpace . l:fails
    else
      " or plural...
      let l:fails = 'failures'
      let l:barMessage .= l:commaSpace . a:failCount . l:oneSpace . l:fails
    endif
    if len(a:firstFailTestFname) > 0
      if a:failCount == 1
        " Singular
        let l:barMessage .= l:commaSpace . 'in file' . l:oneSpace . a:firstFailTestFname
      else
        let l:barMessage .= l:commaSpace . '1st in file' . l:oneSpace . a:firstFailTestFname
      endif
    endif
    if len(string(a:firstFailTestLineNo)) > 0
      let l:barMessage .= l:commaSpace . 'line nr' . l:oneSpace . a:firstFailTestLineNo
    endif
  endif
  if len(string(a:runTime)) > 0
    let l:barMessage .=  l:commaSpace . 'in' . oneSpace . string(a:runTime) . 's'
  endif
  return l:barMessage
endfunction abort
