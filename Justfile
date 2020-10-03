#
# Ari Justfile
#

alias _build-mac := _build-macos
alias _build-win := _build-windows


@_default:
	just _term-wipe
	just --list


# Build compiled app
build target='':
	#!/bin/sh
	just _term-wipe
	if [[ '{{target}}' = '' ]]; then
		just _build-{{os()}}
	else
		just _build-{{target}}
	fi

@_build-macos:
	echo "Building macOS app"
	rm -rf bin/macos
	mkdir -p bin/macos
	GOOS=darwin GOARCH=amd64 go build -o bin/macos/sftpcmp main.go
	# bin/macos/sftpcmp -U mgardner -H 172.21.5.4 -P 2020 -d 1 . /Data/Archive/_FROM_/Marks-MBP/sftpcmp
	just distro macos bin/macos/sftpcmp

@_build-windows:
	echo "Building Windows app"
	rm -rf bin/windows
	mkdir -p bin/windows
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -o bin/win32/sftpcmp.exe main.go
	GOOS=windows GOARCH=amd64 go build -o bin/win64/sftpcmp.exe main.go
	# start bin/windows/sftpcmp -help
	just distro win32 bin/win32/sftpcmp.exe
	just distro win64 bin/win64/sftpcmp.exe


# Setup distrobution archive
distro arch file:
	#!/bin/sh
	path="$(dirname "{{file}}")"
	name="$(basename "{{file}}")"
	ver="$(just version)"
	# echo "path = ${path}"
	# echo "name = ${name}"
	# echo " ver = ${ver}"
	echo "cd ${path}"
	cd "${path}"
	echo "zip sftpcmp-v${ver}-{{arch}}.zip ${name}"
	zip "sftpcmp-v${ver}-{{arch}}.zip" "${name}"
	echo "SCP archive to http://webtools.zone/che/sftpcmp-v${ver}-{{arch}}.zip"
	scp "sftpcmp-v${ver}-{{arch}}.zip" drop2.webtools.zone:/vhost/zone/webtools/www/web/che/


# sftpcmp -U mgardner -H 172.21.5.4 -P 2020 -human -log -T . /Data/Accounting

# Run the program
run +args="*.log":
	just _term-wipe
	go run main.go {{args}}


# Wipes the terminal buffer for a clean start
_term-wipe:
	#!/bin/sh
	if [[ ${#VISUAL_STUDIO_CODE} -gt 0 ]]; then
		clear
	elif [[ ${KITTY_WINDOW_ID} -gt 0 ]] || [[ ${#TMUX} -gt 0 ]] || [[ "${TERM_PROGRAM}" = 'vscode' ]]; then
		printf '\033c'
	elif [[ "$(uname)" == 'Darwin' ]] || [[ "${TERM_PROGRAM}" = 'Apple_Terminal' ]] || [[ "${TERM_PROGRAM}" = 'iTerm.app' ]]; then
		osascript -e 'tell application "System Events" to keystroke "k" using command down'
	elif [[ -x "$(which tput)" ]]; then
		tput reset
	elif [[ -x "$(which reset)" ]]; then
		reset
	else
		clear
	fi


# Output the program version
@version:
	grep '^\tAppVersion' main.go | cut -d'"' -f2

