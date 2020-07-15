@echo off
chcp 65001 > nul
set c=
REM :1
set /p c=提交说明信息:
if "%c%"=="" (
    set c=%date:~3,10%
	REM goto 1
)

@echo on
git add .
git commit -m "%c%"
git push
::pause > nul
::exit