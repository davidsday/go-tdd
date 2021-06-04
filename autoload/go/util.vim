" don't spam the user when Vim is started in Vi compatibility mode
let s:cpo_save = &cpo
set cpo&vim

" PathSep returns the appropriate OS specific path separator.
function! go#util#PathSep() abort
  if go#util#IsWin()
    return '\'
  endif
  return '/'
endfunction

" PathListSep returns the appropriate OS specific path list separator.
function! go#util#PathListSep() abort
  if go#util#IsWin()
    return ";"
  endif
  return ":"
endfunction

" LineEnding returns the correct line ending, based on the current fileformat
function! go#util#LineEnding() abort
  if &fileformat == 'dos'
    return "\r\n"
  elseif &fileformat == 'mac'
    return "\r"
  endif

  return "\n"
endfunction

" Join joins any number of path elements into a single path, adding a
" Separator if necessary and returns the result
function! go#util#Join(...) abort
  return join(a:000, go#util#PathSep())
endfunction

" IsWin returns 1 if current OS is Windows or 0 otherwise
" Note that has('win32') is always 1 when has('win64') is 1, so has('win32') is enough.
function! go#util#IsWin() abort
  return has('win32')
endfunction

" IsMac returns 1 if current OS is macOS or 0 otherwise.
function! go#util#IsMac() abort
  return has('mac') ||
        \ has('macunix') ||
        \ has('gui_macvim') ||
        \ go#util#Exec(['uname'])[0] =~? '^darwin'
endfunction

" restore Vi compatibility settings
let &cpo = s:cpo_save
unlet s:cpo_save

" vim: sw=2 ts=2 et
