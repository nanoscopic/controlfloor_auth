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
    UserLogin( c *gin.Context ) bool
}

type SessionManager interface {
    GetSCSSessionManager() *scs.SessionManager
    GetSession( c *gin.Context ) context.Context
    WriteSession( c *gin.Context )
}

type authUser struct {
    userName string
    password string
}

type demoAH struct {
    sessionManager SessionManager
    testParam      string
    testUser       string
    users          []authUser
}

func NewAuthHandler( authNode uj.JNode, sessionManager SessionManager ) AuthHandler {
    self := &demoAH{
        sessionManager: sessionManager,
        testParam: "",
        testUser: "test",
        users: []authUser{},
    }
    
    if authNode != nil {
        testParamNode := authNode.Get("testparam")
        if testParamNode != nil {
            self.testParam = testParamNode.String()
        } else {
            fmt.Printf("Missing conf auth.testparam node\n")
        }
        testUserNode := authNode.Get("testuser")
        if testUserNode != nil {
            self.testUser = testUserNode.String()
        } else {
            fmt.Printf("Missing conf auth.testuser node\n")
        }
        
        usersNode := authNode.Get("users")
        if usersNode != nil {
            users := []authUser{}
            usersNode.ForEach( func( user uj.JNode ) {
                userNameNode := user.Get("userName")
                pwNode := user.Get("password")
                if userNameNode != nil && pwNode != nil {
                    user := authUser{
                        userName: userNameNode.String(),
                        password: pwNode.String(),
                    }
                    users = append( users, user )
                }
            } )
            self.users = users
        } else {
            fmt.Printf("Missing conf auth.users node\n")
        }
    } else {
        fmt.Printf("Missing conf auth node\n")
    }
    
    return self
}

func (self *demoAH) UserAuth( c *gin.Context ) bool {
    authPass := false
    
    fmt.Printf("uauth\n")
    
    if self.testParam != "" {
        fmt.Printf("checking for %s\n", self.testParam)
        _, uok := c.GetQuery( self.testParam )
        if uok {
            sm := self.sessionManager
            s := sm.GetSession(c)
            scsSM := sm.GetSCSSessionManager()
            scsSM.Put( s, "user", self.testUser )
            sm.WriteSession( c )
            authPass = true
        }
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

func (self *demoAH) UserLogin( c *gin.Context ) bool {
    s := self.sessionManager.GetSession( c )
    scsSM := self.sessionManager.GetSCSSessionManager()
    
    user := c.PostForm("user")
    pass := c.PostForm("pass")
    
    for _,u := range self.users {
        if user == u.userName && pass == u.password {
            fmt.Printf( "login ok; user=%s\n", u.userName )
            
            scsSM.Put( s, "user", u.userName )
            self.sessionManager.WriteSession( c )
            
            return true
        }
    }
    
    fmt.Printf( "login failed; user=%s\n", user )
    return false
}
