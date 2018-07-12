package main

import (
	"fmt"
	"os"

	"github.com/facebookgo/inject"
	"github.com/fatih/color"
	"github.com/srgrn/postmanerator/commands"
	"github.com/srgrn/postmanerator/configuration"
	"github.com/srgrn/postmanerator/postman"
	"github.com/srgrn/postmanerator/themes"
	"github.com/srgrn/postmanerator/utils"
)

var (
	config             = configuration.Config
	errUnknownCmd      = fmt.Errorf("Command not found, please see the documentation at https://github.com/aubm/postmanerator")
	themeManager       = &themes.Manager{}
	themeRenderer      = &themes.Renderer{}
	gitAgent           = &utils.GitAgent{}
	collectionBuilder  = &postman.CollectionBuilder{}
	collectionV1Parser = &postman.CollectionV1Parser{}
	environmentBuilder = &postman.EnvironmentBuilder{}
	defaultCommand     = &commands.Default{}
	getThemeCommand    = &commands.GetTheme{}
	deleteThemeCommand = &commands.DeleteTheme{}
	listThemesCommand  = &commands.ListThemes{}
	availableCommands  = []commands.Command{}
)

func init() {
	checkAndPrintErr(_init())
}

func _init() error {
	configuration.Init()
	if err := inject.Populate(config, themeManager, defaultCommand, getThemeCommand, deleteThemeCommand,
		listThemesCommand, gitAgent, themeRenderer, collectionBuilder, collectionV1Parser, environmentBuilder); err != nil {
		return fmt.Errorf("app initialization failed: %v", err)
	}
	collectionBuilder.Parsers = append(collectionBuilder.Parsers, collectionV1Parser)
	availableCommands = append(availableCommands,
		defaultCommand,
		getThemeCommand,
		deleteThemeCommand,
		listThemesCommand,
	)
	return nil
}

func main() {
	checkAndPrintErr(configuration.InitErr)
	checkAndPrintErr(_main())
}

func _main() (err error) {
	userCommand := evaluateUserCommand()
	for _, availableCommand := range availableCommands {
		if availableCommand.Is(userCommand) {
			return availableCommand.Do()
		}
	}
	return errUnknownCmd
}

func evaluateUserCommand() string {
	if len(config.Args) == 0 {
		return commands.CmdDefault
	}

	switch config.Args[0] {
	case "themes":
		if len(config.Args) < 2 {
			return commands.CmdThemesList
		}
		switch config.Args[1] {
		case "get":
			return commands.CmdThemesGet
		case "delete":
			return commands.CmdThemesDelete
		case "list":
			return commands.CmdThemesList
		}
	}

	return commands.CmdUnknown
}

func checkAndPrintErr(err error) {
	if err == nil {
		return
	}

	fmt.Println(color.RedString(err.Error()))
	os.Exit(1)
}
