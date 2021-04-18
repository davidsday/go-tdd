 let l:pgmdata = {
   \ 'info': {
     \ 'host': '',
     \ 'user': '',
     \ 'begintime': '',
     \ 'endtime': '',
     \ 'commandline': ''
     \ },
   \ 'counts': {
     \ 'runs' : 0,
     \ 'pauses' : 0,
     \ 'continues' : 0,
     \ 'skips' : 0,
     \ 'passes' : 0,
     \ 'fails' : 0,
     \ 'outputs' : 0
     \ },
   \ 'firstfailedtest' : {
     \ 'fname' : '',
     \ 'tname' : '',
     \ 'lineno' : 0
     \ },
   \ 'elapsed' : 0.0,
   \ 'error' : {
     \ 'validjson' : v:true,
     \ 'notestfiles' : v:false,
     \ 'panic' : v:false,
     \ 'buildfailed' : v:false
     \ },
   \ 'qflist' : [],
   \ 'barmessage' : {
     \ 'color': '',
     \ 'message' : ''
     \ }
 \ }

 " Now we can build the QuickFix List
 let l:qfDict = {
      \ 'filename' : a:packageDir . '/' . l:parts[0],
      \ 'lnum'     : l:parts[1],
      \ 'col'      : 1,
      \ 'vcol'     : 1,
      \ 'pattern'  : a:json_line.Test,
      \ 'text'     : l:text
      \ }
