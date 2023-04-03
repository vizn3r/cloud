package handler

// ---------------
// ROUTE
// ---------------

type route struct {
	path 			string
	funcs 			map[string][]RouterFunc
	middlewares 	[]RouterFunc
}

func (r* route) Path() string {
	return r.path
}

func (r* route) Use(fns ...RouterFunc) {
	r.middlewares = append(r.middlewares, fns...)
}

