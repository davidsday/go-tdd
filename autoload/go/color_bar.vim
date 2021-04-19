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
  echon a:msg . repeat(' ',&columns - (length) )
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

  if a:runCount >1
    let l:txtmsg = ' tests run,'
  else
    let l:txtmsg = ' test run,'
  endif
  let l:barMessage = a:runCount . l:txtmsg
  let l:barMessage .= ' ' . a:passCount . ' passed,'
  if a:skipCount > 0
    let l:barMessage .= ' ' . a:skipCount . ' skipped,'
  endif
  " Our failCount will potentially include skipped files
  " So the true error count is failCount - skipCount
  if a:failCount - a:skipCount > 0
    if a:failCount == 1
      " Singular
      let l:fails = 'failure'
      let l:barMessage .=  ' ' . a:failCount .' ' . fails . ','
    else
      " or plural...
      let l:fails = 'failures'
      let l:barMessage .= ' ' . a:failCount .' ' . fails . ','
    endif
    if len(a:firstFailTestFname) > 0
      if a:failCount == 1
        " Singular
        let l:barMessage .= ' in file ' . a:firstFailTestFname . ','
      else
        let l:barMessage .= ' 1st in file ' . a:firstFailTestFname . ','
      endif
    endif
    if len(string(a:firstFailTestLineNo)) > 0
      let l:barMessage .= ' line nr ' . a:firstFailTestLineNo . ','
    endif
  endif
  if len(string(a:runTime)) > 0
    let l:barMessage .= ' in ' . string(a:runTime) . 's'
  endif
  return l:barMessage
endfunction abort
