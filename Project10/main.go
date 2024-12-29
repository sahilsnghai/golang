package main

func main() {
	todos := Todos{}

	storage := NewStorage[Todos]("todos.json")
	storage.Load(&todos)

	cmdFlags := NewCmdFlags()

	cmdFlags.Excute(&todos)
	// todos.add("Hey")
	// todos.add("Hey2")

	// todos.toggle(0)
	// todos.print()

	storage.Save(todos)

}
