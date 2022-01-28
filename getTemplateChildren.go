package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/shurcooL/githubv4"
)

func getTemplateChildren(client githubv4.Client, ctx context.Context) []byte {
	var q struct {
		Viewer struct {
			Repositories struct {
				Edges []struct {
					Node struct {
						TemplateRepository struct {
							NameWithOwner githubv4.String
						}
						NameWithOwner githubv4.String
					}
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage githubv4.Boolean
				}
			} `graphql:"repositories(first: 100, after: $endCursor)"`
		}
	}

	v := map[string]interface{}{
		"endCursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	var list []Repo
	for {
		err := client.Query(ctx, &q, v)
		if err != nil {
			log.Println(err)
		}

		//Filter responses with no parent template repo
		for _, repo := range q.Viewer.Repositories.Edges {
			if repo.Node.TemplateRepository.NameWithOwner != "" {
				list = append(list, Repo{Name: repo.Node.NameWithOwner, Template: repo.Node.TemplateRepository.NameWithOwner})
			}
		}

		//Escape if there are no more pages
		if !q.Viewer.Repositories.PageInfo.HasNextPage {
			break
		}

		//Set next cursor
		v["endCursor"] = githubv4.NewString(q.Viewer.Repositories.PageInfo.EndCursor)
	}

	output, _ := json.MarshalIndent(list, "", "  ")

	return output
}
