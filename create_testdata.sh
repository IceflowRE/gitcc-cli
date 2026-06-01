#!/bin/bash
set -ex

rm -rf "testdata/simpletag"
mkdir -p "testdata/simpletag"
cd "testdata/simpletag"

git init -q

export GIT_AUTHOR_NAME="Dummy"
export GIT_AUTHOR_EMAIL="example@example.com"
export GIT_COMMITTER_NAME="Dummy"
export GIT_COMMITTER_EMAIL="example@example.com"

# Branch A: Initial commits
touch file.txt
echo "A1" > file.txt && git add . && git commit --no-gpg-sign -m "[a1] Initial commit on main"
echo "A2" > file.txt && git add . && git commit --no-gpg-sign -m "[a2] Second commit on main"
echo "A3" > file.txt && git add . && git commit --no-gpg-sign -m "[a3] Point of divergence for branch B"

# Branch B: Create and add commits
git checkout -b branch-B
echo "B1" > file.txt && git add . && git commit --no-gpg-sign -m "[b1] First commit on branch B"
echo "B2" > file.txt && git add . && git commit --no-gpg-sign -m "[b2] Second commit on branch B"
echo "B3" > file.txt && git add . && git commit --no-gpg-sign -m "[b3] Final commit on branch B before merge"

# Branch A: Resume and Merge
git checkout master # or 'main' depending on your default
echo "A4" > file.txt && git add . && git commit --no-gpg-sign -m "[a4] Parallel work on main"

# A5 is the merge commit
git merge branch-B -X theirs --no-gpg-sign -m "[merge] Merge branch B into main"

# Final commit on A
echo "A6" > file.txt && git add . && git commit --no-gpg-sign -m "[a6] Final post-merge commit"

cd ../
rm -rf "simpletag_fail"
mkdir -p "simpletag_fail"
cd "simpletag_fail"

git init -q

# Branch A: Initial commits
touch file.txt
echo "A1" > file.txt && git add . && git commit --no-gpg-sign -m "[a1] Initial commit on main"
echo "A2" > file.txt && git add . && git commit --no-gpg-sign -m "[a2] Second commit on main"
echo "A3" > file.txt && git add . && git commit --no-gpg-sign -m "[a3]"

# Branch B: Create and add commits
git checkout -b branch-B
echo "B1" > file.txt && git add . && git commit --no-gpg-sign -m "[b1] First commit on branch B"
echo "B2" > file.txt && git add . && git commit --no-gpg-sign -m "|tag| Second commit on branch B"
echo "B3" > file.txt && git add . && git commit --no-gpg-sign -m "[B3] Final commit on branch B before merge"

# Branch A: Resume and Merge
git checkout master # or 'main' depending on your default
echo "A4" > file.txt && git add . && git commit --no-gpg-sign -m "Parallel work on main"

# A5 is the merge commit
git merge branch-B -X theirs --no-gpg-sign -m "[merge] Merge branch B into main"

# Final commit on A
echo "A6" > file.txt && git add . && git commit --no-gpg-sign -m "[a6] final post-merge commit"
