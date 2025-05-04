package main

import (
	"to-read/controllers/auth"
	"to-read/model"
	"to-read/shared/llmprocessor"
	"to-read/shared/server"
	"to-read/shared/yamlconfig"
)

func main() {
	configuration, err := yamlconfig.ConfigLoad("config.yml")
	if err != nil {
		panic(err)
	}

	err = model.Connect(configuration.Database)
	if err != nil {
		panic(err)
	}

	err = model.InitModel()
	if err != nil {
		panic(err)
	}

	err = auth.InitAuthorization(configuration.Authorization)
	if err != nil {
		panic(err)
	}

	err = llmprocessor.InitLLMProcessor(configuration.LLMProcessor)
	if err != nil {
		panic(err)
	}

	// err = miniprogram.InitMiniProgramConfig(configuration.Miniprogram)
	// if err != nil {
	// 	panic(err)
	// }

	err = server.Run(configuration.Server)
	if err != nil {
		panic(err)
	}
}
