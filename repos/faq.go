package repos

import (
	"context"
	"database/sql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	models "test-manager/usecase_models/boiler"
)

type FaqsRepository interface {
	UpdateFaqs(ctx context.Context, Faq models.Faq) error
	GetFaqs(ctx context.Context) ([]*models.Faq, error)
	GetFaq(ctx context.Context, id int) (models.Faq, error)
	SaveFaqs(ctx context.Context, Faq models.Faq) (int, error)
}

type faqsRepository struct {
	db *sql.DB
}

func NewFaqsRepository(db *sql.DB) FaqsRepository {
	return &faqsRepository{db: db}
}

func (r *faqsRepository) SaveFaqs(ctx context.Context, Faq models.Faq) (int, error) {
	err := Faq.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return Faq.ID, nil
}

func (r *faqsRepository) UpdateFaqs(ctx context.Context, Faq models.Faq) error {
	_, err := Faq.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (r *faqsRepository) GetFaqs(ctx context.Context) ([]*models.Faq, error) {
	Faq, err := models.Faqs().All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return Faq, nil
}

func (r *faqsRepository) GetFaq(ctx context.Context, id int) (models.Faq, error) {
	Faq, err := models.Faqs(models.FaqWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return models.Faq{}, err
	}
	return *Faq, nil
}
