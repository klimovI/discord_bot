package commands

var testCommand = Command{
	Data: CommandData{
		Name:        "test",
		Description: "Test command",
	},
	Handler: func(ctx Context) {
		ctx.Respond("test")
	},
}
