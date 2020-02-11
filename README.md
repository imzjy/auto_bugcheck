# auto_bugcheck
scan the windows dump files in a folder to get the bug check str

# Usage
```text
auto_bugcheck v0.0.3c 
  -c string
    	command issued to cdb debugger (default "!analyze -v;q")
  -d string
    	folder contains DMP files
  -f string
    	analyze specific dump file, ignore -d if flag set
  -p string
    	cdb file path (default "C:\\Program Files (x86)\\Windows Kits\\10\\Debuggers\\x64\\cdb.exe")
  -raw
    	raw cdb output
  -regex string
    	regular express to exact from cdb output
  -version
    	print version
```
