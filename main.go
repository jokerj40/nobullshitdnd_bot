package main

import (
	"log"
	"time"
	"regexp"
	"strconv"
	"fmt"
	"io/ioutil"
	"math/rand"
	"bytes"

	tele "gopkg.in/telebot.v3"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	TOKEN, err := ioutil.ReadFile("credentials.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	pref := tele.Settings{
		Token:  string(bytes.TrimSpace(TOKEN)),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tele.OnText, func(c tele.Context) error {
		re := regexp.MustCompile("\\.[0-9]?\\,[0-9]{1,2}((\\+|\\-)[0-9]{1,2})?")
		formula := c.Text()

		if !re.MatchString(formula) {
			return nil
		}

		formula = formula[1:]

		re_quant := regexp.MustCompile("[0-9]\\,")
		str_quant := re_quant.FindString(formula)
		var quant int = 1
		if str_quant != "" {
			str_quant = str_quant[:1]
			quant, err = strconv.Atoi(str_quant)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			formula = formula[1:]
		}


		re_dice := regexp.MustCompile("\\,[0-9]{1,2}")
		str_dice := re_dice.FindString(formula)
		var dice int
		if str_dice != "" {
			dice, err = strconv.Atoi(str_dice[1:])
			if err != nil {
				fmt.Println(err)
				return nil
			}
			if len(str_dice) == 3 {
				formula = formula[3:]
			} else if len(str_dice) == 2 {
				formula = formula[2:]
			} else {
				return c.Send("[Error] `d` is too big")
			}
		}

		var mod int = 0
		var str_mod string
		if formula != "" {
			re_mod := regexp.MustCompile("(\\+|\\-)[0-9]{1,2}")
			str_mod = re_mod.FindString(formula)
			if str_mod != "" {
				mod, err = strconv.Atoi(str_mod[1:])
				if err != nil {
					fmt.Println(err)
					return nil
				}
				if str_mod[0] == byte('-') {
					mod = -mod
				}
			}
		}

		Answer := fmt.Sprintf("🎲 %dd%d%s\n\n", quant, dice, str_mod) // Unknown symbol is 'dice' emoji
		nums := make([]int, 0)
		for i := 0; i < quant; i++ {
			num := rand.Intn(dice) + 1
			nums = append(nums, num+mod)
			Answer += fmt.Sprintf("rolled %4d\n", num)
		}

		Answer += "\n["
		var sum int = 0
		for i, n := range nums {
			if i > 0 {
				Answer += ", "
			}
			Answer += fmt.Sprintf("%d", n)
			sum += n
		}
		Answer += "]\n"
		Answer += fmt.Sprintf("\nΣ = %d", sum)

		return c.Send(Answer)
	})
	

	b.Start()
}

