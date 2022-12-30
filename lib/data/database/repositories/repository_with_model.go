package repositories

import "lib/functional"

type RepositoryWithModel[TEntity any, TModel any] struct {
	Repository[TEntity]
	ToModel  func(entity TEntity) (TModel, error)
	ToEntity func(model TModel) (TEntity, error)
}

func (r *RepositoryWithModel[TEntity, TModel]) GetAll() ([]TModel, error) {
	res, err := r.Repository.GetAll()
	if err != nil {
		return nil, err
	}

	var firstError error = nil

	models := functional.Map(res, func(entity TEntity) *TModel {
		model, err := r.ToModel(entity)
		if err != nil {
			if firstError == nil {
				firstError = err
			}

			return nil
		}
		return &model
	})

	if firstError != nil {
		return nil, firstError
	}

	filteredModels := functional.Map(models, func(model *TModel) TModel {
		return *model
	})

	return filteredModels, nil
}

func (r *RepositoryWithModel[TEntity, TModel]) GetById(id uint) (*TModel, error) {
	entity, err := r.Repository.GetById(id)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, nil
	}

	model, err := r.ToModel(*entity)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *RepositoryWithModel[TEntity, TModel]) SaveModel(model TModel) error {
	entity, err := r.ToEntity(model)
	if err != nil {
		return err
	}
	return r.Repository.Save(entity)
}

func (r *RepositoryWithModel[TEntity, TModel]) SaveEntity(entity TEntity) error {
	return r.Repository.Save(entity)
}

func (r *RepositoryWithModel[TEntity, TModel]) DeleteModel(model TModel) error {
	entity, err := r.ToEntity(model)
	if err != nil {
		return err
	}
	return r.Repository.Delete(entity)
}

func (r *RepositoryWithModel[TEntity, TModel]) DeleteEntity(entity TEntity) error {
	return r.Repository.Delete(entity)
}

func (r *RepositoryWithModel[TEntity, TModel]) CreateModel(model TModel) error {
	entity, err := r.ToEntity(model)
	if err != nil {
		return err
	}
	return r.Repository.Create(entity)
}

func (r *RepositoryWithModel[TEntity, TModel]) CreateEntity(entity TEntity) error {
	return r.Repository.Create(entity)
}
