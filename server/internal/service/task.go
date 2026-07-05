package service

import (
	"context"
	"errors"
	"time"

	"starledger/ent"
	"starledger/ent/task"
)

type TaskService struct {
	client *ent.Client
}

func NewTaskService(client *ent.Client) *TaskService {
	return &TaskService{client: client}
}

// List returns tasks for a tenant with pagination.
func (s *TaskService) List(ctx context.Context, tenantID, page, pageSize int, assigneeID int) ([]*ent.Task, int, error) {
	query := s.client.Task.Query().Where(task.TenantID(tenantID))

	// Filter by assignee if specified
	if assigneeID > 0 {
		query = query.Where(task.AssigneeID(assigneeID))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	items, err := query.
		WithAssignee().
		WithCreator().
		Order(ent.Desc(task.FieldCreatedAt)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Get returns a task by ID.
func (s *TaskService) Get(ctx context.Context, id, tenantID int) (*ent.Task, error) {
	return s.client.Task.Query().
		Where(task.ID(id), task.TenantID(tenantID)).
		WithAssignee().
		WithCreator().
		Only(ctx)
}

// Create creates a new task.
func (s *TaskService) Create(ctx context.Context, tenantID, creatorID int, title, description string, assigneeID int, priority string, dueDate *time.Time) (*ent.Task, error) {
	builder := s.client.Task.Create().
		SetTenantID(tenantID).
		SetTitle(title).
		SetDescription(description).
		SetCreatorID(creatorID).
		SetStatus(task.StatusPending)

	if assigneeID > 0 {
		builder.SetAssigneeID(assigneeID)
	}
	if priority != "" {
		builder.SetPriority(task.Priority(priority))
	} else {
		builder.SetPriority(task.PriorityMedium)
	}
	if dueDate != nil {
		builder.SetDueDate(*dueDate)
	}

	return builder.Save(ctx)
}

// Update updates a task.
func (s *TaskService) Update(ctx context.Context, id, tenantID int, title, description string, assigneeID int, status, priority string, dueDate *time.Time) (*ent.Task, error) {
	t, err := s.client.Task.Query().
		Where(task.ID(id), task.TenantID(tenantID)).
		Only(ctx)
	if err != nil {
		return nil, errors.New("任务不存在")
	}

	updater := t.Update()
	if title != "" {
		updater.SetTitle(title)
	}
	if description != "" {
		updater.SetDescription(description)
	}
	if assigneeID > 0 {
		updater.SetAssigneeID(assigneeID)
	}
	if status != "" {
		updater.SetStatus(task.Status(status))
	}
	if priority != "" {
		updater.SetPriority(task.Priority(priority))
	}
	if dueDate != nil {
		updater.SetDueDate(*dueDate)
	}

	return updater.Save(ctx)
}

// Delete deletes a task.
func (s *TaskService) Delete(ctx context.Context, id, tenantID int) error {
	n, err := s.client.Task.Delete().
		Where(task.ID(id), task.TenantID(tenantID)).
		Exec(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("任务不存在")
	}
	return nil
}

// Assign assigns a task to a user.
func (s *TaskService) Assign(ctx context.Context, id, tenantID, assigneeID int) (*ent.Task, error) {
	t, err := s.client.Task.Query().
		Where(task.ID(id), task.TenantID(tenantID)).
		Only(ctx)
	if err != nil {
		return nil, errors.New("任务不存在")
	}

	return t.Update().
		SetAssigneeID(assigneeID).
		SetStatus(task.StatusInProgress).
		Save(ctx)
}

// MyTasks returns tasks assigned to a specific user.
func (s *TaskService) MyTasks(ctx context.Context, tenantID, userID int) ([]*ent.Task, error) {
	return s.client.Task.Query().
		Where(
			task.TenantID(tenantID),
			task.AssigneeID(userID),
			task.StatusNEQ(task.StatusCompleted),
			task.StatusNEQ(task.StatusCancelled),
		).
		WithCreator().
		Order(ent.Asc(task.FieldDueDate)).
		All(ctx)
}
