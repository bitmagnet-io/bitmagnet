package query

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"sync"
)

func GenericBatch[T interface{}](
	ctx context.Context,
	daoQ *dao.Query,
	option Option,
	tableName string,
	factory SubQueryFactory,
	fn func(tx *dao.Query, r []T) error,
) error {
	return daoQ.Transaction(func(tx *dao.Query) error {
		dbCtx := dbContext{
			q:         tx,
			tableName: tableName,
			factory:   factory,
		}
		builder, optionErr := option(newQueryContext(dbCtx))
		ctx = builder.createContext(ctx)
		if optionErr != nil {
			return optionErr
		}
		sq := factory(ctx, tx)
		if selectErr := builder.applySelect(sq); selectErr != nil {
			return selectErr
		}
		if preErr := builder.applyPre(sq); preErr != nil {
			return preErr
		}
		if postErr := builder.applyPost(sq); postErr != nil {
			return postErr
		}
		txCtx := dbContext{
			q:         tx,
			tableName: tableName,
			factory:   factory,
		}
		limit := 1000
		offset := 0
		var items []T
		for {
			if findErr := sq.UnderlyingDB().Limit(limit).Offset(offset).Find(&items).Error; findErr != nil {
				return findErr
			}
			if cbErr := builder.applyCallbacksCtx(ctx, callbackContext{
				dbContext: txCtx,
				Mutex:     &sync.Mutex{},
			}, items); cbErr != nil {
				return cbErr
			}
			if fnErr := fn(tx, items); fnErr != nil {
				return fnErr
			}
			if len(items) < limit {
				break
			}
			offset += limit
			//if offset > limit {
			//	break
			//}
			//if offset > 5000 {
			//	break
			//}
		}
		return nil
	})
}
