package repository

import (
	"context"
	"log"
	"zzz/internal/dto"
	"zzz/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DataRepo struct {
	db *gorm.DB
}

func NewDataRepo(db *gorm.DB) *DataRepo {
	return &DataRepo{db: db}
}

func (r *DataRepo) SaveProducts(ctx context.Context, products []dto.ProductSource) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, p := range products {
			brand := models.Brand{Name: p.Brand}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "name"}},
				DoUpdates: clause.AssignmentColumns([]string{"name"}),
			}).Create(&brand).Error; err != nil {
				return err
			}

			category := models.Category{Name: p.Category}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "name"}},
				DoUpdates: clause.AssignmentColumns([]string{"name"}),
			}).Create(&category).Error; err != nil {
				return err
			}

			price := p.ParsePrice()
			product := models.Product{
				ID:         p.ID,
				Name:       p.Name,
				BrandID:    brand.ID,
				CategoryID: category.ID,
				Price:      price,
				Stock:      p.Stock,
			}
			log.Println("brand id =", brand.ID)

			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"name", "brand_id", "category_id", "price", "stock",
				}),
			}).Create(&product).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *DataRepo) SaveClients(ctx context.Context, clients []dto.ClientSource) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, c := range clients {
			client := models.Client{
				ID:        c.ID,
				FirstName: c.FirstName,
				LastName:  c.LastName,
			}

			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name"}),
			}).Create(&client).Error; err != nil {
				return err
			}

			if err := tx.Where("client_id = ?", client.ID).Delete(&models.ClientProduct{}).Error; err != nil {
				return err
			}

			if len(c.Products) > 0 {
				var links []models.ClientProduct
				for _, prodID := range c.Products {
					links = append(links, models.ClientProduct{
						ClientID:  client.ID,
						ProductID: prodID,
					})
				}

				if err := tx.Clauses(clause.OnConflict{
					DoNothing: true,
				}).Create(&links).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
