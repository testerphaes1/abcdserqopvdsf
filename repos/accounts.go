package repos

import (
	"context"
	"database/sql"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	models "test-manager/usecase_models/boiler"
)

type AccountsRepository interface {
	UpdateAccounts(ctx context.Context, Account models.Account) error
	GetAccounts(ctx context.Context, id int) (models.Account, error)
	GetAccountByEmail(ctx context.Context, username string) (models.Account, error)
	AccountExistsByEmail(ctx context.Context, username string) (bool, error)
	SaveAccounts(ctx context.Context, Account models.Account) (int, error)
}

type accountsRepository struct {
	db *sql.DB
}

func NewAccountsRepositoryRepository(db *sql.DB) AccountsRepository {
	return &accountsRepository{db: db}
}

func (r *accountsRepository) SaveAccounts(ctx context.Context, Account models.Account) (int, error) {
	err := Account.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return Account.ID, nil
}

func (r *accountsRepository) UpdateAccounts(ctx context.Context, Account models.Account) error {
	_, err := Account.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (r *accountsRepository) GetAccounts(ctx context.Context, id int) (models.Account, error) {
	Account, err := models.Accounts(models.AccountWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return models.Account{}, err
	}
	return *Account, nil
}

func (r *accountsRepository) GetAccountByEmail(ctx context.Context, email string) (models.Account, error) {
	Account, err := models.Accounts(models.AccountWhere.Email.EQ(null.NewString(email, true))).One(ctx, r.db)
	if err != nil {
		return models.Account{}, err
	}
	return *Account, nil
}

func (r *accountsRepository) AccountExistsByEmail(ctx context.Context, email string) (bool, error) {
	return models.Accounts(models.AccountWhere.Email.EQ(null.NewString(email, true))).Exists(ctx, r.db)
}
