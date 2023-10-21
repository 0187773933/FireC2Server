#!/bin/bash

FILE="./v1/server/html/login.html"
NEW_LINE="\t\t\t\t\t<form id=\"form-login\" action=\"/ok???\" onSubmit=\"return on_submit();\" method=\"post\">"
awk -v n=17 -v s="$NEW_LINE" '(NR == n) {$0 = s} 1' $FILE > temp.txt
mv temp.txt $FILE

function is_int() { return $(test "$@" -eq "$@" > /dev/null 2>&1); }
ssh-add -D
git init
git config --global --unset user.name
git config --global --unset user.email
git config user.name "0187773933"
git config user.email "collincerbus@student.olympic.edu"
ssh-add -k /Users/morpheous/.ssh/githubWinStitch

LastCommit=$(git log -1 --pretty="%B" | xargs)
# https://stackoverflow.com/a/3626205
if $(is_int "${LastCommit}");
	 then
	 NextCommitNumber=$((LastCommit+1))
else
	echo "Not an integer Resetting"
	NextCommitNumber=1
fi
git add .
if [ -n "$1" ]; then
	git commit -m "$1"
	git tag v1.0.$1
else
	git commit -m "$NextCommitNumber"
	git tag v1.0.$NextCommitNumber
fi
git remote add origin git@github.com:0187773933/FireC2Server.git

# https://proxy.golang.org/github.com/0187773933/FireC2Server/@v/v1.0.2.info
git push origin --tags
git push origin master