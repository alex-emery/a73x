---
title: home 
nav: "/" 
---
Welcome!

# me
*   backend cloud software engineer
*   lang: go
*   infra: kubernetes

# tidbits

## #go

*   layout packages by what they do, not by their abstract type
*   use channels sparingly: write synchronous methods and allow the caller to make it async
*   `append` modifies the underlying slice, you'll only make this mistake once
*   define interfaces where you use them
*   `make([]int, 5)` has a length and capacity of 5. `([]int, 0,5)` has a length of 0 and capacity of 5.  
    `append()` will only do what you want with the latter
*   don't use `init()`
*   TFBO (test, fix, benchmark, optimise)
*   more CPU != more performance  
    more CPU == more contention

## #git

*   `git reflog`
*   `git commit --fixup=<COMMITISH>`  
    `git rebase origin/main --autosquash`

# resources

*   [`proc.go`](https://cs.opensource.google/go/go/+/refs/tags/go1.23.0:src/runtime/proc.go)
*   [proposal: runtime/metrics: define a recommended set of metrics](https://github.com/golang/go/issues/67120)

# books

*   [Designing Data Intensive Applications](https://www.oreilly.com/library/view/designing-data-intensive-applications/9781491903063)
*   [Database Internals](https://www.oreilly.com/library/view/database-internals/9781492040330)
*   [Efficient Go](https://www.oreilly.com/library/view/efficient-go/9781098105709)