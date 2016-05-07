package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Subcommand - структура подкоманды интерфейса командой строки.
type Subcommand struct {
	// Название команды, по которому она будет вызываться.
	Name string
	// Доп. имена команды.
	Aliases []string
	// Подсказка по содержанию аргументов команды без самой команды.
	// Например "<flag1> <flag2> [filename]".
	Usage string
	// Действие команды в повелительном наклонении, с заглавной буквы,
	// с точкой на конце. Например, "Run the tests."
	Hint string
	// Указатель на обработчик команды, которому будут переданы аргументы команды
	// без имени самой команды. Возвращаемое значение сигнализирует,
	// выполнилась ли команда успешно.
	Handler func(args []string) bool
}

// DispatchSubcommand вызывает обработчик нужной команды. args должны быть
// переданы так, как они заданы в os.Args: нулевым элементом путь к бинарнику,
// дальше - аргументы командной строки.
// Возвращает true, если обработчик был найден и команда выполнилась успешно.
// В противном случае false. DispatchSubcommand сам выведет в стандартный вывод
// пример использования при некорректном наборе аргументов.
func DispatchSubcommand(args []string, subs []*Subcommand) bool {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, Usage("cmd", subs))
		return false
	}
	if len(args) < 2 {
		fmt.Fprint(os.Stderr, Usage(args[0], subs))
		return false
	}

	cmd, subcmd, args := args[0], args[1], args[2:]
	var fn func([]string) bool
	for _, sub := range subs {
		if sub.Name == subcmd {
			fn = sub.Handler
			break
		}
		for _, alias := range sub.Aliases {
			if alias == subcmd {
				fn = sub.Handler
				break
			}
		}
		if fn != nil {
			break
		}
	}

	if fn == nil {
		fmt.Fprint(os.Stderr, Usage(cmd, subs))
		return false
	}
	return fn(args)
}

// ByName имплементирует sort.Interface для сортировки списка команд по имени.
type ByName []*Subcommand

func (l ByName) Len() int {
	return len(l)
}

func (l ByName) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

func (l ByName) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Usage возвращает строку помощи по использованию команды.
func Usage(cmd string, subs []*Subcommand) string {
	head := fmt.Sprintf("Usage: %v", cmd)
	if len(subs) == 0 {
		return head + "\n"
	}
	lines := []string{
		head + " <command>",
		"",
		"Available commands:",
	}
	sort.Sort(ByName(subs))
	for _, sub := range subs {
		lines = append(lines, fmt.Sprintf("%-12v%v", sub.Name, sub.Hint))
		if len(sub.Aliases) != 0 {
			lines = append(lines,
				fmt.Sprintf("%12vAlises: %v", " ", strings.Join(sub.Aliases, ", ")),
			)
		}
		usage := fmt.Sprintf("%12vUsage: %v", " ", sub.Name)
		if sub.Usage != "" {
			usage += " " + sub.Usage
		}
		lines = append(lines, usage)
	}
	return strings.Join(lines, "\n") + "\n"
}
