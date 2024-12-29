package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"
)

type CmdFlags struct {
	Add    string
	Del    int
	Edit   string
	Toggle int
	List   bool
}

func NewCmdFlags() *CmdFlags {
	cf := CmdFlags{}
	flag.StringVar(&cf.Add, "add", "", "Add a new Todo specify title.")
	flag.StringVar(&cf.Edit, "edit", "", "Delete a Todo specify title.")

	flag.BoolVar(&cf.List, "list", false, "List all the Todos.")

	flag.IntVar(&cf.Del, "del", -1, "Delete any specify todo by index.")
	flag.IntVar(&cf.Toggle, "toggle", -1, "Toggle any specify todo.")

	flag.Parse()

	return &cf
}

func (cf *CmdFlags) Excute(todos *Todos) {
	switch {

	case cf.List:
		todos.print()

	case cf.Add != "":
		todos.add(cf.Add)
		todos.print()
	case cf.Edit != "":
		parts := strings.SplitN(cf.Edit, ":", 2)
		if len(parts) != 2 {
			log.Errorf("Error, Invalid format for edit. Please use id:New Title")
			os.Exit(1)
		}
		index, err := strconv.Atoi(parts[0])

		if err != nil {
			log.Errorf("Invalid index for edit")

		}
		todos.edit(index, parts[1])
		todos.print()
	case cf.Toggle != -1:
		todos.toggle(cf.Toggle)
		todos.print()
	case cf.Del != -1:
		todos.delete(cf.Del)
		todos.print()
	default:
		fmt.Println("Invalide command")
	}
}
