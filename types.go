package main

import "github.com/shurcooL/githubv4"

type Repo struct {
	Name     githubv4.String
	Template githubv4.String
}
