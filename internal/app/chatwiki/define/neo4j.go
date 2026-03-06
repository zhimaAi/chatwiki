// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import (
	"context"
	"errors"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
)

var Neo4jStatus map[int]bool
var neo4JDriver neo4j.DriverWithContext

func GetNeo4jDriver() (neo4j.DriverWithContext, error) {
	if !cast.ToBool(Config.Neo4j[`enabled`]) {
		return nil, errors.New(`knowledge graph is not enabled`)
	}
	if neo4JDriver != nil {
		return neo4JDriver, nil // has been initialized
	}
	dbUri := fmt.Sprintf("neo4j://%s:%s", Config.Neo4j["host"], Config.Neo4j["port"])
	dbUser := Config.Neo4j["user"]
	dbPassword := Config.Neo4j["password"]
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		return nil, err
	}
	if err = driver.VerifyConnectivity(context.Background()); err != nil {
		return nil, err
	}
	logs.Info(`initialize neo4j finish`)
	neo4JDriver = driver // record the neo4j handle
	return neo4JDriver, nil
}
