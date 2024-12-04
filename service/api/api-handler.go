package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	//USER ENDPOINT
	rt.router.POST("/user/session", rt.createUserHandler)
	rt.router.GET("/users", rt.getAllUsersHandler)
	rt.router.POST("/user", rt.updateUsernameHandler)

	//CONVERSATION ENDPOINT
	rt.router.POST("/conversation", rt.setConversationHandler)

	return rt.router
}
