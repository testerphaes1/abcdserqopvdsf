package repos

import (
	"context"
	"database/sql"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	models "test-manager/usecase_models/boiler"
)

type OrdersRepository interface {
	UpdateOrders(ctx context.Context, Order models.Order) error
	GetOrders(ctx context.Context, accountId int, projectId int) ([]*models.Order, error)
	GetOrder(ctx context.Context, id int) (models.Order, error)
	SaveOrders(ctx context.Context, Order models.Order) (int, error)
	GetOrderByGatewayID(ctx context.Context, gatewayOrderId string) (models.Order, error)
}

type ordersRepository struct {
	db *sql.DB
}

func NewOrdersRepository(db *sql.DB) OrdersRepository {
	return &ordersRepository{db: db}
}

func (r *ordersRepository) SaveOrders(ctx context.Context, Order models.Order) (int, error) {
	err := Order.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return Order.ID, nil
}

func (r *ordersRepository) UpdateOrders(ctx context.Context, Order models.Order) error {
	_, err := Order.Update(ctx, r.db, boil.Blacklist("account_id", "project_id", "package_id", "gateway_id", "amount"))
	if err != nil {
		return err
	}
	return nil
}

func (r *ordersRepository) GetOrders(ctx context.Context, accountId int, projectId int) ([]*models.Order, error) {
	if projectId == 0 {
		Order, err := models.Orders(models.OrderWhere.AccountID.EQ(accountId)).All(ctx, r.db)
		if err != nil {
			return nil, err
		}
		return Order, nil
	} else {
		Order, err := models.Orders(models.OrderWhere.AccountID.EQ(accountId), models.OrderWhere.ProjectID.EQ(projectId)).All(ctx, r.db)
		if err != nil {
			return nil, err
		}
		return Order, nil
	}
}

func (r *ordersRepository) GetOrder(ctx context.Context, id int) (models.Order, error) {
	Order, err := models.Orders(models.OrderWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return models.Order{}, err
	}
	return *Order, nil
}

func (r *ordersRepository) GetOrderByGatewayID(ctx context.Context, gatewayOrderId string) (models.Order, error) {
	Order, err := models.Orders(models.OrderWhere.GatewayOrderID.EQ(null.NewString(gatewayOrderId, true))).One(ctx, r.db)
	if err != nil {
		return models.Order{}, err
	}
	return *Order, nil
}
