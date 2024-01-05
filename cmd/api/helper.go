package main

import "fmt"

func (app *application) background(fn func()) {
	app.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Sprintf("%v", err))
			}
		}()
		fn()
	}()
}
