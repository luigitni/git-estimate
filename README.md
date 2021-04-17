# Git Estimate

Estimates hours and days spent by all developers on a git repository.

```
git-estimate -repo="..\flamel"

commits by luigi.tanzini@decodica.com
=== 0.25 days (2.00 hours)

commits by dev1
=== 0.50 days (4.01 hours)

commits by dev2
=== 1.48 days (11.83 hours)

commits by dev3
=== 0.25 days (2.00 hours)

commits by dev4
=== 10.15 days (81.23 hours)

commits by dev5
=== 27.95 days (223.59 hours)

overall 40.58 days (324.66 hours)
```

The code is written in plain go and uses [go-git](https://github.com/src-d/go-git) to read the repo's commits.

## Install

Clone the repository then simply build git-estimate from source using the goo tools:

```shell script
cd path/to/your/cloned/repo
go build git-estimate
```

## Usage

At a minimum run:

```
git-estimate -repo=/path/to/repo
```

this will use default settings to compute the time spent on the repo at the specified path.

```
-baseline float
        baseline value for session estimate (default 2)
-estimate string
        estimation method. Accepted values are "session" and "day". (default "session")
-from string
        if provided computation starts from the given date and time. Format is yyyy-mm-ddThh-ii
-json
        if true will output estimates in JSON format
-repo string
        git repository path. If no flag is specified the current folder is assumed (default ".")
-group
        group estimates based on comment message content using a predefined or custom pattern. Custom patterns should identify exactly 1 capturing group. See https://github.com/google/re2/wiki/Syntax for syntax.
        Predefined patterns available:
    
                jira - Captures the first Jira issue key based on the smart commit format (https://support.atlassian.com/bitbucket-cloud/docs/use-smart-commits/)
                type - Captures the type component of conventional commit messages (https://www.conventionalcommits.org/en/v1.0.0/)
                scope - Captures the scope component of conventional commit messages (https://www.conventionalcommits.org/en/v1.0.0/)
```

### Output

You can also specify the result to be output in JSON format, should you need to use the program in a pipeline.

```shell script
git-estimate -repo=/path/to/repo -json
```

will output:

```json5
{ 
   "developers":[ 
      { 
         "author":"dev1@decodica.com",
         "hours":223.58805555555554,
         "days":27.948506944444443
      },
      { 
         "author":"luigi.tanzini@decodica.com",
         "hours":2,
         "days":0.25
      },
      { 
         "author":"dev2@decodica.com",
         "hours":4.006111111111111,
         "days":0.5007638888888889
      },
      { 
         "author":"dev3@decodica.com",
         "hours":11.831944444444444,
         "days":1.4789930555555555
      },
      { 
         "author":"dev4@decodica.com",
         "hours":2,
         "days":0.25
      },
      { 
         "author":"dev5@decodica.com",
         "hours":81.23,
         "days":10.15375
      }
   ],
   "overall":{ 
      "author":"all",
      "hours":324.6561111111111,
      "days":40.58201388888889
   }
}
```

### Grouping

By default, commits are grouped by author email address. This can be further subgrouped based on content matching in the commit message. Specify a predefined pattern after the `-group` flag or provide a custom [regex pattern](https://github.com/google/re2/wiki/Syntax) with exactly 1 capture group.

* `-group jira` - Captures the first Jira issue key based on the [smart commit format](https://support.atlassian.com/bitbucket-cloud/docs/use-smart-commits/)
* `-group type` - Captures the type component of [conventional commit messages](https://www.conventionalcommits.org/en/v1.0.0/)
* `-group scope` - Captures the scope component of [conventional commit messages](https://www.conventionalcommits.org/en/v1.0.0/)

## Estimates

Currently the software supports two simple methods of estimation:

#### Working Session

This estimation is perhaps the most accurate.
It assumes the following:

- A "day" output is a working day made of 8 hours.
- If more than 8 hours have passed between a commit and the next one, the former was the last commit of the session.
- The *first* commit of each working session took 2 hours of work. You can configure this padding using the ```baseline``` flag when running *git-estimate*

#### Working Day

This estimation assumes the following:

- If a commit has been done during the day, that day counts toward the total
- A day is made of 8 working hours

## Contributors

Feel free to contribute to the project however you want.

The code should make it easy to add estimation methods and/or output formatter should you need different format or require additional estimation methods.

Please use *go fmt* before a pull request.

#### Inspiration

The software was inspired by [git-hours](https://github.com/kimmobrunfeldt/git-hour). I wasn't able to find anything similar which required less tools to get started so I decided to hack together a quick software that would be as simple and straight to the point.

## License

The software is released under the [MIT](https://github.com/luigitni/git-estimate/blob/master/LICENSE.txt) license.
