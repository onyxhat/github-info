package main

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func main() {
	//Assign cli flags
	pflag.String("token", "", "--token ghp_XXXXXXXXXXXXXXXX")

	//Parse cli flags
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	//Check env
	viper.SetEnvPrefix("github")
	viper.BindEnv("token")

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: viper.GetString("token")})
	httpClient := oauth2.NewClient(context.Background(), src)
	gh_client := githubv4.NewClient(httpClient)
	myContext := context.Background()

	data := getTemplateChildren(*gh_client, myContext)

	fmt.Println(string(data))
}
