# GitHub PR Comments

Create/update GitHub PR comments through an executable (for e.g. in ci/cd scripts)

## Usage

|Option|Description|Required|
|-|-|-|
|`-token=<github_token>`|GitHub token|true|
|`-pr=<pr_number>`|Pull request number|true|
|`-owner=<repo_owner>`|GitHub repository owner|true|
|`-repo=<repo_name>`|GitHub repository name|true|
|`-body=<text>`|The contents of the comment|true|
|`-commentId=<comment_id>`|Update a certain comment by its ID||
|`-bodyIncludes=<text>`|Update a certain comment by its content||

> [!NOTE]
> If no `commentId` or `bodyIncludes` is provided, it will always create a new comment, else it only creates a new comment when not found.
