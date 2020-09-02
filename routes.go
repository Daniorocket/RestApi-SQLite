package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"UserinfoIndex",
		"GET",
		"/todos",
		UserinfoIndex,
	},
	Route{
		"UserinfoShow",
		"GET",
		"/todos/{todoId}",
		UserinfoShow,
	},
	Route{
		"UserinfoCreate",
		"POST",
		"/todos",
		UserinfoCreate,
	},
	Route{
		"EditUserinfo",
		"PUT",
		"/todos/edit/{uid}",
		EditUserinfo,
	},
	Route{
		"DeleteUserinfo",
		"DELETE",
		"/todos/delete/{uid}",
		DeleteUserinfo,
	},
}
