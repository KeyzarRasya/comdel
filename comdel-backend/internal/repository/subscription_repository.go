package repository

import "github.com/KeyzarRasya/comdel-server/internal/dto"

type SubscriptionRepository interface {
	Save(subscription dto.Subscription)					error;
	Activate(plan string, subsId string, userId string)	error;
	Deactivate(userId string)							error;
}