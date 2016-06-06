package api

import (
	"github.com/kwo/rakewire/api/msg"
	"github.com/kwo/rakewire/auth"
	"github.com/kwo/rakewire/model"
	"golang.org/x/net/context"
)

// GroupList lists a user's groups.
func (z *API) GroupList(ctx context.Context, req *msg.GroupListRequest) (*msg.GroupListResponse, error) {

	user := ctx.Value("user").(*auth.User)

	rsp := &msg.GroupListResponse{}

	err := z.db.Select(func(tx model.Transaction) error {

		groups := model.G.GetForUser(tx, user.ID)

		for _, group := range groups {
			g := &msg.Group{
				Name: group.Name,
			}
			rsp.Groups = append(rsp.Groups, g)
		}

		return nil

	})

	return rsp, err

}
