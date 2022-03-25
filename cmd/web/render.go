package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type templateData struct {
	Data map[string]interface{}
}

func (app *application) addDefaultTData(td *templateData) *templateData {
	// If we did not have data, create data.
	if td == nil {
		var newTD templateData
		td = &newTD
	}
	data := make(map[string]interface{})

	if len(td.Data) == 0 {
		td.Data = data
	}

	if app.vueDist != nil {
		td.Data["vueGlue"] = app.VueGlue
	}
	return td
}

func (app *application) RenderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData) {
	templateToLoad := fmt.Sprintf("templates/%s.page.gohtml", page)
	td = app.addDefaultTData(td)

	ts, err := template.ParseFiles(templateToLoad)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// We then use the Execute() method on the template set to write the template
	// content as the response body. The last parameter to Execute() represents any
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	err = ts.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}
