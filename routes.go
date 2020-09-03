package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func initRoutes(handler Handler) Routes {
	return Routes{
		Route{
			"Index",
			"GET",
			"/",
			handler.Index,
		},
		Route{
			"UserinfoIndex",
			"GET",
			"/todos",
			handler.UserinfoIndex,
		},
		Route{
			"UserinfoShow",
			"GET",
			"/todos/{todoId}",
			handler.UserinfoShow,
		},
		Route{
			"UserinfoCreate",
			"POST",
			"/todos",
			handler.UserinfoCreate,
		},
		Route{
			"EditUserinfo",
			"PUT",
			"/todos/edit/{uid}",
			handler.EditUserinfo,
		},
		Route{
			"DeleteUserinfo",
			"DELETE",
			"/todos/delete/{uid}",
			handler.DeleteUserinfo,
		},
	}
}
