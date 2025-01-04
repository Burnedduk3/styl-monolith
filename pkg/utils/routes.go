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
func buildRoutePath(pathPrefix, methodName string) (string, string) {
	var sb strings.Builder
	routePath := ""
	routeParam := ""
	extractedRouteName := ""
	if strings.Contains(strings.ToLower(methodName), "by") {
		// split the function name by the word 'by' and use the second part as param
		splittedName := strings.Split(strings.ToLower(methodName), "by")
		extractedRouteName = strings.ToLower(splittedName[0])
		routeParam = strings.ToLower(splittedName[1])
	} else {
		extractedRouteName = strings.ToLower(methodName)
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
		routePath += fmt.Sprintf("/%s/:%s", routeParam, routeParam)
	}
	// Combine the processed path and convert it to lowercase

	return routePath, routeParam
}

// RegisterRoutesAutomatically registers routes dynamically based on the methods in the handler.
func RegisterRoutesAutomatically(e *echo.Echo, handler interface{}, methodPrefix string, log *logrus.Logger) {
	t := reflect.TypeOf(handler)
	v := reflect.ValueOf(handler)

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		routePath, routeParam := buildRoutePath(methodPrefix, method.Name)
		httpVerb := getHTTPVerbFromMethodName(method.Name)
		log.Debug(fmt.Sprintf("http_verb: %s route_path: %s", httpVerb, routePath))

		// Dynamic route registration
		e.Add(httpVerb, routePath, func(c echo.Context) error {
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
		})
	}
}
