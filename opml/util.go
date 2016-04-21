package opml

import (
	"rakewire/model"
)

// Flatten pulls nested groups up to the top level, modifing the group name.
func flatten(outlines Outlines) map[*Outline]Outlines {
	result := make(map[*Outline]Outlines)
	_flatten(outlines, &Outline{}, result)
	return result
}

func _flatten(outlines Outlines, branch *Outline, result map[*Outline]Outlines) {
	for _, outline := range outlines {
		if outline.Type == "rss" {
			result[branch] = append(result[branch], outline)
		} else {
			// inherited from parent branch
			b := &Outline{}
			b.SetAutoRead(branch.IsAutoRead() || outline.IsAutoRead())
			b.SetAutoStar(branch.IsAutoStar() || outline.IsAutoStar())
			if len(branch.Title) == 0 {
				b.Title = outline.Title
				b.Text = outline.Text
				_flatten(outline.Outlines, b, result) // 2nd level
			} else {
				b.Title = branch.Title + "/" + outline.Title
				b.Text = branch.Text + "/" + outline.Text
				_flatten(outline.Outlines, b, result) // 3rd + levels only
			}
		}
	}
}

func groupOutlinesByURL(flatOPML map[*Outline]Outlines) map[string]*Outline {
	result := make(map[string]*Outline)
	for _, outlines := range flatOPML {
		for _, outline := range outlines {
			if _, ok := result[outline.XMLURL]; !ok {
				result[outline.XMLURL] = outline
			}
		}
	}
	return result
}

func groupSubscriptionsByGroup(subscriptions model.Subscriptions, groups map[string]*model.Group) map[*model.Group]model.Subscriptions {

	result := make(map[*model.Group]model.Subscriptions)
	for _, subscription := range subscriptions {
		for _, groupID := range subscription.GroupIDs {
			result[groups[groupID]] = append(result[groups[groupID]], subscription)
		}
	}
	return result

}

func groupSubscriptionsByURL(subscriptions model.Subscriptions, feedsByID map[string]*model.Feed) (map[string]*model.Subscription, model.Subscriptions) {

	result := make(map[string]*model.Subscription)
	duplicates := model.Subscriptions{}
	for _, subscription := range subscriptions {
		if feed := feedsByID[subscription.FeedID]; feed != nil {
			if _, ok := result[feed.URL]; !ok {
				result[feed.URL] = subscription
			} else {
				duplicates = append(duplicates, subscription)
			}
		}
	}

	return result, duplicates

}
