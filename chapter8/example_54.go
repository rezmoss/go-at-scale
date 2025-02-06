// Example 54
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