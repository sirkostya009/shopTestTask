package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"shopTestTask/db"
)

type Handlers struct {
	db db.DB
}

func New() *Handlers {
	return &Handlers{db.New()}
}

type httpErr struct {
	code int
	msg  string
}

func (err *httpErr) Error() string {
	return fmt.Sprintf("%d: %s", err.code, err.msg)
}

type Context struct {
	context.Context
	http.ResponseWriter
	*http.Request
}

func (ctx *Context) JSON(i any) error {
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(ctx.ResponseWriter).Encode(i)
}

type CtxHandler func(ctx *Context) error

func (f CtxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := Context{r.Context(), w, r}
	err := f(&ctx)
	if err != nil {
		var httpE *httpErr
		if errors.As(err, &httpE) {
			http.Error(w, httpE.msg, httpE.code)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
