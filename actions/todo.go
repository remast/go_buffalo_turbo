package actions

import (
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/validate"
)

var todoItems []*TodoItem

func init() {
	todoItems = []*TodoItem{}
}

// TodoItem for a single todo
type TodoItem struct {
	ID   int
	Name string
	Done bool
}

// TodoIndex default implementation.
func TodoIndex(c buffalo.Context) error {
	c.Set("todoItems", todoItems)
	return c.Render(http.StatusOK, r.HTML("todo/index.plush.html"))
}

// TodoToggle default implementation.
func TodoToggle(c buffalo.Context) error {
	todoItemID := c.Params().Get("ID")

	for _, todoItem := range todoItems {
		if strconv.Itoa(todoItem.ID) == todoItemID {
			todoItem.Done = !todoItem.Done
			c.Set("todoItem", todoItem)
			break
		}
	}

	if acceptsTurboStream(c.Request()) {
		id := "todo_item_" + todoItemID
		return c.Render(http.StatusOK, r.Func("text/vnd.turbo-stream.html", createTurboWriter("todo/todo_item.plush.html", "replace", id)))
	}

	return c.Redirect(302, "/")
}

// TodoCreate default implementation.
func TodoCreate(c buffalo.Context) error {
	newTodoItem := &TodoItem{ID: rand.Intn(1000)}
	if err := c.Bind(newTodoItem); err != nil {
		return err
	}
	if newTodoItem.Name == "" {
		c.Set("todoItem", newTodoItem)

		verrs := validate.NewErrors()
		verrs.Add("name", "Name missing.")
		c.Set("errors", verrs)

		if acceptsTurboStream(c.Request()) {
			turboAction := "replace"
			turboDomID := "todo_new_form"
			return c.Render(http.StatusOK, r.Func("text/vnd.turbo-stream.html", createTurboWriter("todo/new.plush.html", turboAction, turboDomID)))
		}

		return c.Render(http.StatusOK, r.HTML("todo/new.plush.html"))
	}

	todoItems = append(todoItems, newTodoItem)
	return c.Redirect(302, "/")
}

// TodoNew default implementation.
func TodoNew(c buffalo.Context) error {
	c.Set("todoItem", &TodoItem{})
	return c.Render(http.StatusOK, r.HTML("todo/new.plush.html"))
}

func acceptsTurboStream(request *http.Request) bool {
	for _, acceptValue := range request.Header["Accept"] {
		if strings.Contains(acceptValue, "text/vnd.turbo-stream.html") {
			return true
		}
	}
	return false
}

func createTurboWriter(template, action, target string) render.RendererFunc {
	return func(w io.Writer, d render.Data) error {
		d["action"] = action
		d["target"] = target
		r.HTML(template, "turbo/turbo_stream.plush.html").Render(w, d)
		return nil
	}
}
