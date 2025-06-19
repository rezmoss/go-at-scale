// Example 54
package main

import (
	"errors"
	"fmt"
)

type Command interface {
	Execute() error
	Undo() error
}

type DocumentEditor struct {
	content string
}

type InsertTextCommand struct {
	editor *DocumentEditor
	pos    int
	text   string
}

func (c *InsertTextCommand) Execute() error {
	c.editor.content = c.editor.content[:c.pos] + c.text + c.editor.content[c.pos:]
	return nil
}

func (c *InsertTextCommand) Undo() error {
	c.editor.content = c.editor.content[:c.pos] + c.editor.content[c.pos+len(c.text):]
	return nil
}

type CommandInvoker struct {
	commands []Command
	current  int
}

func (i *CommandInvoker) Execute(cmd Command) error {
	if err := cmd.Execute(); err != nil {
		return err
	}
	i.commands = append(i.commands[:i.current], cmd)
	i.current++
	return nil
}

func (i *CommandInvoker) Undo() error {
	if i.current == 0 {
		return errors.New("nothing to undo")
	}
	i.current--
	return i.commands[i.current].Undo()
}

func main() {
	// Create a document editor
	editor := &DocumentEditor{content: ""}

	// Create a command invoker
	invoker := &CommandInvoker{}

	// Create and execute commands
	cmd1 := &InsertTextCommand{editor: editor, pos: 0, text: "Hello"}
	err := invoker.Execute(cmd1)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}

	cmd2 := &InsertTextCommand{editor: editor, pos: 5, text: " World"}
	err = invoker.Execute(cmd2)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}

	fmt.Printf("Content after commands: %s\n", editor.content)

	// Undo the last command
	err = invoker.Undo()
	if err != nil {
		fmt.Printf("Error undoing command: %v\n", err)
	}

	fmt.Printf("Content after undo: %s\n", editor.content)

	// Undo again
	err = invoker.Undo()
	if err != nil {
		fmt.Printf("Error undoing command: %v\n", err)
	}

	fmt.Printf("Content after second undo: %s\n", editor.content)

	// Try to undo when there's nothing left
	err = invoker.Undo()
	if err != nil {
		fmt.Printf("Error undoing command: %v\n", err)
	}
}