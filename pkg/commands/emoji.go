package commands

var emojiCommand = Command{
	Data: CommandData{
		Name:        "emoji",
		Description: "Says hello",
	},
	Handler: func(ctx Context) {
		ctx.Respond("🤡")
	},
}
