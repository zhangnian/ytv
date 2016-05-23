// GENERATED CODE - DO NOT EDIT
package routes

import "github.com/revel/revel"


type tApiBaseController struct {}
var ApiBaseController tApiBaseController


func (_ tApiBaseController) RenderError(
		code int,
		msg string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "code", code)
	revel.Unbind(args, "msg", msg)
	return revel.MainRouter.Reverse("ApiBaseController.RenderError", args).Url
}


type tTestRunner struct {}
var TestRunner tTestRunner


func (_ tTestRunner) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.Index", args).Url
}

func (_ tTestRunner) Run(
		suite string,
		test string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	revel.Unbind(args, "test", test)
	return revel.MainRouter.Reverse("TestRunner.Run", args).Url
}

func (_ tTestRunner) List(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.List", args).Url
}


type tStatic struct {}
var Static tStatic


func (_ tStatic) Serve(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.Serve", args).Url
}

func (_ tStatic) ServeModule(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModule", args).Url
}


type tApiInfoController struct {}
var ApiInfoController tApiInfoController


func (_ tApiInfoController) Announcement(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("ApiInfoController.Announcement", args).Url
}


type tApiUserController struct {}
var ApiUserController tApiUserController


func (_ tApiUserController) Login(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("ApiUserController.Login", args).Url
}

func (_ tApiUserController) Register(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("ApiUserController.Register", args).Url
}

func (_ tApiUserController) GetCode(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("ApiUserController.GetCode", args).Url
}

func (_ tApiUserController) CheckCode(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("ApiUserController.CheckCode", args).Url
}


