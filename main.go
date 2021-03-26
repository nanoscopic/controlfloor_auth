package controlfloor_auth

import (
    "context"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/alexedwards/scs/v2"
    uj "github.com/nanoscopic/ujsonin/v2/mod"
)

type AuthHandler interface {
    UserAuth( c *gin.Context ) bool
}

type SessionManager interface {
    GetSCSSessionManager() *scs.SessionManager
    GetSession( c *gin.Context ) context.Context
    WriteSession( c *gin.Context )
}

type demoAH struct {
    sessionManager SessionManager
}

func NewAuthHandler( confRoot uj.JNode, sessionManager SessionManager ) AuthHandler {
    return &demoAH{ sessionManager }
}

func (self *demoAH) UserAuth( c *gin.Context ) bool {
    authPass := false
    
    _, uok := c.GetQuery("ok")
    if uok {
        sm := self.sessionManager
        s := sm.GetSession(c)
        scsSM := sm.GetSCSSessionManager()
        scsSM.Put( s, "user", "test" )
        sm.WriteSession( c )
        authPass = true
    }
    
    if authPass {
        c.Next()
        return true
    }
    
    c.Redirect( 302, "/login" )
    c.Abort()
    fmt.Println("user fail")
    
    return false
}

