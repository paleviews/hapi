package logic

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/paleviews/hapi/example/todo/apidesign/golang/codes"
	"github.com/paleviews/hapi/example/todo/apidesign/golang/common"
	"github.com/paleviews/hapi/example/todo/apidesign/golang/todo"
)

type Todo struct {
	mu           sync.Mutex
	todoByID     map[int64]*todo.Todo
	deletedTodos map[int64]*todo.Todo
	lastID       int64
}

func NewTodo() todo.V1Server {
	return &Todo{
		todoByID:     make(map[int64]*todo.Todo),
		deletedTodos: make(map[int64]*todo.Todo),
	}
}

func (t *Todo) get(id int64) (*todo.Todo, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	td, ok := t.todoByID[id]
	return td, ok
}

func (t *Todo) set(td *todo.Todo) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.todoByID[td.ID] = td
}

func (t *Todo) Create(_ context.Context, req *todo.CreateRequest) (*todo.CreateResponse, error) {
	id := atomic.AddInt64(&t.lastID, 1)
	t.set(&todo.Todo{
		ID:            id,
		Title:         req.Title,
		Detail:        req.Detail,
		Labels:        nil,
		Completeness:  0,
		CreatedTime:   time.Now().Unix(),
		CompletedTime: 0,
	})
	return &todo.CreateResponse{
		ID: id,
	}, nil
}

func (t *Todo) Get(_ context.Context, req *todo.GetRequest) (*todo.Todo, error) {
	td, ok := t.get(req.ID)
	if !ok {
		return nil, codes.APIErrorFromResponseCode(
			codes.ResponseCode_RESPONSE_CODE_NOT_FOUND,
			fmt.Errorf("todo by id %d not found", req.ID),
		)
	}
	return td, nil
}

func (t *Todo) List(_ context.Context, req *todo.ListRequest) (*todo.ListResponse, error) {
	filter := func(td *todo.Todo) bool {
		if req.TitleContains != "" {
			if !strings.Contains(td.Title, req.TitleContains) {
				return false
			}
		}
		if req.DetailContains != "" {
			if !strings.Contains(td.Detail, req.DetailContains) {
				return false
			}
		}
		return true
	}
	var list []*todo.Todo
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, v := range t.todoByID {
		if filter(v) {
			list = append(list, v)
		}
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})
	start := req.Page * req.PageSize
	end := start + req.PageSize
	total := int64(len(list))
	if start >= total {
		return &todo.ListResponse{
			Total: total,
			List:  nil,
		}, nil
	}
	if end >= total {
		return &todo.ListResponse{
			Total: total,
			List:  list[start:],
		}, nil
	}
	return &todo.ListResponse{
		Total: total,
		List:  list[start:end],
	}, nil
}

func (t *Todo) Update(_ context.Context, req *todo.Todo) (*common.Empty, error) {
	if _, ok := t.todoByID[req.ID]; !ok {
		return nil, codes.APIErrorFromResponseCode(
			codes.ResponseCode_RESPONSE_CODE_NOT_FOUND,
			fmt.Errorf("todo id %d not found", req.ID),
		)
	}
	t.set(req)
	return &common.Empty{}, nil
}

func (t *Todo) Delete(_ context.Context, req *todo.DeleteRequest) (*common.Empty, error) {
	td, ok := t.get(req.ID)
	if !ok {
		return nil, codes.APIErrorFromResponseCode(
			codes.ResponseCode_RESPONSE_CODE_NOT_FOUND,
			fmt.Errorf("todo id %d not found", req.ID),
		)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.todoByID, req.ID)
	if req.SoftDelete {
		t.deletedTodos[td.ID] = td
	}
	return &common.Empty{}, nil
}
