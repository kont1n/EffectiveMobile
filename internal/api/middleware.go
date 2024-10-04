package api

import (
	"EffectiveMobile/internal/models"
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
	defaultOffset    = "0"
	defaultPageToken = "0"
	defaultСouplet   = "0"
)

var (
	sortList = map[string]string{
		"song":    "",
		"group":   "",
		"release": "",
		"text":    "",
		"link":    "",
	}
	filterList = map[string]string{
		"song":    "",
		"group":   "",
		"release": "",
		"text":    "",
		"link":    "",
	}
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

func (api *ApiHandler) Sorting(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		api.loger.Debugln("Sorting Middleware")

		reqID := middleware.GetReqID(request.Context())
		sortBy := request.URL.Query().Get("sort_by")
		sortOrder := request.URL.Query().Get("sort_order")

		if sortBy != "" {
			if _, ok := sortList[sortBy]; !ok {
				api.JSONError(writer, "Sort_by is not valid", http.StatusBadRequest, reqID)
				return
			}
		} else {
			sortBy = defaultSortBy
		}

		if sortOrder != "" {
			if strings.ToLower(sortOrder) != "asc" && strings.ToLower(sortOrder) != "desc" {
				api.JSONError(writer, "Sort_order must be asc or desc", http.StatusBadRequest, reqID)
				return
			}
		} else {
			sortOrder = defaultSortOrder
		}

		sortOptions := models.SortOptions{
			Field: sortBy,
			Order: sortOrder,
		}

		ctx := context.WithValue(request.Context(), "sort_options", sortOptions)
		handler.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (api *ApiHandler) Filtering(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		api.loger.Debugln("Filtering Middleware")

		m := make(map[string]string)
		for key, val := range filterList {
			m[key] = val
		}

		for key, _ := range m {
			val := request.URL.Query().Get(key)
			if val != "" {
				m[key] = val
			} else {
				delete(m, key)
			}
		}

		ctx := context.WithValue(request.Context(), "filter_options", m)

		handler.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (api *ApiHandler) Pagination(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		api.loger.Debugln("Pagination Middleware")

		limit := request.URL.Query().Get("limit")
		offset := request.URL.Query().Get("offset")
		pageToken := request.URL.Query().Get("page_token")

		if limit == "" {
			limit = defaultLimit
		}
		if pageToken == "" && offset == "" {
			offset = defaultOffset
		}

		paginationOptions := models.PaginationOptions{
			Limit:     limit,
			Offset:    offset,
			PageToken: pageToken,
		}
		ctx := context.WithValue(request.Context(), "pagination_options", paginationOptions)
		handler.ServeHTTP(writer, request.WithContext(ctx))
	})
}
