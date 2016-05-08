package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/yamnikov-oleg/pss/Godeps/_workspace/src/github.com/howeyc/gopass"
	"github.com/yamnikov-oleg/pss/lib/pss"
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
		"",
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
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n") + "\n"
}

func openStorage() (pss.Storage, string, bool) {
	fmt.Print("Master-password: ")
	pwd := string(gopass.GetPasswd())

	stg, err := pss.DecryptDefault(pwd)
	if err == pss.WrongPwdErr {
		fmt.Println("Wrong password")
		return nil, "", false
	} else if os.IsNotExist(err) {
		fmt.Println("Error: storage file does not exist.")
		fmt.Printf("Try calling `%v init` for setup.\n", os.Args[0])
		return nil, "", false
	} else if err != nil {
		fmt.Println("Error opening storage:")
		fmt.Println(err)
		return nil, "", false
	}

	return stg, pwd, true
}

func saveStorage(stg pss.Storage, pwd string) bool {
	if err := pss.EncryptDefault(stg, pwd); err != nil {
		fmt.Println("Error writing storage:")
		fmt.Println(err)
		return false
	}
	return true
}

func printRecord(r *pss.Record) {
	fmt.Printf("%24v @ %v\n", r.Username, r.Website)
}

func findRecord(stg pss.Storage, args []string) (*pss.Record, bool) {
	var (
		website  string
		username string
	)
	if len(args) < 1 || len(args) > 2 {
		fmt.Println("Please, provide website and username via arguments")
		return nil, false
	}
	website = args[0]
	if len(args) == 2 {
		username = args[1]
	}

	var results []*pss.Record
	for _, rec := range stg {
		siteMatch := strings.Contains(rec.Website, website)
		nameMatch := strings.Contains(rec.Username, username)
		crossMatch := strings.Contains(rec.Username, website)
		if username == "" && (siteMatch || crossMatch) {
			results = append(results, rec)
			continue
		}
		if username != "" && siteMatch && nameMatch {
			results = append(results, rec)
			continue
		}
	}

	if len(results) == 0 {
		query := ""
		if username != "" {
			query = fmt.Sprintf("%v at %v", username, website)
		} else {
			query = website
		}
		fmt.Printf("Could not find a record of %v\n", query)
		return nil, false
	}
	if len(results) > 1 {
		fmt.Println("Multiple choices occured:")
		for _, r := range results {
			printRecord(r)
		}
		fmt.Println("Please, make your request more specific")
		return nil, false
	}

	return results[0], true
}

func openAndFindRecord(args []string) (*pss.Record, pss.Storage, string, bool) {
	if len(args) < 1 || len(args) > 2 {
		fmt.Println("Please, provide website and username via arguments")
		return nil, nil, "", false
	}
	stg, pwd, ok := openStorage()
	if !ok {
		return nil, nil, "", false
	}
	rec, ok := findRecord(stg, args)
	if !ok {
		return nil, nil, "", false
	}
	return rec, stg, pwd, true
}

func promtPwd(text string, allowEmpty bool) (pwd string) {
	for {
		fmt.Print(text + ": ")
		pwd = string(gopass.GetPasswd())
		fmt.Print("Repeat: ")
		pwd2 := string(gopass.GetPasswd())
		if pwd != pwd2 {
			fmt.Println("Passwords do not match!")
			continue
		}
		if pwd == "" && !allowEmpty {
			fmt.Println("Empty password is not allowed.")
			continue
		}
		return
	}
}
