package utils

import (
	"api/internal/types"
	"context"

	"github.com/machinebox/graphql"
)

func FindData(page int) []types.HashnodePost {

	graphqlClient := graphql.NewClient("https://api.hashnode.com")
	graphqlRequest := graphql.NewRequest(`
        query($page: Int!){
          user(username: "solanaroundup") {
            publication {
              posts(page: $page) {
                title
                brief
                slug
                dateAdded
                contentMarkdown
                coverImage
              }
            }
          }
        }
    `)

	graphqlRequest.Var("page", page)

	var graphqlResponse types.HashNodeUser

	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		panic(err)
	}

	return graphqlResponse.User["publication"].Posts
}
