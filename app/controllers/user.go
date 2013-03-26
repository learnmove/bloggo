package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/jgraham909/bloggo/app/models"
	"github.com/robfig/revel"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Application
}

func (c User) Index() revel.Result {
	if c.User != nil {
		action := "/User/SaveExistingUser"
		user := c.User
		// TODO Populate form & Save properly
		return c.Render(user, action)
	}
	return c.Redirect(User.Login)
}

func (c User) SaveExistingUser(user *models.User, verifyPassword string) revel.Result {
	return c.SaveUser(user)
}

func (c User) SaveNewUser(user *models.User, verifyPassword string) revel.Result {
	if exists := user.GetUserByEmail(c.MSession, user.Email); exists.Email == user.Email {
		msg := fmt.Sprint("Account with ", user.Email, " already exists.")
		c.Validation.Required(user.Email != exists.Email).
			Message(msg)
	} else {
		user.Id = bson.NewObjectId()
	}

	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).
		Message("Password does not match")

	return c.SaveUser(user)
}

func (c User) SaveUser(user *models.User) revel.Result {
	fmt.Printf("SaveUser(user): %v\n", user)

	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Please correct the errors below.")
		return c.Redirect(User.RegisterForm)
	}

	user.Save(c.MSession)

	c.Session["user"] = user.Email
	c.Flash.Success("Welcome, " + user.String())
	return c.Redirect(Application.Index)
}

func (c User) Login(Email, Password string) revel.Result {
	user := new(models.User)
	user = user.GetUserByEmail(c.MSession, Email)

	if user.Email != "" {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(Password))
		if err == nil {
			c.Session["user"] = Email
			c.Flash.Success("Welcome, " + Email)
			return c.Redirect(Application.Index)
		}
	}

	c.Flash.Out["mail"] = Email
	c.Flash.Error("Incorrect email address or password.")
	return c.Redirect(User.LoginForm)
}

func (c User) LoginForm() revel.Result {
	if c.UserAuthenticated() == false {
		return c.Render()
	}

	// User already logged in bounce to main page
	return c.Redirect(Application.Index)
}

func (c User) RegisterForm() revel.Result {
	action := "/User/SaveNewUser"
	return c.Render(action)
}

func (c User) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(Application.Index)
}
