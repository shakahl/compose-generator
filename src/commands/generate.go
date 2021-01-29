package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/otiai10/copy"

	"compose-generator/parser"
	"compose-generator/utils"
)

func Generate(flag_advanced bool, flag_run bool, flag_demonized bool) {
	// Execute SafetyFileChecks
	utils.ExecuteSafetyFileChecks()

	// Welcome Message
	utils.Heading("Welcome to Compose Generator!")
	fmt.Println("Please continue by answering a few questions:")
	fmt.Println()

	// Project name
	project_name := utils.TextQuestion("What is the name of your project: ")
	if project_name == "" {
		utils.Error("Error. You must specify a project name!", true)
	}
	project_name_container := strings.ReplaceAll(strings.ToLower(project_name), " ", "-")

	// Docker Swarm compatability (default: no)
	docker_swarm := utils.YesNoQuestion("Should your compose file be used for distributed deployment with Docker Swarm?", false)
	fmt.Println(docker_swarm)

	// Predefined stack (default: yes)
	use_predefined_stack := utils.YesNoQuestion("Do you want to use a predefined stack?", true)
	if use_predefined_stack {
		// Load stacks from templates
		template_data := parser.ParseTemplates()
		// Predefined stack menu
		var items []string
		for _, t := range template_data {
			items = append(items, t.Label)
		}
		index, _ := utils.MenuQuestion("Predefined software stack", items)
		fmt.Println()

		// Ask configured questions to the user
		envMap := make(map[string]string)
		envMap["PROJECT_NAME"] = project_name
		envMap["PROJECT_NAME_CONTAINER"] = project_name_container
		for _, q := range template_data[index].Questions {
			if !q.Advanced || (q.Advanced && flag_advanced) {
				switch q.Type {
				case 1: // Yes/No
					default_value, _ := strconv.ParseBool(q.Default_value)
					envMap[q.Env_var] = strconv.FormatBool(utils.YesNoQuestion(q.Text, default_value))
				case 2: // Text
					envMap[q.Env_var] = utils.TextQuestionWithDefault(q.Text, q.Default_value)
				}
			} else {
				envMap[q.Env_var] = q.Default_value
			}
		}
		fmt.Println()

		// Copy templates
		fmt.Print("Copying template ...")
		src_path := utils.GetTemplatesPath() + "/" + template_data[index].Dir
		dst_path := "."

		os.Remove(dst_path + "/docker-compose.yml")
		os.Remove(dst_path + "/environment.env")
		os.Remove(dst_path + "/volumes")

		opt := copy.Options{
			Skip: func(src string) (bool, error) {
				return strings.HasSuffix(src, "config.json") || strings.HasSuffix(src, "README.md") || strings.HasSuffix(src, ".gitkeep"), nil
			},
			OnDirExists: func(src string, dst string) copy.DirExistsAction {
				return copy.Replace
			},
		}
		err := copy.Copy(src_path, dst_path, opt)
		if err != nil {
			utils.Error("Could not copy template files.", true)
		}
		color.Green(" [done]")

		// Replace variables
		fmt.Print("Applying customizations ...")
		utils.ReplaceVarsInFile("./docker-compose.yml", envMap)
		utils.ReplaceVarsInFile("./environment.env", envMap)
		color.Green(" [done]")

		// Generate secrets
		fmt.Print("Generate secrets ...")
		secretsMap := utils.GenerateSecrets("./environment.env", template_data[index].Secrets)
		color.Green(" [done]")
		// Print secrets to console
		fmt.Println()
		fmt.Println("Following secrets were automatically generated:")
		for key, secret := range secretsMap {
			fmt.Print("   " + key + ": ")
			color.Yellow(secret)
		}
	} else {
		// Create custom stack
		utils.Heading("Let's create a custom stack for you!")
	}

	// Run if the regarding flag is set
	if flag_run || flag_demonized {
		fmt.Println()
		fmt.Println("Running docker-compose ...")
		fmt.Println()

		cmd := exec.Command("docker-compose", "up")
		if flag_demonized {
			cmd = exec.Command("docker-compose", "up", "-d")
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
