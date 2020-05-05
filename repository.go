package main

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) AddTasks(tasks []Task) ([]Task, error) {
	var addedTasks []Task

	tx, err := r.db.Begin()
	if err != nil {
		return addedTasks, err
	}

	for _, v := range tasks {
		var insertId int64
		err = tx.QueryRow(
			"INSERT INTO task(url, count) VALUES($1, $2) RETURNING id",
			v.Url,
			v.NumberOfRequests).Scan(&insertId)
		if err != nil {
			tx.Rollback()
			return addedTasks, err
		}

		addedTasks = append(addedTasks, Task{Id: insertId, Url: v.Url, NumberOfRequests: v.NumberOfRequests})
	}

	err = tx.Commit()
	return addedTasks, err
}

func (r *Repository) IncSuccessTask(task Task) (err error) {
	stmt, err := r.db.Prepare("UPDATE task SET success = success + 1 WHERE id = $1")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(task.Id)
	return
}

func (r *Repository) IncErrorTask(task Task) (err error) {
	stmt, err := r.db.Prepare("UPDATE task SET error = error + 1 WHERE id = $1")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(task.Id)
	return
}

func (r *Repository) GetStatus() (status Status, err error) {
	err = r.db.QueryRow(`select
       count(*) as task_total,
       COALESCE(sum(case when success + error != count then 1 else 0 end), 0) as task_queue,
       COALESCE(sum(case when success = count then 1 else 0 end), 0) as task_success,
       COALESCE(sum(case when error > 0 then 1 else 0 end), 0) as task_error,
       COALESCE(sum(count), 0) as url_total,
       COALESCE(sum(count) - sum(success) - sum(error), 0) as url_queue,
       COALESCE(sum(success), 0) as url_success,
       COALESCE(sum(error), 0) as url_error
       from task`).Scan(&status.Tasks.Total,
		&status.Tasks.InQueue,
		&status.Tasks.Complete,
		&status.Tasks.Error,
		&status.Urls.Total,
		&status.Urls.InQueue,
		&status.Urls.Complete,
		&status.Urls.Error)
	return
}
