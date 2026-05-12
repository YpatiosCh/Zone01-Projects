
## Authors

- ychaniot (Ypatios Chaniotakos)

Project [link](https://platform.zone01.gr/git/ychaniot/git).

## _The ultimate Git Guide_    


[![Build Status](image/zone-logo.jpg)

[![Build Status](image/git.png)

This project is designed to introduce you to the world of version control and collaboration using Git. Git is a powerful and widely used tool for tracking changes in your projects, collaborating with others, and ensuring the integrity of your code.

- Start from the bascis
- Explore more advanced topis
- Equip yourself with the essential knowledge and practices for effective version control and collaboration
--------------
## Setting up Git 

- Install Git on your local machine by following the instructions for your operating system on the official Git [website.](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- Configure Git with your username and email address.
    ```
    git config user.name "Your Name"
    git config user.email "your.email@example.com"
    ```
--------------
## Git commits to commit
- Within the work directory, establish a subdirectory named hello. Inside this directory, generate a file titled hello.sh and input the following content:
```
echo "Hello, World"
```

- Initialize the git repository in the hello directory:
    ```
    git init           //initializes an empty local repository
    ```

- Check the status and act accordingly with the output of the executed command.
    ```
    git status         //checks the status and provides info of unstaged files
    git add hello.sh   //staged the hello.sh file so that it can be tracked
    ```
-  Change the hello.sh content to the following:
    ```
    #!/bin/bash
    echo "Hello, $1"
    ```

- Stage the changed file and commit the changes, the working tree should be clean.
    1. Stage specific file:
    ```
    git add <filename>           
    ```
    or stage everything that was changed:
    ```
    git add .
    ```
    2. Commit
    ```
    git commit -m "first commit"  
    ```
    this comment usually contains the change that was made and/or what stage of the program are 

    3. Check status
    ```
    git status                   
    ```
---------------------------
-  Modify the hello.sh file to include comments and stage it.
```
 #!/bin/bash
 
 # Default is "World"
 name=${1:-"World"}
 echo "Hello, $name"
```
- Make two separate commits:
    1. The first commit should be for the comment in line 3:
    ```
    git add hello.sh
	git commit -m "comment in line 3"
    ```

    2. The second commit should include changes made to lines 4 and 5:
    ```
	git add hello.sh
	git commit -m "changes made in line 4 and 5"
    ```
-------------------------------
## History

- Show the history of the working directory.
    ```
    git log      //This will show a list of commits in the current branch, including the commit hash, author, date, and commit message
    ```

- Show One-Line History for a condensed view showing only commit hashes and messages.
    ```
	git log --oneline     //This displays each commit in a single line, with the commit hash and the message—great for an overview of changes
    ```

- Controlled Entries:

    You need to customize the log output by specifying the number of entries or a time range. Customize it to display the last 2 entries and to view the commits made within the last 5 minutes.
    ```
	git log -n 2 --since="5 minutes ago"     //the log output will only show up to two commits that were made within the last 5 minutes
    ```

- Personalized Format:

    Show logs in a personalized format, including the commit hash, date, message, branch information, and author name, resembling * e4e3645 2023-06-10 | Added a comment (HEAD -> main) [John Doe]
    ```
	git log --pretty=format:"* %h %ad | %s (%d) [%an]" --date=short
    ```
        
| Command| Explanation
| -----| --------------------------------------------------------- |        
| *%h |  Shows the commit hash (abbreviated) with a leading asterisk. |
| %ad |    Shows the commit date. |
| %s    |  Displays the commit message. |
| %d    |  Displays branch information and any references (e.g., HEAD -> main). |
| %an   |  Shows the author’s name. |
| --date=short |  Formats the date in a short format (YYYY-MM-DD). |

--------------
## Check it out 

- Restore First Snapshot:
    
    Revert the working tree to its initial state, as captured in the first snapshot, and then print the content of the hello.sh file.

1. Find the First Commit
	```
	git rev-list --max-parents=0 HEAD    //This will print the hash of the very first commit
	```
	
2. Reset the Working Tree to the First Commit:
	
	Option 1: Soft Reset (keeps changes in the working directory but moves HEAD)
	```
	git reset --soft <first_commit_hash>
	```
		
	Option 2: Hard Reset (reverts everything, including staging and working directory)
	```
	git reset --hard <first_commit_hash>
	```
		
3. Print the Contents of hello.sh
	```
	cat hello.sh
    ```
--------------
- Restore Second Recent Snapshot:
Revert the working tree to the second most recent snapshot and print the content of the hello.sh file:

1. Identify the Second Most Recent Commit
	```
	git log -n 2 --pretty=format:"%h"
	```
		
2. Reset to the Second Most Recent Commit
	```
	git reset --soft <second_recent_commit_hash>      //if you want to keep changes and only move head
	```
    or
	```
	git reset --hard <second_recent_commit_hash>    //revert everything, including staging and working directory
	```
	    
3. Print the Content of hello.sh
	```
	cat hello.sh
    ```
--------------

- Return to Latest Version:
Ensure that the working directory reflects the latest version of the hello.sh file present in the main branch, without referring to specific commit hashes:

1. Switch to the Latest Version in the main Branch
    ```	
	git checkout main
	```
		
2. Reset hello.sh to Match the Latest Commit on main:
	```
	git checkout main --hello.sh
	```
		
3. Print the Content of hello.sh
	```
	cat hello.sh
	```
	Alternative: Reset the Entire Working Directory (if Needed) (Not recommended for this project)
	```
	git reset --hard origin/main 
	```
	(This approach will revert any local changes and align your working directory fully with the latest state of the main branch, ensuring hello.sh and all other files are up to date. Note: it will also discard any uncommited changes)       

--------------
## TAG me

- Referencing Current Version:

    Tag the current version of the repository as v1.
    ```
    git tag v1  //This command will create a lightweight tag named v1 pointing to the latest commit on the current branch (usually main)
	```
--------------	    
- Tagging Previous Version:

    Tag the version immediately prior to the current version as v1-beta, without relying on commit hashes to navigate through the history.

1. Identify the Previous Commit
	```
	git tag v1-beta HEAD^  //This command tags the commit right before the latest one as v1-beta
	```
2. Verify the Tag
	```
	git show v1-beta    //This will display details about the v1-beta tag, including the commit message, hash, and date
	```
--------------		
- Navigating Tagged Versions:

    Move back and forth between the two tagged versions, v1 and v1-beta.

1. Checkout v1-beta
	```
	git checkout v1-beta
	```
		
2. Checkout v1
	```
	git checkout v1
	```
		
3. Return to the Latest Commit on main (Optional)
	```
	git checkout main
	```
		
	Note: When you’re on a tag, Git enters a "detached HEAD" state, meaning you're not on a branch. Any new commits made in this state won’t belong to a branch unless you create a new branch from that tag or return to an existing branch.
------------------------------
- Listing Tags:

    Display a list of all tags present in the repository to verify successful tagging.
    ```
	git tag //This command will show a simple list of all tags, including v1 and v1-beta if they were created successfully
	```
--------------  
## Changed Your mind?

- Reverting Changes:

    Modify the latest version of the file with unwanted comments, then revert it back to its original state before staging using a Git command:
```
#!/bin/bash
# This is a bad comment. We want to revert it.
name=${1:-"World"}
echo "Hello, $name"
```

```
git checkout -- hello.sh //revert to previous state if file was not staged
```
--------------

- Staging and Cleaning:

    Introduce unwanted changes to the file, stage them, then clean the staging area to discard the changes:
```
#!/bin/bash
# This is an unwanted but staged comment
name=${1:-"World"}
echo "Hello, $name"
```

1. Make Unwanted Changes to the File 
2. Stage the Unwanted Changes
	```
	git add hello.sh
	```
3. Remove the Changes from the Staging Area
	```
	git restore --staged hello.sh //This command removes the file from the staging area but keeps the unwanted changes in the working directory
    ```
    or
    ```
	git restore hello.sh //return the file completely to its last committed state
	```
--------------
- Committing and Reverting:

    Add the following unwanted changes again, stage the file, commit the changes, then revert them back to their original state.
```
#!/bin/bash
# This is an unwanted but committed change
name=${1:-"World"}
echo "Hello, $name"
```

1. Add the Unwanted Changes
2. Stage the Changes
	```
	git add hello.sh
	```
3. Commit the Unwanted Changes
	```
	git commit -m "added unwanted comment to hello.sh"
	```
4. Revert the commit:
	```
	git revert HEAD 
	```
Note: This command will create a new commit that undoes the changes from the previous commit (the unwanted changes), preserving your commit history.
Git will open an editor to allow you to edit the commit message. You can leave it as the default (something like "Revert 'Unwanted changes added to hello.sh'") or modify it if needed.

--------------
- Tagging and Removing Commits:

    Tag the latest commit with oops, then remove commits made after the v1 version. Ensure that the HEAD points to v1.

1. Tag the Latest Commit as oops
	```
	git tag oops
	```
2. Remove Commits After v1 and Point HEAD to v1
	```
	git reset --hard v1 //This command moves HEAD to the commit tagged as v1 and removes any later commits from the history
	```
Note: Using --hard will also reset the working directory and staging area to match the v1 commit. Be careful with this command if you have any uncommitted changes, as it will discard them.

3. Verify Your Tags and HEAD Position
	```
	git log --oneline --decorate
	```
	or
	```
	git tag //To see all tags (including v1 and oops)
    ```
 --------------    
- Displaying Logs with Deleted Commits:

    Show the logs with the deleted commits displayed, particularly focusing on the commit tagged oops.

1. Use git reflog to Display Logs with Deleted Commits
	```
	git reflog  
    ```
	This will show a list of actions (like commits, resets, and checkouts) along with their respective commit hashes. The commit tagged as oops should be listed in the reflog with its hash
		
2. Show Details of the oops Commit
	```
	git show oops 
	```
	This will display detailed information about the commit, including changes made, commit message, author, and date.
--------------		
- Cleaning Unreferenced Commits:

    Ensure that unreferenced commits are deleted from the history, meaning there should be no logs for these deleted commits.

1. Clear the Reflog for All Branches
	```
	git reflog expire --expire=now --all 
	```
    _--expire=now_: Forces all reflog entries to expire immediately.

    _--all_: Applies this operation to all branches and refs, clearing the reflog history.)

2. Run Garbage Collection with Pruning
	```
	git gc --prune=now 
	```
	The git gc command cleans up unnecessary files and optimizes the local repository. --prune=now: This option removes all unreachable commits right away, rather than waiting for the default grace period (usually 30 days))
		
3. Verify the Deletion of Unreferenced Commits	
	```
	git reflog
	```
	You should no longer see entries for the deleted commits.
	```
	git log 
	```
	This will only show the current history, without the removed commits.
		
--------------
- Author Information:

    Add an author comment to the file and commit the changes.
```
#!/bin/bash

# Default is World
# Author: Jim Weirich
name=${1:-"World"}

echo "Hello, $name"
```

1. Make the changes, stage and commit
	```
	git add hello.sh
	git commit -m "Author comment"
	```
--------------
- Oops the author email was forgotten, update the file to include the email without making a new commit, but include the change in the last commit.

1. Change the file adding the author's email:
```
#!/bin/bash
 
# Default is World
# Author: Jim Weirich <jim@example.com>
name=${1:-"World"}

echo "Hello, $name"
```

2. Stage the Updated File
	```
	git add hello.sh
	```
		
3. Amend the Last Commit
	```
	git commit --amend --no-edit (--amend: Modifies the last commit.
	```
    _--no-edit_: Keeps the original commit message
    
4. Verify the Updated Commit
	```
	git show
	```	
--------------
## Move it

- Moving hello.sh:

    Using Git commands, move the program hello.sh into a lib/ directory, and then commit the move.

1. Create the lib/ Directory (if it doesn’t already exist)
	```
	mkdir -p lib
	```
		
2. Move hello.sh to the lib/ Directory
	```
	git mv hello.sh lib/hello.sh
	```
		
3. Commit the Move
	```
	git commit -m "Move hello.sh into the lib/ directory"
	```
--------------
- Create a Makefile in the root directory of the repository with the provided content and commit it to the repository.
```
TARGET="lib/hello.sh"
 
run:
 	bash ${TARGET}
```

1. Create the Makefile
    ```
	echo -e 'TARGET="lib/hello.sh"\n\nrun:\n\tbash ${TARGET}' > Makefile
	```
		
2. Stage the Makefile
	```
	git add Makefile
	```
		
3. Commit the Makefile
	```
	git commit -m "Add Makefile to run lib/hello.sh"
	```
(Note: if you want to run the Makefile use the command: make run )

--------------
## blobs, trees and commits

- Navigate to the .git/ directory in your project and examine its contents. You will have to explain the purpose of each subdirectory, including objects/, config, refs, and HEAD in the audit.

_objects/_

		1. Purpose: This directory contains all the actual data (like commits, trees, and blobs) for your repository in a compressed format.
		2. Types of Objects:
			Blobs: Store the contents of each file.
			Trees: Represent directory structures and contain references to blobs and other trees.
			Commits: Store metadata about each commit, such as author, message, and pointers to parent commits and associated trees.
_refs/_

		1.Purpose: Stores references to branches and tags, helping Git manage different versions of the repository.
		2.Subdirectories:
			refs/heads/: Contains files for each local branch, each file pointing to the latest commit in that branch.
			refs/tags/: Stores pointers to specific commit hashes associated with tags, allowing you to mark significant points in history.
			refs/remotes/: Contains files for each remote-tracking branch, representing branches from remote repositories.
_config_

		1. Purpose: The config file holds repository-specific configurations, such as user information, remote repository URLs, and settings for pushing/pulling.
		2 .Common Configurations:
			[user]: Stores the username and email for commits.
			[remote "origin"]: Defines the URL for the remote repository.
			[branch "main"]: Configuration specific to the main branch, like remote tracking information.

_HEAD_

		1. Purpose: The HEAD file is a pointer to the current branch’s latest commit. It tells Git where you are in the commit history.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
- Latest Object Hash:

    Find the latest object hash within the .git/objects/ directory using Git commands and print the type and content of this object using Git commands.

1. Find the Latest Commit Object Hash
	```
	git log -1 --pretty=format:"%H" 
	```
	copy the hash
		
2. Display the Object Type and Content
	```
	git cat-file -t <commit_hash>  //show object type
	```
	result: commit
	```	
	git cat-file -p <commit_hash> //show object content
	```
	result: 
    ```
	tree b983668ff77aa5f3941643a832772bacaba75901
	parent e95719025323236f566b5cab9a524ba072684d25
	author Ypatios Chaniotakos <ypatioschaniotakos@gmail.com> 1730986555 +0200
	committer Ypatios Chaniotakos <ypatioschaniotakos@gmail.com> 1730986555 +0200
    Added Makefile to run lib/hello.sh
    ```     
--------------
- Dumping Directory Tree:

   Use Git commands to dump the directory tree referenced by this commit.

1. Get the Latest Commit Hash
	```
	latest_commit_hash=$(git rev-parse HEAD)
	```
		
2. Find the Tree Hash Referenced by the Commit
	```
	tree_hash=$(git cat-file -p $latest_commit_hash | grep '^tree' | awk '{print $2}')
	```
		
3. Dump the Directory Tree
	```
	git cat-file -p $tree_hash
	```
The output format will look like this:
```
100644 blob <blob_hash>    filename
040000 tree <subtree_hash> dirname
```
_Blob entries_ represent individual files.

_Tree entries_ represent subdirectories, which you can inspect further by using git cat-file -p <subtree_hash>.

-------------------------------------------------------
- Dump the contents of the lib/ directory and the hello.sh file using Git commands.
  
1. Locate the Latest Commit and Tree Hash 
    ```
	latest_commit_hash=$(git rev-parse HEAD) (get the latest commit hash)
	tree_hash=$(git cat-file -p $latest_commit_hash | grep '^tree' | awk '{print $2}') 
	```
	(Get the Tree Hash of the Commit)
		
2. Locate the Hash for the lib/ Directory
	```
	git cat-file -p $tree_hash 
	```
	Dump the Root Tree to Find lib/ !!!(Note down the <lib_tree_hash> for the lib/ directory)
		
3. Dump the Contents of the lib/ Directory
	```
	git cat-file -p <lib_tree_hash> 
	```
	Replace <lib_tree_hash> with the actual hash you found in Step 2. This will list the files inside lib/, including hello.sh, and provide its object hash.) !!!(note down the hello-blob-hash)
		
4. Dump the Contents of the hello.sh File
	```
	hello_blob_hash=$(git cat-file -p $lib_tree_hash | grep 'hello.sh' | awk '{print $3}')
	git cat-file -p $hello_blob_hash
	```
		
Result:
```
#!/bin/bash
 
# Default is World
# Author: Jim Weirich <jim@example.com>
name=${1:-"World"}

echo "Hello, $name"
```
--------------
## Branching 
It’s time to do a major rewrite of the hello world functionality. Since this might take a while, you’ll want to put these changes into a separate branch to isolate them from changes in the main branch.
- Create and Switch to New Branch:
Create a local branch named greet and switch to it.

1. Create the greet Branch
	```
	git branch greet
	```
		
2. Switch to the greet Branch
	```
	git checkout greet
	```
		
3. Verify the Branch Switch
	```
	git branch
	```
--------------
- In the lib directory, create a new file named greeter.sh and add the provided code to it. Commit these changes.
```
#!/bin/bash
 
Greeter() {
     who="$1"
    echo "Hello, $who"
}
```

1. Navigate to the lib/ Directory
	```
	cd lib
	```
		
2. Create the greeter.sh File with the Provided Code
	```
	echo -e '#!/bin/bash\n\nGreeter() {\n    who="$1"\n    echo "Hello, $who"\n}' > greeter.sh
	```
		
3. Stage the New File
	```
	git add greeter.sh
	```
		
4. Commit the New File
	```
	git commit -m "Add greeter.sh with Greeter function to greet specific users"
    ```
--------------
- Update the lib/hello.sh file by adding the content below, stage and commit the changes.
```
#!/bin/bash

source lib/greeter.sh

name="$1"
if [ -z "$name" ]; then
    name="World"

Greeter "$name"
```

1. Edit lib/hello.sh with the content above.
2. Stage the Changes
    ```
	git add lib/hello.sh
	```
3. Commit the Changes
	```
	git commit -m "Update hello.sh to use Greeter function from greeter.sh"
	```
--------------
- Update the Makefile with the following comment and commit the changes.
```
#Ensure it runs the updated lib/hello.sh file
TARGET="lib/hello.sh"

run:
	bash ${TARGET}
```

1. Edit the Makefile with the content above
2. Stage the changes 
    ```
	git add Makefile
	```
		
3. Commit the changes
	```
	git commit -m "Update Makefile with comment to ensure it runs the updated lib/hello.sh file"
	```
-----------------------
- Switch back to the main branch, compare and show the differences between the main and greet branches for Makefile, hello.sh, and greeter.sh files.

1. Switch to the main Branch
	```
	git checkout master
	```
2. Compare Differences Between the Branches
    
    compare Makefile:
	```
	git diff greet..master -- Makefile 
	```

	compare hello.sh:
	```
	git diff greet..master -- lib/hello.sh 
	```

	compare greeter.sh:
	```
	git diff greet..master -- lib/greeter.sh 
	```

	compare everything:
	```
	git diff greet..master 
	```
--------------------------
- Generate a README.md file for the project with the provided content. Commit this file.
```
This is the Hello World example from the git project.
```

1. Create the README.md File with the content above.
2. Stage the README.md File
    ```
	git add README.md
	```
3. Commit the File
	```
	git commit -m "Add README.md with project description"
	```
------------------
## Conflicts, merging and rebasing

- Merge Main into Greet Branch:
Start by merging the changes from the main branch into the greet branch.

1. Switch to the greet Branch
	```
	git checkout greet
	```
2. Merge the main Branch into greet
	```
	git merge master 
	```
	If there are no conflicts, Git will complete the merge and add a merge commit on the greet branch
		
3. Resolve Any Merge Conflicts (if they occur)

	If there are conflicts, Git will indicate which files are conflicting. Open each conflicted file, resolve the differences, then stage the resolved files:
	```
	git add <file-name>
	```
	After resolving conflicts and staging the changes, complete the merge commit:
    ```
	git commit 
    ```
            
------------------

- Switch to main branch and make the changes below to the hello.sh file, save and commit the changes.
```
#!/bin/bash

echo "What's your name"
read my_name

echo "Hello, $my_name"
```

1. Switch to main branch
    ```
	git checkout master
	```
		
2. Make changes to hello.sh with the above content
	
3. Stage the changes and commit
    ```	
	git add lib/hello.sh
    git commit -m "Update hello.sh to prompt for user name and display personalized greeting"
	```
------------------
- Merging Main into Greet Branch (Conflict):

    Attempt to merge the main branch into greet. Bingooo! There you have it, a conflict.
    Resolve the conflict (manually or using graphical merge tools), accept changes from main branch, then commit the conflict resolution.

1. Switch to the greet Branch
	```
	git checkout greet 
	```
		
2. Attempt to Merge main into greet

    ```
	git merge master
	```
		
3. Resolve the Conflict by Accepting Changes from main
	
    To resolve the conflict and keep the changes from main, edit hello.sh to reflect only the content from the main branch
		
4. Stage the Resolved File
	```
	git add lib/hello.sh
	```
		
5. Complete the Merge with a Commit
	```
	git commit -m "Resolve merge conflict in hello.sh by accepting changes from main"
	```
		
---------------

- Rebasing Greet Branch:

    Go back to the point before the initial merge between main and greet.
    Rebase the greet branch on top of the latest changes in the main branch.

1. Switch to the greet Branch
	```
	git checkout greet
	```
		
2. Reset the greet Branch to Before the Merge
    ```	
	git log --oneline --graph (Note down the hash)
	git reset --hard <noted-hash-from-above> (this will reset greet to previous state before merge)
	```
		
3. Rebase the greet Branch onto main	
	```
	git rebase main
	```
		
4. Resolve Any Conflicts Manually (If Any)
	
5. Stage the resolved files
	```
	git add <filename>
	```
		
6. Continue the rebase after resolving each conflict
	```
	git rebase --continue
	```
		
7. Verify the Rebase (verify that the greet branch is now based on top of main’s latest commits)
	```
	git log 
	```
------------------

- Merging Greet into Main:

    Merge the changes from the greet branch into the main branch.

1. Switch to the main Branch
	```
	git checkout master
	```
		
2. Merge the greet Branch into main (If there are no conflicts, Git will complete the merge and create a merge commit in the main branch)
	```
	git merge greet 
	```
		
3. Resolve Any Conflicts Manually (if they arise)	
	
4. Stage the files and commit
	```
	git add <filename>
	git commit
	```
--------------------

- Understanding Fast-Forwarding and Differences:
Explain fast-forwarding and the difference between merging and rebasing.

A fast-forward happens when a branch can move directly to the latest commit of another branch without creating a merge commit. 
	Merging combines branches with a new merge commit, keeping both histories, while rebasing rewrites commits onto a new base, creating a linear history without diverging branches.

| Feature | Merge | Rebase |
| ------ | ------ | ------ |
| History |  Maintains full history, showing branches and merges |  Linearizes history, eliminating branch divergence |
| Commit Structure |  Creates a new merge commit | Moves (rebases) commits onto a new base |
| Use Case | Collaborative work with multiple feature branches	 | Moves (rebases) commits onto a new base |


----------------
## Local and Remote Repositories

- In the work/ directory, make a clone of the repository hello as cloned_hello. (Do not use copy command)

1. Navigate to the work/ Directory
	
2. Clone the Repository as cloned_hello:
	```
	git clone ./hello hello_cloned
	```
		
3. Verify the clone
	```
	cd hello_cloned
    git status 
	```
	Should get: nothing to commit, working tree clean
		
--------------
- Show the logs for the cloned repository.
    ```
	git log 
	```
- Display the name of the remote repository and provide more information about it. (in hello_cloned)
    ```
	git remote show origin
	```
		
	This command provides detailed information, including:
	```
	Remote URL: The path or URL of the original repository.
	Fetch and Push URLs: Locations Git uses to sync with the remote.
	Remote Tracking Branches: Branches available on the remote.
	Merge Information: Default branch tracking setup for pulling changes.
    ```
- List all remote and local branches in the cloned_hello repository.
    ```
	git branch -a
    ```
------------------
- Make changes to the original repository, update the README.md file with the provided content, and commit the changes.
```
This is the Hello World example from the git project.
(changed in the original)
```

1. Navigate to original repo

2. Make changes to README with above content

3. Stage and commit 
	```
	git add README.md
	git commit -m "Update README.md with Hello World example content"
    ```
-----------
- Inside the cloned repository (cloned_hello), fetch the changes from the remote repository and display the logs. Ensure that commits from the hello repository are included in the logs.

1. Navigate to the cloned repository

2. Fetch the latest changes from the remote repository:
	```
	git fetch
	```
		    
3. View the commit logs:
	```
	git log --oneline --graph --all
	```
----------
- Merge the changes from the remote main branch into the local main branch.

1. Navigate to the cloned repository:
2. Check out the local main branch:
	```
	git checkout master
    ```
		
3. Fetch the latest changes from the remote repository:
	```
	git fetch
	```
		
4. Merge the remote main branch into your local main branch:
	```
	git merge origin/master
    ```
---------------------

- Add a local branch named greet tracking the remote origin/greet branch.

1. Fetch the latest branches from the remote:
	```
	git fetch
	```
		
2. Create and check out the greet branch:	
    ```
	git checkout -b greet origin/greet 
	```
	This command does two things:
    ```
	git checkout -b greet: Creates and switches to a new local branch named greet.
	origin/greet: Specifies that the new local greet branch will track the remote origin/greet branch.
	```		

3. Verify the new branch and tracking:
	```
	git branch -vv (This will show you a list of all local branches, along with their remote tracking information)
	```
--------------------

- Add a remote to your Git repository and push the main and greet branches to the remote. (in hello/)

1. Add a remote repository
	```
	git remote add origin <git-URL>
	```
		
2. Verify the remote
	```
	git remote -v
	```
		
3. Push the main branch to the remote
	```
	git push -u origin master
	```
		
4. Push the greet branch to the remote
	```
	git push -u origin greet
	```
		
5. Verify the branches on the remote
	```
	git branch -r
	```
		
----------

- Be ready for this question in the audit!

"What is the single git command equivalent to what you did before to bring changes from remote to local main branch?"
        ```
	git pull origin master
	    ```
Explanation:
```
- git pull: This command fetches changes from a remote repository and automatically merges them into the current branch.
- origin: This specifies the remote repository (usually named origin by default).
- master: This specifies the branch you want to pull from (in this case, the remote main branch).
```
-----------------

- What is a bare repository and why is it needed?

	A bare repository in Git is a repository that does not have a working directory (i.e., it does not contain the actual files you’re working on). Instead, it only contains the .git directory structure, including all the version history (commits, branches, tags, etc.) and references (such as refs/heads/main). In other words, a bare repository stores all the Git metadata and object data but does not store the actual working files that you would typically edit.
	
-------------------

- Create a bare repository named hello.git from the existing hello repository.

1. Navigate to the existing hello repository:
	
2. Create a bare repository:
	```
	git clone --bare . /path/to/hello.git
	```
		
3. Verify the bare repository:
	```
	ls /path/to/hello.git
	```
		
---------------
- Add the bare hello.git repository as a remote to the original repository hello.

1. Navigate to the original hello repository
	
2. Add the hello.git bare repository as a remote:
	```
	git remote add bare /path/to/hello.git
	```
		
3. Verify that the remote has been added:
	```
	git remote -v
	```
		
	Step 4: Push your local changes to the bare repository:
	```
	git push origin master
	```
		
------------------------

- Change the README.md file in the original repository with the provided content:
```
This is the Hello World example from the git project.
(Changed in the original and pushed to shared)
```

- Commit the changes and push them to the shared repository.

1. Navigate to the original repository directory (hello)
	
2. Edit the README.md file with the content above
	
3. Stage and commit the changes
	```
	git add README.md
    git commit -m  "Updated README.md with new content: This is the Hello World example from the git project."
    ```
------------------------
- Switch to the cloned repository cloned_hello and pull down the changes just pushed to the shared repository.

1. Navigate to the cloned_hello repository:
	
2. Pull the changes from the shared repository:
	```
	git pull origin master
	```
----------------


