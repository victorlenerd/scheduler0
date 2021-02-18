package project

import (
	"errors"
	"github.com/go-pg/pg"
	"github.com/victorlenerd/scheduler0/server/src/models"
	"github.com/victorlenerd/scheduler0/server/src/utils"
	"net/http"
)

type ProjectManager models.ProjectModel

func (projectManager *ProjectManager) CreateOne(pool *utils.Pool) (string, *utils.GenericError) {
	conn, err := pool.Acquire()
	if err != nil {
		return "", utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	db := conn.(*pg.DB)
	defer pool.Release(conn)

	if len(projectManager.Name) < 1 {
		return "", utils.HTTPGenericError(http.StatusBadRequest, "name field is required")
	}

	if len(projectManager.Description) < 1 {
		return "", utils.HTTPGenericError(http.StatusBadRequest, "description field is required")
	}

	projectWithName := ProjectManager{
		Name: projectManager.Name,
	}

	_ = projectWithName.GetOneByName(pool)
	if len(projectWithName.UUID) > 5 {
		return "", utils.HTTPGenericError(http.StatusBadRequest, "Another project exist with the same name")
	}

	_, err = db.Model(projectManager).Insert()
	if err != nil {
		return "", utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	return projectManager.UUID, nil
}

func (projectManager *ProjectManager) GetOneByName(pool *utils.Pool) *utils.GenericError {
	conn, err := pool.Acquire()
	if err != nil {
		return utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	db := conn.(*pg.DB)
	defer pool.Release(conn)

	projectManagers := []ProjectManager{}

	err = db.
		Model(&projectManagers).
		Where("name = ?", projectManager.Name).
		Select()

	if len(projectManagers) < 1 {
		return utils.HTTPGenericError(http.StatusNotFound, "project with name : "+projectManager.Name+" does not exist")
	}

	if err != nil {
		return utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	projectManager.Name = projectManagers[0].Name
	projectManager.UUID = projectManagers[0].UUID
	projectManager.Description = projectManagers[0].Description
	projectManager.DateCreated = projectManagers[0].DateCreated
	projectManager.ID = projectManagers[0].ID

	return nil
}

func (projectManager *ProjectManager) GetOneByUUID(pool *utils.Pool) *utils.GenericError {
	conn, err := pool.Acquire()
	if err != nil {
		return utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	db := conn.(*pg.DB)
	defer pool.Release(conn)

	err = db.
		Model(projectManager).
		Where("uuid = ?", projectManager.UUID).
		Select()

	if err != nil {
		return utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	return  nil
}

func (projectManager *ProjectManager) GetAll(pool *utils.Pool, offset int, limit int) ([]ProjectManager, *utils.GenericError) {
	conn, err := pool.Acquire()
	if err != nil {
		return nil, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	db := conn.(*pg.DB)
	defer pool.Release(conn)

	projects := make([]ProjectManager, 0, limit)

	err = db.Model(&projects).
		Order("date_created").
		Offset(offset).
		Limit(limit).
		Select()

	if err != nil {
		return nil, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	return projects, nil
}

func (projectManager *ProjectManager) GetTotalCount(pool *utils.Pool) (int, *utils.GenericError) {
	conn, err := pool.Acquire()
	if err != nil {
		return 0, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	db := conn.(*pg.DB)
	defer pool.Release(conn)

	count, err := db.Model(projectManager).
		Order("date_created").
		Count()

	if err != nil {
		return 0, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	return count, nil
}

func (projectManager *ProjectManager) UpdateOne(pool *utils.Pool) (int, *utils.GenericError) {
	conn, err := pool.Acquire()
	if err != nil {
		return 0, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	db := conn.(*pg.DB)
	defer pool.Release(conn)

	savedProject := ProjectManager{UUID: projectManager.UUID}

	if e := savedProject.GetOneByUUID(pool); e != nil {
		return 0, e
	}


	if len(savedProject.UUID) < 1 {
		return 0, utils.HTTPGenericError(http.StatusBadRequest, "project does not exist")
	}

	if savedProject.Name != projectManager.Name {
		projectWithSimilarName := ProjectManager{
			Name: projectManager.Name,
		}

		e := projectWithSimilarName.GetOneByName(pool)
		if e != nil && e.Type != http.StatusNotFound {
			return 0, e
		}
	}

	res, err := db.Model(projectManager).Where("UUID = ?", projectManager.UUID).Update(projectManager)
	if err != nil {
		return 0, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	return res.RowsAffected(), nil
}

func (projectManager *ProjectManager) DeleteOne(pool *utils.Pool) (int, *utils.GenericError) {
	conn, err := pool.Acquire()
	if err != nil {
		return -1, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}
	db := conn.(*pg.DB)
	defer pool.Release(conn)

	jobs := []models.JobModel{}

	err = db.Model(&jobs).Where("project_uuid = ?", projectManager.UUID).Select()
	if err != nil {
		return -1, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	if len(jobs) > 0 {
		err = errors.New("cannot delete projects with jobs")
		return -1, utils.HTTPGenericError(http.StatusBadRequest, err.Error())
	}

	r, err := db.Model(projectManager).Where("uuid = ?", projectManager.UUID).Delete()
	if err != nil {
		return -1, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	return r.RowsAffected(), nil
}