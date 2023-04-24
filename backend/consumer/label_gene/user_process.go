package label_gene

import (
	"backend/cmd/dal/query"
	"context"
	"fmt"
)

func processUserLabel(ctx context.Context, appId int64, labelId int64) map[int64]string {
	res := make(map[int64]string)

	// 用户数据
	users, err := query.User.WithContext(ctx).Where(query.User.AppID.Eq(appId)).Find()
	if err != nil {
		return res
	}

	switch labelId {
	case Gender:
		for _, user := range users {
			res[user.UserID] = fmt.Sprintf("%d", user.UserGender)
		}
	case Age:
		for _, user := range users {
			res[user.UserID] = fmt.Sprintf("%d", user.UserAge)
		}
	case Career:
		for _, user := range users {
			if user.UserCareer == nil {
				continue
			}
			res[user.UserID] = fmt.Sprintf("%s", *user.UserCareer)
		}
	default:

	}

	return res
}
