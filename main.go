package main

import (
	"als/msa/configuration"
	"als/msa/util"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dixonwille/wmenu"
	"github.com/jxeng/shortcut"
)

func main() {

    argLength := len(os.Args[1:])
	if argLength == 0{
		ShowMenu()
	}else{
		StartSteam(os.Args[1])
	}
}



func ShowMenu() {

	reader := bufio.NewReader(os.Stdin)
	
	fmt.Println("Multi Steam Account")
	fmt.Println("---------------------")

	menu := wmenu.NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Action(MenuAction)

	menu.Option("Create ShortCut", nil, false, GetSteam)	

	menu.Option("Exit", nil, false, func(opt wmenu.Opt) error {
	  fmt.Printf("Bye !!!")
	  return nil
	})
	err := menu.Run()
	if err != nil{
	  log.Fatal(err)
	}
}

func MenuAction(opts []wmenu.Opt) error {
	//fmt.Println(opts[0].Text)
	return nil
}

func GetSteam(opt wmenu.Opt) error{
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Steam user name-> ")
	steam, _ := reader.ReadString('\n')
	// convert CRLF to LF
	steam = strings.Replace(steam, "\r\n", "", -1)

	result, err := CreateShortCut(steam)
	if result == false {
		log.Fatal(err)
		return err
	} else {
		fmt.Printf("Sucess!")	
		return nil
	}	
	
}

func CreateShortCut(steam string) (bool, error) {

	sc := shortcut.Shortcut{
		ShortcutPath:     configuration.DesktopPath() + "Steam " + steam + ".lnk",
		Target:           configuration.CurrentPath(),
		IconLocation:     configuration.GetSteamPath() + ",0" ,
		Arguments:        steam,
		Description:      "Steam Multi Account",
		Hotkey:           "",
		WindowStyle:      "1",
		WorkingDirectory: configuration.GetSteamDir(),
	}
	err := shortcut.Create(sc)
	if err != nil {
		return false, err
	} else {
		return true, nil	
	}
}

func StartSteam(steam string){
	fmt.Println("Starting Steam "  + steam)
	
	configuration.SaveReg(steam)

	explorer := util.FindProcessByName("Steam.exe")
	if explorer != nil {
		// found it
		process := os.Process{Pid: explorer.ProcessID}
		util.TerminateProcess(process,0)
	}	
	
	cmdToRun := configuration.GetSteamPath()
	args := []string{"", ""}
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	if process, err := os.StartProcess(cmdToRun, args, procAttr); err != nil {
		fmt.Printf("Unable to run %s: %s\n", cmdToRun, err.Error())
	} else {
		fmt.Printf("%s running as pid %d\n", cmdToRun, process.Pid)
	}
}

