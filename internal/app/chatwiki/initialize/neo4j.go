// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package initialize

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/cast"
)

func initNeo4j() {
	if !cast.ToBool(define.Config.Neo4j["enabled"]) {
		return
	}
	var err error
	ctx := context.Background()
	dbUri := fmt.Sprintf("neo4j://%s:%s", define.Config.Neo4j["host"], define.Config.Neo4j["port"])
	dbUser := define.Config.Neo4j["user"]
	dbPassword := define.Config.Neo4j["password"]
	define.Neo4jDriver, err = neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	if err = define.Neo4jDriver.VerifyConnectivity(ctx); err != nil {
		panic(err)
	}
}
