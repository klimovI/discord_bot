package commands

var helloCommand = Command{
	Data: CommandData{
		Name:        "hello",
		Description: "Says hello",
	},
	Handler: func(ctx Context) {
		ctx.Respond("Hello")
	},
}
