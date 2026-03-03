package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/repository"
)

type CategoryHandler struct {
	Repo *repository.Queries
}

type GetCategoryRequest struct{}
type GetCategoryResponse struct {
	Body struct {
		Categories []CategoryEntry `json:"categories"`
	}
}

type CategoryEntry struct {
	CategoryInfo  repository.Category      `json:"category_info"`
	Subcategories []repository.Subcategory `json:"subcategories"`
}

func (h *CategoryHandler) Get(ctx context.Context, req *GetCategoryRequest) (*GetCategoryResponse, error) {
	resp := new(GetCategoryResponse)

	categories, err := h.Repo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	// Optimized
	resp.Body.Categories = make([]CategoryEntry, 0, len(categories))

	for _, cat := range categories {
		subcategories, err := h.Repo.GetSubcategoriesFor(ctx, repository.GetSubcategoriesForParams{
			CategoryID: cat.ID,
		})
		if err != nil {
			return nil, err
		}

		resp.Body.Categories = append(resp.Body.Categories, CategoryEntry{
			CategoryInfo:  cat,
			Subcategories: subcategories,
		})
	}

	return resp, nil
}
