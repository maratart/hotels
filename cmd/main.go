package main

import "hotels/app"

func main() {
	a := app.NewApp(app.NewConfig())
	a.Run()
}
