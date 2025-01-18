package utils

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

// getHTTPVerbFromMethodName infers the HTTP verb from a method's name.
func getHTTPVerbFromMethodName(methodName string) string {
	// Define a mapping of prefixes to HTTP verbs
	methodToVerb := map[string]string{
		"Get":    "GET",
		"Create": "POST",
		"Post":   "POST",
		"Update": "PUT",
		"Delete": "DELETE",
		"Patch":  "PATCH",
		"List":   "GET",
	}

	// Check the prefix of the methodName
	for prefix, verb := range methodToVerb {
		if strings.HasPrefix(methodName, prefix) {
			return verb
		}
	}

	// Default to GET if no match is found, or return an empty string/appropriate fallback
	return "GET"
}

// buildRoutePath determines the route path from a method's name and optionally adds route parameters.
func buildRoutePath(pathPrefix, methodName, httpVerb string) (string, string) {
	var sb strings.Builder
	routePath := ""
	routeParam := ""
	extractedRouteName := ""
	if strings.Contains(strings.ToLower(methodName), "by") {
		// split the function name by the word 'by' and use the second part as param
		splittedName := strings.Split(strings.ToLower(methodName), "by")
		extractedRouteName = splittedName[0]
		routeParam = strings.ToLower(splittedName[1])
	} else {
		extractedRouteName = methodName
	}
	// Strip HTTP verb prefix from the method name
	verbs := []string{"Get", "Delete", "Update", "Create", "Patch", "List", "Post"}
	for _, verb := range verbs {
		if strings.HasPrefix(strings.ToLower(extractedRouteName), strings.ToLower(verb)) {
			extractedRouteName = extractedRouteName[len(verb):] // Remove the prefix
			break
		}
	}

	// Build the route path by adding slashes before uppercase letters
	for i, r := range extractedRouteName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			sb.WriteRune('-')
		}
		sb.WriteRune(r)
	}
	routePath = pathPrefix + "/" + strings.ToLower(sb.String()) + routePath
	if routeParam != "" {
		if routeParam == "id" {
			routePath += fmt.Sprintf("/:%s", routeParam)
		}
	}
	// Combine the processed path and convert it to lowercase

	return routePath, routeParam
}

func callHandlerMethod(v reflect.Value, c echo.Context, method reflect.Method, routeParam string) error {
	// Check if the route includes a parameter
	paramValue := c.Param(routeParam)

	var args []reflect.Value
	args = append(args, v)                  // Add the handler instance as the first argument
	args = append(args, reflect.ValueOf(c)) // Add the Echo context (`c`) as the second argument

	// If a parameter is present in the route, pass it as an argument to the method
	if paramValue != "" {
		args = append(args, reflect.ValueOf(paramValue))
	}

	// Call the handler method with the appropriate arguments
	results := method.Func.Call(args)

	// Handle any returned results
	if len(results) > 0 && results[0].CanInterface() {
		if results[0].Interface() != nil {
			return results[0].Interface().(error)
		}
		return nil
	}
	return nil
}

// RegisterRoutesAutomatically registers routes dynamically based on the methods in the handler.
func RegisterRoutesAutomatically(e *echo.Echo, handler interface{}, methodPrefix string, log *logrus.Logger, isPrivate bool, authFunc echo.MiddlewareFunc) {
	t := reflect.TypeOf(handler)
	v := reflect.ValueOf(handler)

	// Set up group and ensure proper auth setup for private routes
	var echoGroup *echo.Group
	if isPrivate && authFunc != nil {
		echoGroup = e.Group("/private", authFunc)
	} else {
		echoGroup = e.Group("") // No additional prefix for public routes
	}

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		httpVerb := getHTTPVerbFromMethodName(method.Name)
		routePath, routeParam := buildRoutePath(methodPrefix, method.Name, httpVerb)

		log.Debug(fmt.Sprintf("Registering route: %s %s (private: %t)", httpVerb, routePath, isPrivate))

		if isPrivate {
			echoGroup.Add(httpVerb, routePath, func(c echo.Context) error {
				return callHandlerMethod(v, c, method, routeParam)
			})
		} else {
			e.Add(httpVerb, routePath, func(c echo.Context) error {
				return callHandlerMethod(v, c, method, routeParam)
			})
		}
	}
}
