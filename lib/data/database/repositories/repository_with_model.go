package repositories

import "lib/functional"

type DbRepositoryWithModel[TEntity any, TModel any] struct {
	DbRepository[TEntity]
	ToModel  func(entity TEntity) (TModel, error)
	ToEntity func(model TModel) (*TEntity, error)
}

func (r *DbRepositoryWithModel[TEntity, TModel]) GetAll() ([]TModel, error) {
	res, err := r.DbRepository.GetAll()
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

func (r *DbRepositoryWithModel[TEntity, TModel]) GetById(id uint) (*TModel, error) {
	entity, err := r.DbRepository.GetById(id)
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

func (r *DbRepositoryWithModel[TEntity, TModel]) SaveModel(model TModel) error {
	entity, err := r.ToEntity(model)
	if err != nil {
		return err
	}
	return r.DbRepository.Save(*entity)
}

func (r *DbRepositoryWithModel[TEntity, TModel]) SaveEntity(entity TEntity) error {
	return r.DbRepository.Save(entity)
}

func (r *DbRepositoryWithModel[TEntity, TModel]) DeleteModel(model TModel) error {
	entity, err := r.ToEntity(model)
	if err != nil {
		return err
	}
	return r.DbRepository.Delete(*entity)
}

func (r *DbRepositoryWithModel[TEntity, TModel]) DeleteEntity(entity TEntity) error {
	return r.DbRepository.Delete(entity)
}

func (r *DbRepositoryWithModel[TEntity, TModel]) CreateModel(model TModel) error {
	entity, err := r.ToEntity(model)
	if err != nil {
		return err
	}
	return r.DbRepository.Create(entity)
}

func (r *DbRepositoryWithModel[TEntity, TModel]) CreateEntity(entity *TEntity) error {
	return r.DbRepository.Create(entity)
}
