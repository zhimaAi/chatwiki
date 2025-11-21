package migrations

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func init() {
	goose.AddMigrationNoTxContext(func(ctx context.Context, db *sql.DB) error {
		return VectorEmbedding2000Migration()
	}, nil)
}

func VectorEmbedding2000Migration() error {
	var maxId int
	indexModel := msql.Model(`chat_ai_library_file_data_index`, define.Postgres)
	if maxIdStr, err := indexModel.Max(`id`); err != nil {
		return err
	} else {
		maxId = cast.ToInt(maxIdStr)
	}
	logs.Other(`migration`, `获取最大ID:%d`, maxId)
	var size = 1000 //每一批次数
	for i := 0; ; i++ {
		start, end := i*size, (i+1)*size
		logs.Other(`migration`, `第%d轮:%d~%d`, i+1, start, end)
		affect, err := indexModel.Where(`id`, `>`, cast.ToString(start)).
			Where(`id`, `<=`, cast.ToString(end)).
			Where(`vector_dims(embedding)=2000`).
			Update2(`embedding2000=embedding,embedding=NULL`)
		if err != nil {
			return err
		}
		logs.Other(`migration`, `第%d轮:affect(%d)`, i+1, affect)
		if end >= maxId {
			break //处理完毕,结束循环
		}
	}
	return nil
}
