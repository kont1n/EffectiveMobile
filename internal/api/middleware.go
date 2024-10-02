package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/urfave/negroni"
)

const (
	defaultSortBy    = "id"
	defaultSortOrder = "desc"
	defaultLimit     = "10"
	defaultPageToken = "0"
	defaultСouplet   = "0"
)

func (api *ApiHandler) LogAPI(h http.Handler) http.Handler {
	logFn := func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		uri := request.RequestURI
		method := request.Method
		reqID := middleware.GetReqID(request.Context())
		lrw := negroni.NewResponseWriter(writer)

		h.ServeHTTP(lrw, request)

		statusCode := lrw.Status()
		duration := time.Since(start)

		api.loger.Debugln(
			"RequestID:", reqID,
			"statusCode:", statusCode,
			"uri:", uri,
			"method:", method,
			"duration:", duration,
		)
	}
	return http.HandlerFunc(logFn)
}

func Sorting(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		sortBy := request.URL.Query().Get("sort_by")
		sortOrder := request.URL.Query().Get("sort_order")

		if sortBy == "" {
			sortBy = defaultSortBy
		}
		if sortOrder == "" {
			sortOrder = defaultSortOrder
		}

		if strings.ToLower(sortOrder) != "asc" && strings.ToLower(sortOrder) != "desc" {
			http.Error(writer, "sort_order must be asc or desc", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(request.Context(), "sort_by", sortBy)
		ctx = context.WithValue(request.Context(), "sort_order", sortOrder)

		handler.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func Pagination(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		limit := request.URL.Query().Get("limit")
		pageToken := request.URL.Query().Get("page_token")

		if limit == "" {
			limit = defaultLimit
		}
		if pageToken == "" {
			pageToken = defaultPageToken
		}

		ctx := context.WithValue(request.Context(), "limit", limit)
		ctx = context.WithValue(request.Context(), "page_token", pageToken)

		handler.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func FiledPagination(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		couplet := request.URL.Query().Get("couplet")

		if couplet == "" {
			couplet = defaultСouplet
		}

		ctx := context.WithValue(request.Context(), "couplet", couplet)

		handler.ServeHTTP(writer, request.WithContext(ctx))
	})
}
