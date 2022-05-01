

#prune merged git branches locally
git remote prune origin

local_branches=$(git branch --merged | grep -v 'master$' | grep -v "master")
if [ -n "$local_branches" ]; then
  echo "$local_branches"
fi


git branch -d `git branch --merged | grep -v 'master$' | sed 's/origin\///g' | tr -d '\n'`


$SHELL